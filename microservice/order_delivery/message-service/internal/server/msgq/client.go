package msgq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"message_service/internal/container"
	"message_service/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	wg       *sync.WaitGroup // app에서 관리하는 wg
	txrxwg   sync.WaitGroup  // pub / sub routine group
	conn     *amqp.Connection
	msg_ch   *amqp.Channel
	que      amqp.Queue
	exchange string

	ctx    context.Context // master context
	cancel context.CancelFunc
	mu     sync.RWMutex
	ct     *container.Container

	close_ch  chan *amqp.Error
	connected bool

	routingKey string
	bindingKey []string
}

func NewClient(wg *sync.WaitGroup, ct *container.Container) (*RabbitMQClient, error) {
	ctx, cancel := context.WithCancel(context.Background())
	mqclient := &RabbitMQClient{
		wg:         wg,
		ctx:        ctx,
		cancel:     cancel,
		ct:         ct,
		connected:  false,
		exchange:   "topic_logs",
		routingKey: "info",
		bindingKey: []string{"info", "warn", "critical"},
	}

	err := mqclient.Connect() // Dial, and open channel
	if err != nil {
		logger.Log.Error("Rabbit server not connected...")
	} else {
		mqclient.ExchangeDeclareTopic("topic_logs")
		err = mqclient.QueueDeclare()
		if err != nil {
			logger.Log.Error("%v", err)
		}
	}

	return mqclient, err
}

func (q *RabbitMQClient) Connect() error {
	addr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(addr)
	if err != nil {
		q.connected = false
		return fmt.Errorf("failed to connect to RabbitMQ")
	}

	q.close_ch = conn.NotifyClose(make(chan *amqp.Error))

	ch, err := conn.Channel()
	if err != nil {
		q.connected = false
		return fmt.Errorf("failed to open a channel")
	}

	q.conn = conn
	q.msg_ch = ch
	q.connected = true
	return nil
}

func (q *RabbitMQClient) ExchangeDeclareTopic(exchange string) error {
	// ex) exchange := "topic_logs"

	err := q.msg_ch.ExchangeDeclare(
		exchange,
		"topic",
		true,  // durable: 운영 시 true 권장
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to declare exchange")
	}

	return nil
}

func (q *RabbitMQClient) QueueDeclare() error {
	var err error = nil
	q.que, err = q.msg_ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return fmt.Errorf("Failed to declare a queue (%v)", err)
	}

	for _, key := range q.bindingKey {
		err = q.msg_ch.QueueBind(q.que.Name, key, q.exchange, false, nil)
		if err != nil {
			logger.Log.Error("Failed to bind queue..(%v)", err)
		}
	}

	return err
}

func (m *RabbitMQClient) isConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.connected
}

func (m *RabbitMQClient) Start() {
	logger.Log.Print(2, "RabbitMQ client start..")
	defer m.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Print(2, "RabbitMQ start panic..%v ", r)
		}
	}()

	//

	subpubCtx, subpubCancel := context.WithCancel(m.ctx) // master ctx로 하위 루틴 컨텍스트 생성

	var mu sync.Mutex
	var subRunning, pubRunning, connRunning bool
	subDone := make(chan struct{}, 1)
	pubDone := make(chan struct{}, 1)
	connectDone := make(chan struct{}, 1)

	// 최초 시작
	m.startSubRoutine(subpubCtx, &mu, &subRunning, subDone)
	m.startPubRoutine(subpubCtx, &mu, &pubRunning, pubDone)
	m.startManageConnect(subpubCtx, &mu, &connRunning, connectDone)

	for {
		select {
		case <-m.ctx.Done():
			logger.Log.Print(2, "mp client master context canceled..")
			subpubCancel()
			goto WAIT

		case <-subDone:
			logger.Log.Print(2, "sub routine ended..")
			if m.ctx.Err() == nil {
				m.startSubRoutine(subpubCtx, &mu, &subRunning, subDone)
			}

		case <-pubDone:
			logger.Log.Print(2, "pub routine ended..")
			if m.ctx.Err() == nil {
				m.startPubRoutine(subpubCtx, &mu, &pubRunning, pubDone)
			}

		case <-connectDone:
			logger.Log.Print(2, "conn routine ended..")
			if m.ctx.Err() == nil {
				m.startManageConnect(subpubCtx, &mu, &connRunning, connectDone)
			}
		default:
			time.Sleep(time.Second * 1)
		}
	}

WAIT:
	done := make(chan struct{})
	go func() {
		for {
			mu.Lock()
			if !subRunning && !pubRunning && !connRunning {
				mu.Unlock()
				break
			}
			mu.Unlock()
			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case <-done:
		logger.Log.Print(2, "Graceful shutdown!")
	case <-time.After(time.Second * 5):
		logger.Log.Print(2, "Gracefule shutdown.. timeout..")
	}

	logger.Log.Print(2, "Quit rabbitmq client for pub..")
}

func (m *RabbitMQClient) Shutdown() {
	logger.Log.Print(2, "Shutdown RabbitMQ client")

	m.cancel()
	m.txrxwg.Wait()

	logger.Log.Print(2, "rabbitmq client for pub finished")
}

func (m *RabbitMQClient) close() {
	m.mu.Lock()
	m.connected = false
	m.mu.Unlock()

	if m.msg_ch != nil {
		m.msg_ch.Close()
	}

	if m.conn != nil {
		m.conn.Close()
	}

	m.msg_ch = nil
	m.conn = nil
}
