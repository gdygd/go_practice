package msgq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"message_service/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (m *RabbitMQClient) startManageConnect(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	m.txrxwg.Add(1)

	go func() {
		defer m.txrxwg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()

		m.manageConnect(ctx)
	}()
}

func (m *RabbitMQClient) manageConnect(ctx context.Context) {
	logger.Log.Print(2, "manageconnect..")
	for {
		select {
		case <-m.close_ch:
			logger.Log.Print(2, "conn close.")
			m.close()
			m.Connect()
		case <-ctx.Done():
			logger.Log.Print(2, "quit routine [manageconnect]..")
			m.close() // channel & conn close
			return
		}
	}
}

func (m *RabbitMQClient) startSubRoutine(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	m.txrxwg.Add(1)

	go func() {
		defer m.txrxwg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()

		m.subRoutine(ctx)
	}()
}

func (m *RabbitMQClient) subRoutine(ctx context.Context) {
	logger.Log.Print(2, "subRoutine...")
	msgs, err := m.msg_ch.Consume(m.que.Name, "", true, false, false, false, nil)
	if err != nil {
		logger.Log.Print(2, "Failed to register a consumer")
		m.close() // channel & conn close
		return
	}
	for {
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "quit routine [subscribe]..")
			m.close() // channel & conn close
			return
		case d := <-msgs:
			logger.Log.Print(2, "subscribe message [%s][%s]", d.RoutingKey, d.Body)

		default:
			logger.Log.Print(2, "subscribe...")
			time.Sleep(time.Second * 1)
		}
	}
}

func (m *RabbitMQClient) startPubRoutine(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	m.txrxwg.Add(1)

	go func() {
		defer m.txrxwg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()

		m.pubRoutine(ctx)
	}()
}

func (m *RabbitMQClient) pubRoutine(ctx context.Context) {
	logger.Log.Print(2, "pubRoutine...")
	var count int = 1
	for {
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "quit routine [publish]..")
			m.close() // channel & conn close
			return
		default:
			logger.Log.Print(2, "pub...")

			if !m.isConnected() {
				logger.Log.Print(2, "mq server not connected..")
				time.Sleep(time.Second * 1)
				continue
			}

			msg := fmt.Sprintf("hello_%d", count)
			count++

			// publishing 함수로 뺼것
			err := m.msg_ch.Publish(
				m.exchange,
				m.routingKey,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(msg),
				},
			)
			if err != nil {
				logger.Log.Error("message publish failed..")
				m.close() // channel & conn close
				return
			}
			logger.Log.Error("[x] Sent %s: %s", m.routingKey, msg)

			time.Sleep(time.Second * 1)
		}
	}
}
