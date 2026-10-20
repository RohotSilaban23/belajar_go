package middleware

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
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

	publicKeyData, err := os.ReadFile("keys/public.pem")
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return &JWTService{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
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

func (j *JWTService) ValidateToken(
	tokenString string,
) (int64, string, error) {

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodRS256 {
				return nil, fmt.Errorf("algoritma JWT tidak valid")
			}

			return j.PublicKey, nil
		},
	)

	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", fmt.Errorf("token tidak valid")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", fmt.Errorf("user_id tidak ditemukan")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return 0, "", fmt.Errorf("email tidak ditemukan")
	}

	return int64(userID), email, nil
}
