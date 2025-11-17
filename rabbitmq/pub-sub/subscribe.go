package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	const exchange = "logs"
	err = ch.ExchangeDeclare(
		exchange,
		"fanout", // publisher와 동일
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	q, err := ch.QueueDeclare(
		"",    // name: empty → server-named
		false, // durable
		true,  // autodelete
		true,  // exclusive
		false, // no-wait
		nil,
	)

	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)
	failOnError(err, "Failed to bind queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true, // autoAck
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register consumer")
	log.Println(" [*] Waiting for logs. CTRL+C to exit")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [*] Received: %s", d.Body)
		}
	}()

	<-forever
}
