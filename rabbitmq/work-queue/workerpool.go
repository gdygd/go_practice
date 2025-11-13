package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// func worker(id int, wg *sync.WaitGroup, ctx context.Context) {
// 	defer wg.Done()

// 	select {
// 	case <-time.After(time.Second * 20):
// 		log.Printf("timeout .. %d", id)
// 		return
// 	case <-ctx.Done():
// 		log.Printf("worker ended .. %d", id)
// 		return
// 	}
// }

func worker(id int, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	log.Printf("Start worker %d", id)

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

	for {
		select {
		case d := <-msgs:
			log.Printf("Received a message: (%d) %s", id, d.Body)
			// dotCount := bytes.Count(d.Body, []byte("."))
			dotCount := 5
			time.Sleep(time.Duration(dotCount) * time.Second)

			log.Printf("(%d)Done", id)
			d.Ack(false) // 작업 완료 후 ack
		case <-ctx.Done():
			log.Printf("worker ended .. %d", id)
			return
			// default:
			// 	time.Sleep(time.Millisecond * 100)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	var workerCnt int = 10
	work_ctx, work_cancel := context.WithCancel(context.Background())

	ch_signal := make(chan os.Signal, 10)
	signal.Notify(ch_signal, syscall.SIGINT)

	go func() {
		select {
		case <-ch_signal:
			work_cancel()
		}
	}()

	for i := 0; i < workerCnt; i++ {
		wg.Add(1)
		go worker(i+1, &wg, work_ctx)
	}

	wg.Wait()
}
