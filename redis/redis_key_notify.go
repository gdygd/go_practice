package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	key := "mykey"
	eventChannel := fmt.Sprintf("__keyspace@0__:%s", key)

	// 1. 키 변경 이벤트 구독
	sub := rdb.Subscribe(ctx, eventChannel)
	defer sub.Close()

	fmt.Println("Subscribed to", eventChannel)

	ch := sub.Channel()

	// 2. 메시지 수신 루프
	go func() {
		for msg := range ch {
			fmt.Printf("KEY EVENT: [%s] %s\n", msg.Channel, msg.Payload)
		}
	}()

	// 3. 예시로 키 변경
	rdb.Set(ctx, key, "hello", 0)
	rdb.Del(ctx, key)

	select {} // block forever
}
