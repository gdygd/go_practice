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
	// usage: go run receive_topic.go <binding_key>...
	// examples:
	// go run receive_topic.go "kern.*"
	// go run receive_topic.go "*.critical"
	// go run receive_topic.go "order.*.eu" "order.*.us"

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <binding_key>...", os.Args[0])
	}
	bindingKeys := os.Args[1:]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failedto open a channel")
	defer ch.Close()

	exchange := "topic_logs"
	err = ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	failOnError(err, "Failed to declare exchange")

	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	failOnError(err, "Failed to declare a queue")

	for _, key := range bindingKeys {
		err = ch.QueueBind(q.Name, key, exchange, false, nil)
		failOnError(err, "Faield to bind queue")
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Faield to register a consumer")

	log.Printf(" [*] Waiting for logs. Bound to: %v", bindingKeys)
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s: %s", d.RoutingKey, d.Body)
		}
	}()
	<-forever
}
