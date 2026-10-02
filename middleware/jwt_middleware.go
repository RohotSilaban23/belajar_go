package middleware

import (
	"belajar_go/redis"
	"crypto/rsa"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type JWTMiddleware struct {
	PublicKey    *rsa.PublicKey
	RedisService *redis.RedisService
}

func NewJWTMiddleware(redisService *redis.RedisService) (*JWTMiddleware, error) {

	publicKeyData, err := os.ReadFile("keys/public.pem")
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return &JWTMiddleware{
		PublicKey:    publicKey,
		RedisService: redisService,
	}, nil
}

func (m *JWTMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c *echo.Context) error {

		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Token tidak ditemukan",
			})
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Format token tidak valid",
			})
		}

		tokenString := parts[1]

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				// Pastikan token menggunakan RS256
				if token.Method != jwt.SigningMethodRS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return m.PublicKey, nil
			},
		)

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Token tidak valid atau sudah expired",
			})
		}

		exists, err := m.RedisService.Client.Exists(
			c.Request().Context(),
			"active_token:"+tokenString,
		).Result()

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Gagal mengecek token",
			})
		}
		if exists == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Token sudah tidak aktif",
			})
		}

		return next(c)
	}
}
