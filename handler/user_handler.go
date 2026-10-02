package handler

import (
	"belajar_go/model"
	"belajar_go/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
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

	return c.JSON(http.StatusCreated, CreateUser)
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
