package main

import (
	"databus/databus"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var bus *databus.DataBus

func App1(wg *sync.WaitGroup) {
	defer wg.Done()
	// time.Sleep(time.Second * 2)

	app1Ch := bus.Subscribe("app1")

	for {
		select {
		case msg := <-app1Ch:
			log.Println("[APP1] Received:", msg.Data)
		default:

		}
		time.Sleep(time.Millisecond * 100)
	}

}

func App2(wg *sync.WaitGroup) {
	defer wg.Done()

	app2Ch := bus.Subscribe("app2")

	for {
		select {
		case msg := <-app2Ch:
			log.Println("[APP2] Received:", msg.Data)
		default:

		}
		time.Sleep(time.Millisecond * 100)
	}

}

func App3(wg *sync.WaitGroup) {
	defer wg.Done()

	app3Ch := bus.Subscribe("app3")

	for {
		select {
		case msg := <-app3Ch:
			log.Println("[APP3] Received:", msg.Data)
		default:

		}
		time.Sleep(time.Millisecond * 100)
	}

}

func Producer(wg *sync.WaitGroup) {
	defer wg.Done()

	var arrchannel []string = []string{"app1", "app2", "app3"}
	for i := 0; i < 100; i++ {
		idx := rand.Intn(3)
		chann := arrchannel[idx]
		data := fmt.Sprintf("%s, data (%d)", chann, i)

		log.Printf("Publish app:[%s] data:[%s]", chann, data)

		msg := databus.Message{
			Topic: chann,
			Data:  data,
		}
		bus.Publish(msg)

		time.Sleep(time.Millisecond * 200)
	}

}

func main() {
	log.Println("start program..")

	//----------------------------
	bus = databus.NewDataBus()

	//----------------------------

	var wg sync.WaitGroup
	wg.Add(8)

	go App1(&wg)
	go App2(&wg)
	go App3(&wg)
	go Producer(&wg)
	go Producer(&wg)
	go Producer(&wg)
	go Producer(&wg)
	go Producer(&wg)

	wg.Wait()

	log.Println("exit program..#1")
	ch_signal := make(chan os.Signal, 10)
	signal.Notify(ch_signal, syscall.SIGINT)
	<-ch_signal

	log.Println("exit program..#2")

}
