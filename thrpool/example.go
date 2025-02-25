package main

import (
	"fmt"
	"thrpool/pool"
	"time"
)

// Id, CloseHandle은 필수 변수
type MyTask struct {
	Id          int64
	CloseHandle chan bool
}

func main() {
	pool := pool.NewThreadPool(2, 2)
	task := &MyTask{Id: 111, CloseHandle: make(chan bool)}
	task2 := &MyTask{Id: 211, CloseHandle: make(chan bool)}
	task3 := &MyTask{Id: 311, CloseHandle: make(chan bool)}
	err := pool.Execute(task)
	if err != nil {
		fmt.Printf("task err : %v", err)
	}
	err = pool.Execute(task2)
	if err != nil {
		fmt.Printf("task2 err : %v", err)
	}
	err = pool.Execute(task3)
	if err != nil {
		fmt.Printf("task3 err : %v", err)
	}

	time.Sleep(time.Second * 10)
}

func (t *MyTask) Run() {
	fmt.Printf("MyTask Run. %d \n", t.Id)

	// return 시 End job..처리
}
