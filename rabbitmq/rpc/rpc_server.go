package main

import (
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func fib(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Faield to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // queue name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // args
	)

	failOnError(err, "Failed to declare rpc queue")

	err = ch.Qos(1, 0, false)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer
		false, // autoAck=false (명시적 Ack 필요)
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")
	log.Println(" [x] Awaiting RPC requests")

	for d := range msgs {
		n, err := strconv.Atoi(string(d.Body))
		if err != nil {
			log.Printf("Invalid arg : %s", d.Body)

			_ = ch.Publish(
				"",        // default exchange
				d.ReplyTo, // reply queue
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte("error: invalid arg"),
				},
			)
			d.Ack(false)
			continue
		}

		log.Printf(" [.] fib(%d)", n)

		result := fib(n)

		err = ch.Publish(
			"",        // default exchange (direct)
			d.ReplyTo, // reply-to queue
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          []byte(strconv.Itoa(result)),
			},
		)
		if err != nil {
			log.Printf("Faield to publish response: %v", err)
		}
		d.Ack(false)

		time.Sleep(50 * time.Millisecond)
	}
}
