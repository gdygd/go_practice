package databus

// import (
// 	"log"
// 	"sync"
// )

// type Message struct {
// 	Topic string
// 	Data  interface{}
// }

// type DataBus struct {
// 	subscribers map[string][]chan Message
// 	mu          sync.RWMutex
// }

// func NewDataBus() *DataBus {
// 	return &DataBus{
// 		subscribers: make(map[string][]chan Message),
// 	}
// }

// func (bus *DataBus) Subscribe(topic string) <-chan Message {
// 	ch := make(chan Message, 10) // 버퍼 채널
// 	bus.mu.Lock()
// 	defer bus.mu.Unlock()

// 	bus.subscribers[topic] = append(bus.subscribers[topic], ch)
// 	return ch
// }

// func (bus *DataBus) Publish(msg Message) {
// 	bus.mu.RLock()
// 	defer bus.mu.RUnlock()

// 	for _, ch := range bus.subscribers[msg.Topic] {
// 		select {
// 		case ch <- msg:
// 		default:
// 			log.Printf("Subscriber channel full for topic %s, dropping message", msg.Topic)
// 		}
// 	}
// }

// func (bus *DataBus) replaceSubscribers(topic string) {
// 	newChan := make(chan Message, CHAN_BUF)
// 	bus.subscribers[topic] = []chan Message{newChan}
// 	log.Printf("Replaced subscriber channels for topic %s", topic)
// }

// func clearChannel(chs []chan Message) {
// 	for _, ch := range chs {
// 		for {
// 			select {
// 			case <-ch:
// 			default:
// 				break
// 			}
// 		}
// 	}
// }

/*

bus := NewDataBus()

// REST Server subscribes to "log"
restCh := bus.Subscribe("log")

// Socket Server subscribes to "chat"
socketCh := bus.Subscribe("chat")

// 메시지 발행
bus.Publish(Message{Topic: "log", Data: "user logged in"})
bus.Publish(Message{Topic: "chat", Data: "hello from socket"})

// 처리 루틴 예
go func() {
	for msg := range restCh {
		log.Println("[REST] Received:", msg.Data)
	}
}()

go func() {
	for msg := range socketCh {
		log.Println("[Socket] Received:", msg.Data)
	}
}()


*/
