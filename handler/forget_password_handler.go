package handler

import (
	"belajar_go/service"
	"net/http"

	"github.com/labstack/echo/v5"
)

type ForgotPasswordHandler struct {
	ForgotPasswordService *service.ForgotPasswordService
}

type ChangePasswordRequest struct {
	Email           string `json:"email"`
	OTP             string `json:"otp"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

func NewForgotPasswordHandler(
	forgotPasswordService *service.ForgotPasswordService,
) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		ForgotPasswordService: forgotPasswordService,
	}
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func (h *ForgotPasswordHandler) ForgotPassword(
	c *echo.Context,
) error {
	var request ForgotPasswordRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Format request tidak valid",
		})
	}

	if request.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Email wajib diisi",
		})
	}

	otp, err := h.ForgotPasswordService.CreateOTP(
		c.Request().Context(),
		request.Email,
	)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Permintaan forgot password gagal",
		})
	}

	// OTP belum dikirim melalui email.
	// Variabel ini nantinya digunakan dalam integrasi RabbitMQ.
	_ = otp

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Permintaan berhasil diproses. Jika email terdaftar, OTP akan dikirim.",
	})
}

func (h *ForgotPasswordHandler) ChangePassword(
	c *echo.Context,
) error {
	var request ChangePasswordRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Format request tidak valid",
		})
	}

	if request.Email == "" ||
		request.OTP == "" ||
		request.NewPassword == "" ||
		request.ConfirmPassword == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Semua field wajib diisi",
		})
	}

	err := h.ForgotPasswordService.ChangePassword(
		c.Request().Context(),
		request.Email,
		request.OTP,
		request.NewPassword,
		request.ConfirmPassword,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Password berhasil diubah",
	})
}
