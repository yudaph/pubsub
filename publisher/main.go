package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main() {
	now := time.Now()
	nowString := now.Format(time.DateTime)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = ch.Publish(
		"",          // exchange
		"queueName", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(nowString),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err = redis.Publish(context.Background(), "channel", nowString).Err()
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
}
