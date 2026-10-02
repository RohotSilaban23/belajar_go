package middleware

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	PrivateKey *rsa.PrivateKey
}

func NewJWTService() (*JWTService, error) {

	privateKeyData, err := os.ReadFile("keys/private.pem")
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	return &JWTService{
		PrivateKey: privateKey,
	}, nil
}

func (j *JWTService) GenerateToken(userID int64, email string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)

	tokenString, err := token.SignedString(j.PrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
