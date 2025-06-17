package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // 필요 시 변경
	})

	channels := []string{"channel1", "channel2", "channel3"}

	// context로 graceful shutdown 구현
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go subscribeMultipleChannels(ctx, rdb, channels)

	time.Sleep(1 * time.Second) // 구독 준비 시간

	for i := 0; i < 5; i++ {

		for chidx := 0; chidx < len(channels); chidx++ {
			channel := channels[chidx]
			msg := fmt.Sprintf("(%s) message %d", channel, i)

			err := rdb.Publish(ctx, channel, msg).Err()
			if err != nil {
				log.Fatalf("Publish error: %v", err)
			}
		}

		time.Sleep(1 * time.Second)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down...")
}

func subscribeMultipleChannels(ctx context.Context, rdb *redis.Client, channels []string) {
	pubsub := rdb.Subscribe(ctx, channels...)
	defer pubsub.Close()

	log.Println("Subscribed to channels:", strings.Join(channels, ", "))

	// 메시지 처리 루프
	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			log.Println("Context canceled. Exiting subscription loop.")
			return
		case msg := <-ch:
			fmt.Printf("[Channel: %s] Message: %s\n", msg.Channel, msg.Payload)
		}
	}
}
