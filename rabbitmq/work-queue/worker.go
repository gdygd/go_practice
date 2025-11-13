package main

import (
	"bytes"
	"log"
	"time"

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

	q, err := ch.QueueDeclare(
		"task_queue", // queue name
		true,         // durable: 서버 재시작 후에도 유지
		false,        // autoDelete
		false,        // exclusive
		false,        // noWait
		nil,          // arguments
	)

	failOnError(err, "Failed to declare a queue")

	// 브로커가 보내는 메세지 제어 consumer가 동시에 처리 중인 메세지를 1개로 제한ㅛ
	err = ch.Qos(1, 0, false)
	failOnError(err, "Failed to register a consumer")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false, // autoAck=false : 명시적으로 Ack 필요
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			time.Sleep(time.Duration(dotCount) * time.Second)

			log.Printf("Done")
			d.Ack(false) // 작업 완료 후 ack
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
