package service

import (
	"belajar_go/rabbitmq"
	"belajar_go/redis"
	"belajar_go/repository"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ForgotPasswordService struct {
	UserRepository   *repository.UserRepository
	RedisService     *redis.RedisService
	RabbitMQProducer *rabbitmq.Producer
}

func NewForgotPasswordService(
	userRepository *repository.UserRepository,
	redisService *redis.RedisService,
	producer *rabbitmq.Producer,
) *ForgotPasswordService {
	return &ForgotPasswordService{
		UserRepository:   userRepository,
		RedisService:     redisService,
		RabbitMQProducer: producer,
	}
}

// Membuat OTP enam digit.
func generateOTP() (string, error) {
	number, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", number.Int64()), nil
}

// Mengubah OTP menjadi hash.
func hashOTP(otp string) string {
	hash := sha256.Sum256([]byte(otp))
	return fmt.Sprintf("%x", hash)
}

// Memeriksa email, membuat OTP, lalu menyimpan hash ke Redis.
func (s *ForgotPasswordService) CreateOTP(
	ctx context.Context,
	email string,
) (string, error) {

	// Pastikan email terdaftar.
	_, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("email tidak ditemukan")
		}
		return "", err
	}

	// if otpType != "forgot_password" &&
	// 	otpType != "register_confirmation" {
	// 	return "", fmt.Errorf("jenis OTP tidak valid")
	// }

	// Buat OTP.
	otp, err := generateOTP()
	if err != nil {
		return "", err
	}

	// Simpan hash OTP selama 5 menit.
	err = s.RedisService.SaveOTP(
		ctx,
		email,
		hashOTP(otp),
		5*time.Minute,
	)
	if err != nil {
		return "", err
	}

	err = s.RabbitMQProducer.PublishOTP(ctx, email, otp)
	if err != nil {
		// Hapus OTP agar pengguna tidak memiliki kode
		// yang gagal dijadwalkan untuk dikirim.
		_ = s.RedisService.DeleteOTP(ctx, email)
		return "", err
	}

	// OTP asli nantinya dikirim melalui RabbitMQ.
	return otp, nil
}

func (s *ForgotPasswordService) ChangePassword(
	ctx context.Context,
	email string,
	otp string,
	newPassword string,
	confirmPassword string,
) error {
	if newPassword == "" || confirmPassword == "" {
		return fmt.Errorf("password wajib diisi")
	}

	if newPassword != confirmPassword {
		return fmt.Errorf("konfirmasi password tidak cocok")
	}

	// Ambil hash OTP dari Redis.
	storedHash, err := s.RedisService.GetOTP(ctx, email)
	if err != nil {
		return fmt.Errorf("OTP tidak ditemukan atau sudah kedaluwarsa")
	}

	// Bandingkan OTP yang dikirim pengguna.
	if hashOTP(otp) != storedHash {
		return fmt.Errorf("OTP tidak valid")
	}

	// Cari akun berdasarkan email.
	user, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("akun tidak ditemukan")
	}

	// Hash password baru dengan bcrypt.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// Simpan password baru ke database.
	_, err = s.UserRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	// OTP hanya boleh digunakan sekali.
	if err := s.RedisService.DeleteOTP(ctx, email); err != nil {
		return err
	}

	return nil
}

func (s *ForgotPasswordService) ConfirmRegisterOTP(
	ctx context.Context,
	email string,
	otp string,
) error {
	// 1. Cari user berdasarkan email.
	user, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("akun tidak ditemukan")
	}

	// 2. Jangan proses akun yang sudah aktif.
	if user.IsActive {
		return fmt.Errorf("akun sudah aktif")
	}

	// 3. Ambil hash OTP dari Redis.
	storedHash, err := s.RedisService.GetOTP(ctx, email)
	if err != nil {
		return fmt.Errorf("OTP tidak ditemukan atau sudah kedaluwarsa")
	}

	// 4. Cocokkan OTP yang dikirim dengan hash di Redis.
	if hashOTP(otp) != storedHash {
		return fmt.Errorf("OTP tidak valid")
	}

	// 5. Aktifkan akun.
	if err := s.UserRepository.ActivateUser(email); err != nil {
		return err
	}

	// 6. Hapus OTP setelah berhasil diverifikasi.
	if err := s.RedisService.DeleteOTP(ctx, email); err != nil {
		return fmt.Errorf("akun aktif, tetapi OTP gagal dihapus")
	}

	return nil
}
