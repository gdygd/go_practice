package main

import (
	"log"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

// go run emit_topic.go "kern.critical" "Kernel panic!"
// go run emit_topic.go "auth.info" "User logged in"
// go run emit_topic.go "order.created.eu" "Order #123 created"

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s<routeing_key> <message...>", os.Args[0])
	}

	routingKey := os.Args[1]
	body := strings.Join(os.Args[2:], " ")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	exchange := "topic_logs"
	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,  // durable: 운영 시 true 권장
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	err = ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s: %s", routingKey, body)
}
