package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go consumeRedis(ctx)
	go consumeRabbitMQ(ctx)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}

func consumeRedis(ctx context.Context) {
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pubsub := redis.Subscribe(ctx, "channel")
	defer pubsub.Close()

	ch := pubsub.Channel()

	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongo.Disconnect(ctx)

	err = mongo.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping to MongoDB: %v", err)
	}

	db := mongo.Database("database")

	collection := db.Collection("collection")

	for msg := range ch {
		_, err := collection.InsertOne(ctx, bson.M{"message": msg.Payload})
		if err != nil {
			log.Fatalf("Failed to insert to MongoDB: %v", err)
		}
	}
}

func consumeRabbitMQ(ctx context.Context) {
	rabbitMq, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMq.Close()

	ch, err := rabbitMq.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"queueName", // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongo.Disconnect(ctx)

	err = mongo.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping to MongoDB: %v", err)
	}

	db := mongo.Database("database")

	collection := db.Collection("collection")

	for d := range msgs {
		_, err := collection.InsertOne(ctx, bson.M{"message": string(d.Body)})
		if err != nil {
			log.Fatalf("Failed to insert to MongoDB: %v", err)
		}
	}
}
