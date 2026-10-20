package main

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	password := os.Getenv("RABBITMQ_PASSWORD")
	if password == "" {
		log.Fatal("Environment variable RABBITMQ_PASSWORD belum diisi")
	}

	url := fmt.Sprintf(
		"amqp://alert_app:%s@127.0.0.1:5672/",
		password,
	)

	fmt.Println("Mencoba menghubungkan ke RabbitMQ...")

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Koneksi gagal: %v", err)
	}
	defer conn.Close()

	fmt.Println("Koneksi RabbitMQ berhasil!")

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Gagal membuka channel: %v", err)
	}
	defer channel.Close()

	fmt.Println("Channel RabbitMQ berhasil dibuat!")
}
