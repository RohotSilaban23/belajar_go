package main

import (
	"belajar_go/database"
	"belajar_go/handler"
	"belajar_go/middleware"
	"belajar_go/redis"
	"belajar_go/repository"
	"context"

	"belajar_go/routes"
	"belajar_go/service"
	"fmt"

	"github.com/labstack/echo/v5"
)

func main() {
	db, err := database.InitDB()

	if err != nil {
		panic(err)
	}

	defer db.Close()
	fmt.Println("Database berhasil Connet")

	redisService := redis.NewRedisService()

	err = redisService.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis berhasil terhubung")

	//Repository
	userRepository := repository.NewUserRepository(db)
	//Service
	jwtService, err := middleware.NewJWTService()
	if err != nil {
		panic(err)
	}

	userService := service.NewUserService(userRepository, jwtService, redisService)
	//Handler
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)
	jwtMiddleware, err := middleware.NewJWTMiddleware(redisService)
	if err != nil {
		panic(err)
	}

	// Start server
	e := echo.New()
	routes.Init(e, userHandler, authHandler, jwtMiddleware)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
