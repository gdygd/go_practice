package databus

import (
	"log"
	"sync"
)

const (
	CHAN_BUF  = 10
	RETRY_BUF = 10
)

type Message struct {
	Topic string
	Data  interface{}
}

type DataBus struct {
	subscribers map[string][]chan Message
	retryQueue  chan Message
	mu          sync.RWMutex
	wg          sync.WaitGroup
	closed      chan struct{}
}

func NewDataBus() *DataBus {
	bus := &DataBus{
		subscribers: make(map[string][]chan Message),
		retryQueue:  make(chan Message, RETRY_BUF),
		closed:      make(chan struct{}),
	}
	bus.wg.Add(1)
	go bus.retryDispatcher()
	return bus
}

func (bus *DataBus) retryDispatcher() {
	bus.wg.Done()

	for {
		select {
		case msg, ok := <-bus.retryQueue:
			if !ok {
				log.Println("Retry dispatcher shutting down.")
				return
			}
			bus.mu.RLock()
			for _, ch := range bus.subscribers[msg.Topic] {
				select {
				case ch <- msg:
				default:
					log.Printf("Retry failed to send topic %s", msg.Topic)
				}
			}
			bus.mu.RUnlock()
		case <-bus.closed:
			log.Println("Received shutdown signal.")
			return
		}
	}
}

func (bus *DataBus) Subscribe(topic string) <-chan Message {
	ch := make(chan Message, CHAN_BUF)
	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.subscribers[topic] = append(bus.subscribers[topic], ch)
	return ch
}

func (bus *DataBus) Unsubscribe(topic string, targetCh <-chan Message) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	channels := bus.subscribers[topic]
	newList := channels[:0]

	for _, ch := range channels {
		if ch != targetCh {
			newList = append(newList, ch)
		} else {
			close(ch)
		}
	}

	if len(newList) > 0 {
		bus.subscribers[topic] = newList
	} else {
		delete(bus.subscribers, topic)
	}
}

func (bus *DataBus) Publish(msg Message) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	for _, ch := range bus.subscribers[msg.Topic] {
		select {
		case ch <- msg:
		default:
			select {
			case bus.retryQueue <- msg:
				// ch full.. enqueue retry Q
				log.Printf("Enqueue retry topic %s\n", msg.Topic)
			default:
				log.Printf("Retry queue full for topic %s, dropping message", msg.Topic)
			}
		}
	}
}

func (bus *DataBus) ShutDown() {
	log.Println("==============================================Shutting down Databus...==============================================")
	close(bus.retryQueue)
	close(bus.closed) // retryDispatcher 종료

	bus.wg.Wait()

	bus.mu.Lock()
	for topic, chans := range bus.subscribers {
		for _, ch := range chans {
			close(ch)
		}
		delete(bus.subscribers, topic)
	}
	bus.mu.Unlock()

	log.Println("==============================================Shut down Databus!==============================================")
}

func (bus *DataBus) clearRetryQueue() {
	bus.retryQueue = make(chan Message, 1000)
}

func (bus *DataBus) replaceSubscribers(topic string) {
	newChan := make(chan Message, CHAN_BUF)
	bus.subscribers[topic] = []chan Message{newChan}
	log.Printf("Replaced subscriber channels for topic %s", topic)
}

func clearChannel(chs []chan Message) {
	for _, ch := range chs {
		for {
			select {
			case <-ch:
			default:
				break
			}
		}
	}
}
