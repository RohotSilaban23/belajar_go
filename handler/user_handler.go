package handler

import (
	"belajar_go/model"
	"belajar_go/service"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	UserService *service.UserService
	OtpService  *service.ForgotPasswordService
}

func NewUserHandler(userService *service.UserService, otpService *service.ForgotPasswordService) *UserHandler {
	return &UserHandler{
		UserService: userService,
		OtpService:  otpService,
	}
}

func (h *UserHandler) CreateUser(c *echo.Context) error {

	var user model.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Request tidak valid",
		})
	}

	CreateUser, err := h.UserService.CreateUser(&user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	_, err = h.OtpService.CreateOTP(
		c.Request().Context(),
		user.Email,
	)

	if err != nil {
		log.Println("DEBUG CreateOTP gagal:", err)

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Akun berhasil dibuat, tetapi OTP gagal dikirim",
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Registrasi berhasil diproses. Silakan periksa email untuk OTP verifikasi.",
		"email":   CreateUser.Email,
	})
}

func (h *UserHandler) FindByID(c *echo.Context) error {

	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "ID tidak valid",
		})
	}

	user, err := h.UserService.FindByID(id)

	if err != nil {
		log.Println("DEBUG FindByID - ID:", id)
		log.Println("DEBUG FindByID - Error:", err)
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User tidak ditemukan",
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "ID tidak valid",
		})
	}

	err = h.UserService.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User berhasil dihapus",
	})
}

func (h *UserHandler) UpdateUser(c *echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "ID tidak valid",
		})
	}

	var user model.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Request tidak valid",
		})
	}

	user.ID = id

	updatedUser, err := h.UserService.UpdateUser(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) ConfirmRegisterOTP(c *echo.Context) error {
	var request model.RequestOtpRegister

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Request tidak valid",
		})
	}

	if request.Email == "" || request.OTP == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Email dan OTP wajib diisi",
		})
	}

	err := h.OtpService.ConfirmRegisterOTP(
		c.Request().Context(),
		request.Email,
		request.OTP,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Konfirmasi berhasil. Akun kamu sekarang aktif.",
	})
}
