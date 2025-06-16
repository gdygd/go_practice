package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	students := []Student{
		Student{"gildong1", 10},
		Student{"gildong2", 20},
		Student{"gildong3", 30},
	}

	// fruits := []string{"aa", "bb", "cc"}

	data, err := json.Marshal(students)
	if err != nil {
		log.Fatalf("JSON marshal error: %v", err)
	}

	err = rdb.Set(ctx, "students", data, time.Hour).Err()
	if err != nil {
		log.Fatalf("Redis Set error: %v", err)
	}

	raw, err := rdb.Get(ctx, "students").Result()
	if err != nil {
		log.Fatalf("Redis Get error: %v", err)
	}

	var result []Student
	err = json.Unmarshal([]byte(raw), &result)
	if err != nil {
		log.Fatalf("JSON unmarshal error: %v", err)
	}

	fmt.Println("Recovered from Redis:", result)
}
