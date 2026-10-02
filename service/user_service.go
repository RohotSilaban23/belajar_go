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

func NewUserService(userRepository *repository.UserRepository, jwtService *middleware.JWTService, redisService *redis.RedisService) *UserService {
	return &UserService{
		UserRepository: userRepository,
		JWTService:     jwtService,
		RedisService:   redisService,
	}
}

func (s *UserService) CreateUser(user *model.User) (*model.User, error) {

	existingUser, err := s.UserRepository.FindByEmai(user.Email)

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

func (s *UserService) Login(email string, password string) (string, error) {

	user, err := s.UserRepository.FindByEmai(email)
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

	token, err := s.JWTService.GenerateToken(user.ID, user.Name)
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
