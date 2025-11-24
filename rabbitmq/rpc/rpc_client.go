package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run rpc_client.go <n>")
	}
	n := os.Args[1]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	replyQ, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // autoDelete
		true,  // exclusive
		false, // noWait
		nil,   // args
	)
	failOnError(err, "Failed to declare a reply queue")
	msgs, err := ch.Consume(
		replyQ.Name,
		"",    // consumer
		true,  // autoAck = true (간단한 예제. 필요한 경우 false로 처리)
		false, // exclusive
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	corrId := uuid.New().String()

	err = ch.Publish(
		"",          // default exchange
		"rpc_queue", // routing key = queue name
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       replyQ.Name,
			Body:          []byte(n),
		},
	)
	failOnError(err, "Failed to publish a message")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		select {
		case d := <-msgs:
			if d.CorrelationId == corrId {
				log.Printf("Got response: %s", d.Body)
				return
			}
		case <-ctx.Done():
			log.Println("Timeout waiting for RPC response")
			return
		}
	}
}
