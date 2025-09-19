package main

import (
	"fmt"
	"time"
)

func listenToChannel(ch chan int) {
	for {
		i := <-ch
		fmt.Println("Got i", i, "from channel")
		time.Sleep(time.Second * 1)
	}
}

func main() {

	var ch = make(chan int, 200)
	go listenToChannel(ch)

	for i := 0; i <= 100; i++ {
		fmt.Println("sending.. ", i, "to channel..")
		ch <- i
		fmt.Println("sent", i, "to channel")
	}

	fmt.Println("Done")
	close(ch)

}
