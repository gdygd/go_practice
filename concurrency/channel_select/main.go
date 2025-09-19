package main

import (
	"fmt"
	"time"
)

func server1(ch chan string) {
	for {
		ch <- fmt.Sprintf("server1")
		time.Sleep(time.Second * 3)
	}

}

func server2(ch chan string) {
	for {
		ch <- fmt.Sprintf("server2")
		time.Sleep(time.Second * 3)
	}
}

func main() {

	var ch1 = make(chan string)
	var ch2 = make(chan string)

	go server1(ch1)
	go server2(ch2)

	for {
		select {
		case s1 := <-ch1:
			fmt.Println("s1 : ", s1)
		case s2 := <-ch1:
			fmt.Println("s2 : ", s2)
		case s3 := <-ch2:
			fmt.Println("s3 : ", s3)
		case s4 := <-ch2:
			fmt.Println("s4 : ", s4)
		}
	}
}
