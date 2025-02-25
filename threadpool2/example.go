package main

import (
	"fmt"
	"time"

	//"github.com/shettyh/threadpool"
	"threadpool2/thpool"
)

func main() {
	pool := thpool.NewThreadPool(2000, 100000)

	task := &myTask{ID: 123}
	task1 := &myTask{ID: 1231}
	task2 := &myTask{ID: 1232}
	task3 := &myTask{ID: 1233}
	pool.Execute(task)
	pool.Execute(task1)
	pool.Execute(task2)
	pool.Execute(task3)
	time.Sleep(10 * time.Second)
}

type myTask struct {
	ID int64
}

func (m *myTask) Run() {
	fmt.Println("Running my task ", m.ID)
}
