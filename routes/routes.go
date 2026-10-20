package routes

import (
	"belajar_go/handler"
	middleware "belajar_go/middleware"
	"net/http"

	"github.com/labstack/echo/v5"
)

func Init(e *echo.Echo,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	forgetHandler *handler.ForgotPasswordHandler,
	jwtMiddleware *middleware.JWTMiddleware) {

	// Route

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, "Hello, this Echo and rohot")
	})
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.FindByID, jwtMiddleware.Authenticate)
	e.DELETE("/users/:id", userHandler.DeleteUser, jwtMiddleware.Authenticate)
	e.PUT("/users/:id", userHandler.UpdateUser, jwtMiddleware.Authenticate)

	e.POST("/login", authHandler.Login)
	e.POST("/logout", authHandler.Logout, jwtMiddleware.Authenticate)
	e.POST("/forget_password", forgetHandler.ForgotPassword)
	e.POST("/change-password", forgetHandler.ChangePassword)
	e.POST("/refresh-Token", authHandler.RefreshToken)
	e.POST("/users/confirm-register-otp", userHandler.ConfirmRegisterOTP)
}
