package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s [binding_key]...", os.Args[0])
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	exchange := "direct_logs"
	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	failOnError(err, "Failed to declare exchange")

	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	failOnError(err, "Failed to declare a queue")

	for _, key := range os.Args[1:] {
		err = ch.QueueBind(q.Name, key, exchange, false, nil)
		failOnError(err, "Failed to bind queue")
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	log.Printf(" [*] Waiting for logs. Bound to : %v", os.Args[1:])
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s: %s", d.RoutingKey, d.Body)
		}
	}()

	<-forever
}
