package utils

import (
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp091.Connection
var RabbitChannel *amqp091.Channel

func InitRabbitMQ(url string) (*amqp091.Connection, error) {
	var conn *amqp091.Connection
	var err error

	for i := 0; i < 10; i++ {
		conn, err = amqp091.Dial(url)
		if err == nil {
			fmt.Println("Connected to RabbitMQ!")
			RabbitConn = conn

			RabbitChannel, err = conn.Channel()
			if err != nil {
				return nil, err
			}

			return conn, nil
		}

		fmt.Println("RabbitMQ not ready, retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after retries: %w", err)
}

func Publish(queueName string, body []byte) error {
	_, err := RabbitChannel.QueueDeclare(
		queueName,
		true,
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	return RabbitChannel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
