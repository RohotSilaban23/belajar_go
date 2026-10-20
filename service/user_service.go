package service

import (
	"belajar_go/middleware"
	"belajar_go/model"
	"belajar_go/redis"
	"belajar_go/repository"
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository *repository.UserRepository
	JWTService     *middleware.JWTService
	RedisService   *redis.RedisService
}

func NewUserService(
	userRepository *repository.UserRepository,
	jwtService *middleware.JWTService,
	redisService *redis.RedisService) *UserService {
	return &UserService{
		UserRepository: userRepository,
		JWTService:     jwtService,
		RedisService:   redisService,
	}
}

func (s *UserService) CreateUser(user *model.User) (*model.User, error) {

	existingUser, err := s.UserRepository.FindByEmail(user.Email)

	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email sudah digunakan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	user.Role = string("menber")
	user.IsActive = false

	return s.UserRepository.Create(user)
}

func (s *UserService) FindByID(id int64) (*model.User, error) {
	return s.UserRepository.FindByID(id)
}

func (s *UserService) DeleteUser(id int64) error {
	return s.UserRepository.DeleteUser(id)
}

func (s *UserService) UpdateUser(user *model.User) (*model.User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	return s.UserRepository.UpdateUser(user)
}

func (s *UserService) CheckUserActive(user *model.User) error {
	if !user.IsActive {
		return fmt.Errorf("akun tidak aktif")
	}

	return nil
}

func (s *UserService) Login(email string, password string) (string, error) {

	user, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return "", err
	}

	err = s.CheckUserActive(user)
	if err != nil {
		return "", err
	}

	token, err := s.JWTService.GenerateToken(user.ID, user.Nama)
	if err != nil {
		return "", err
	}

	err = s.RedisService.SaveToken(
		context.Background(),
		token,
		user.ID,
		24*time.Hour,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) Logout(
	token string,
) error {

	return s.RedisService.DeleteToken(
		context.Background(),
		token,
	)
}

func (s *UserService) RefreshToken(
	oldToken string,
) (string, error) {

	ctx := context.Background()

	// 1. Validasi JWT lama.
	id, email, err := s.JWTService.ValidateToken(oldToken)
	if err != nil {
		return "", fmt.Errorf("token tidak valid atau kedaluwarsa")
	}

	// 2. Periksa apakah token lama masih aktif di Redis.
	active, err := s.RedisService.IsTokenActive(ctx, oldToken)
	if err != nil {
		return "", fmt.Errorf("gagal mengecek token: %w", err)
	}

	if !active {
		return "", fmt.Errorf("token sudah tidak aktif")
	}

	// 3. Buat JWT baru.
	newToken, err := s.JWTService.GenerateToken(id, email)
	if err != nil {
		return "", err
	}

	// 4. Simpan JWT baru ke Redis.
	err = s.RedisService.SaveToken(
		ctx,
		newToken,
		id,
		24*time.Hour,
	)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan token baru: %w", err)
	}

	// 5. Hapus JWT lama dari Redis.
	err = s.RedisService.DeleteToken(ctx, oldToken)
	if err != nil {
		// Batalkan token baru jika token lama gagal dihapus.
		_ = s.RedisService.DeleteToken(ctx, newToken)
		return "", fmt.Errorf("gagal menonaktifkan token lama: %w", err)
	}

	// 6. Kembalikan JWT baru.
	return newToken, nil
}
