package main

import (
	"fmt"
	"net"
	"time"
)

// func main() {

// 	fmt.Printf("client... \n")

// 	// Connect to the server
// 	conn, err := net.Dial("tcp", "localhost:23000")
// 	if err != nil {
// 		fmt.Println("Error connecting:", err.Error())
// 		return
// 	}
// 	defer conn.Close()

// 	for {

// 		// Receive response from the server
// 		buffer := make([]byte, 1024)
// 		n, err := conn.Read(buffer)
// 		if err != nil {
// 			fmt.Println("Error receiving data:", err.Error())
// 			continue
// 		}
// 		// Process the response
// 		response := string(buffer[:n])
// 		fmt.Println("Response from server:", response)

// 		//---------------------------------------------

// 		// Send data to the server
// 		message := []byte("Hello, server!")
// 		_, err = conn.Write(message)
// 		if err != nil {
// 			fmt.Println("Error sending data:", err.Error())
// 			return
// 		}

// 		time.Sleep(time.Millisecond * 1000)
// 	}
// }

func main() {

	conn, err := net.Dial("tcp", "192.168.56.2:20030")
	if err != nil {
		panic(err)
	}

	// Set the read deadline to 5 seconds
	//err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	//err = conn.SetDeadline(time.Now().Add(1 * time.Second))
	err = conn.SetDeadline(time.Now().Add(2000 * time.Millisecond))
	if err != nil {
		panic(err)
	}

	// Perform a read operation on the connection
	for {

		fmt.Printf("tcp client...(1)\n")
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		fmt.Printf("tcp client...(2)\n")
		if err != nil {
			panic(err)
		}
		if n > 1 {
			fmt.Printf("Received %d bytes: %s\n", n, string(buf[:n]))
		} else {
			fmt.Printf("no recv data...\n")
		}

		time.Sleep(time.Millisecond * 1000)
	}

}
