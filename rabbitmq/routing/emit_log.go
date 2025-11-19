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

func severityFrom(args []string) string {
	if len(args) < 2 {
		return "info"
	}
	return args[1]
}

func bodyFrom(args []string) string {
	if len(args) < 3 {
		return "Hello World!"
	}
	return strings.Join(args[2:], " ")
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Faield to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	exchange := "direct_logs"

	err = ch.ExchangeDeclare(
		exchange,
		"direct",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	)

	failOnError(err, "Failed to decleare an exchange")

	severity := severityFrom(os.Args)
	body := bodyFrom(os.Args)

	err = ch.Publish(
		exchange,
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			// DeliveryMode: amqp.Persistent,
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Printf("[*] Sent %s: %S", severity, body)
}
