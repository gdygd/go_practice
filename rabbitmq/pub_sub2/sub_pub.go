package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ch   *amqp.Channel    = nil
	conn *amqp.Connection = nil
	que  amqp.Queue
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func publishing(exchange, routingKey, body string) {
	var count int = 1
	for {
		msg := fmt.Sprintf("%s_%d", body, count)
		count++
		err := ch.Publish(
			exchange,
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			},
		)

		failOnError(err, "failed to publish message..")

		log.Printf(" [x] Sent %s: %s", routingKey, msg)
		time.Sleep(time.Second)
	}
}

func main() {
	routingKey := "info"
	body := "hello"
	var err error = nil

	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err = conn.Channel()
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

	go publishing(exchange, routingKey, body)

	// sub...
	que, err = ch.QueueDeclare("", false, true, true, false, nil)

	err = ch.QueueBind(que.Name, "info", exchange, false, nil)
	failOnError(err, "Faield to bind queue")

	msgs, err := ch.Consume(que.Name, "", true, false, false, false, nil)
	failOnError(err, "Faield to register a consumer")

	log.Printf(" [*] Waiting for logs. Bound to: %v", "info")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s: %s", d.RoutingKey, d.Body)
		}
	}()
	<-forever
}
