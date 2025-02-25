package main

import (
	"fmt"
	"time"
)

func main() {
	requestChan := make(chan chan string)

	go goroutineC(requestChan)
	go goroutineD(requestChan)

	time.Sleep(time.Second * 2)
}

func goroutineC(requestChan chan chan string) {
	responseChan := make(chan string)
	requestChan <- responseChan
	//response := <-responseChan

	var strRes string

	var ok bool = false
	for {
		select {
		case res := <-responseChan:
			strRes = res
			fmt.Printf("Response : %v\n", res)
			ok = true
			break
		default:
			if !ok {
				fmt.Printf("default...")
			}
		}
	}

	fmt.Printf("Response : %v\n", strRes)
}

func goroutineD(requestChan chan chan string) {
	responseChan := <-requestChan
	time.Sleep(time.Second * 1)
	responseChan <- "wassup!"
}
