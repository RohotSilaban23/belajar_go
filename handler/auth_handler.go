package handler

import (
	"net/http"
	"strings"

	"belajar_go/model"
	"belajar_go/service"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	UserService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		UserService: userService,
	}
}

func (h *AuthHandler) Login(c *echo.Context) error {

	var req model.RequestLogin

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Request tidak valid",
		})
	}

	user, err := h.UserService.Login(
		req.Email,
		req.Password,
	)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Email atau password salah",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login berhasil",
		"Token":   user,
	})
}

func (h *AuthHandler) Logout(c *echo.Context) error {

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

	err := h.UserService.Logout(tokenString)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Gagal logout",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logout berhasil",
	})
}

func (h *AuthHandler) RefreshToken(c *echo.Context) error {
	// 1. Ambil token dari header Authorization.
	authHeader := c.Request().Header.Get("Authorization")

	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Token tidak ditemukan",
		})
	}

	// 2. Pisahkan Bearer dan JWT.
	parts := strings.Fields(authHeader)

	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Format token tidak valid",
		})
	}

	oldToken := parts[1]

	// 3. Kirim token lama ke Service.
	newToken, err := h.UserService.RefreshToken(oldToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Token tidak valid, kedaluwarsa, atau tidak aktif",
		})
	}

	// 4. Kembalikan token baru.
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Token berhasil diperbarui",
		"token":   newToken,
	})
}
