package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func subscribe(rdb *redis.Client, channel string) {
	sub := rdb.Subscribe(ctx, channel)

	// 최초 응답에서 Subscription 확인을 위해 수신
	_, err := sub.Receive(ctx)
	if err != nil {
		log.Fatalf("Subscribe error: %v", err)
	}

	ch := sub.Channel()

	for msg := range ch {
		fmt.Printf("Received message from %s: %s\n", msg.Channel, msg.Payload)
	}

	// 종료 시 unsubscribe도 가능
	sub.Close()

}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	go subscribe(rdb, "my-channel")

	time.Sleep(1 * time.Second) // 구독 준비 시간
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("message %d", i)
		err := rdb.Publish(ctx, "my-channel", msg).Err()
		if err != nil {
			log.Fatalf("Publish error: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

	time.Sleep(3 * time.Second)
}
