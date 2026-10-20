package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

const OTPQueue = "otp_email_queue"

type OTPMessage struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type Producer struct {
	Channel *amqp.Channel
}

func NewProducer(amqpURL string) (*Producer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = channel.QueueDeclare(
		OTPQueue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	return &Producer{
		Channel: channel,
	}, nil
}

func (p *Producer) PublishOTP(
	ctx context.Context,
	email string,
	otp string,
) error {
	body, err := json.Marshal(OTPMessage{
		Email: email,
		OTP:   otp,
	})
	if err != nil {
		return err
	}

	return p.Channel.PublishWithContext(
		ctx,
		"",       // default exchange
		OTPQueue, // nama queue
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
