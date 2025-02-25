package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":9100")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 9100...")
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error accepting:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024*1024) // 1 MB 버퍼
	var totalBytes int64

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		totalBytes += int64(n)
	}

	fmt.Printf("Total bytes received: %d\n", totalBytes)
}
