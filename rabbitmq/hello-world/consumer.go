package main

import (
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	ID      int
	Name    string
	Message string
	Time    time.Time
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Faield to open a channel : %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queeue : %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer tag
		true,  // auto-ack (자동으로 Ack 처리)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer : %v", err)
	}

	log.Println("[*] Waiting for message. To exit press CTRL+C")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			// log.Printf("[*] Received: %s", d.Body)
			var data Message
			if err := json.Unmarshal(d.Body, &data); err != nil {
				log.Printf("failed to parse json.. %v", err)
				continue
			}
			log.Printf("[*] Received: %v", data)
		}
	}()

	<-forever
}
