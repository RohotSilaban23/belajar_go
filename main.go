package main

import (
	"belajar_go/database"
	"belajar_go/handler"
	"belajar_go/middleware"
	"belajar_go/rabbitmq"
	"belajar_go/redis"
	"belajar_go/repository"
	"context"
	"log"

	"belajar_go/routes"
	"belajar_go/service"
	"fmt"

	"github.com/labstack/echo/v5"
)

func main() {
	db, err := database.InitDB()

	if err := database.SeedAdmin(db); err != nil {
		log.Fatal("Gagal melakukan seeding:", err)
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

	authHandler := handler.NewAuthHandler(userService)
	jwtMiddleware, err := middleware.NewJWTMiddleware(redisService)
	if err != nil {
		panic(err)
	}
	producer, err := rabbitmq.NewProducer(
		"amqp://alert_app:Silaban123@127.0.0.1:5672/",
	)
	if err != nil {
		log.Fatal("Gagal terhubung ke RabbitMQ: ", err)
	}
	log.Println("RabbitMQ berhasil terhubung")

	forgotPasswordService := service.NewForgotPasswordService(
		userRepository,
		redisService,
		producer,
	)
	userHandler := handler.NewUserHandler(userService, forgotPasswordService)

	forgetHandler := handler.NewForgotPasswordHandler(
		forgotPasswordService,
	)

	// Start server
	e := echo.New()
	routes.Init(e, userHandler, authHandler, forgetHandler, jwtMiddleware)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
