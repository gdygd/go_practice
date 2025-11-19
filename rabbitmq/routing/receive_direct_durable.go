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
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <queue_name> <binding_key>...", os.Args[0])
	}
	queueName := os.Args[0]
	bindingKeys := os.Args[2:]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabiitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	exchange := "direct_logs"
	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	failOnError(err, "Failed to declare exchange")

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queeu")

	for _, key := range bindingKeys {
		err = ch.QueueBind(q.Name, key, exchange, false, nil)
		failOnError(err, "Failed to bind queue")
	}

	err = ch.Qos(1, 0, false)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register consumer")

	log.Printf(" [*] Waiting for messages in %s. Bound to: %v", q.Name, bindingKeys)
	for d := range msgs {
		log.Printf(" [x] Received %s: %s", d.RoutingKey, d.Body)
		d.Ack(false)
	}
}
