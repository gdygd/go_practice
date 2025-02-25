package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9100")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024*1024) // 1 MB 버퍼

	start := time.Now()

	for i := 0; i < 100; i++ { // 100 MB 전송
		_, err = conn.Write(buffer)
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("eltm:%v\n", elapsed)

	totalBytes := 100 * 1024 * 1024 // 100 MB
	bandwidth := float64(totalBytes*8) / (elapsed * 1e6)
	bandwidth1 := float64(totalBytes*8) / (elapsed * 1000000)

	fmt.Printf("Bandwidth: %.2f Mbps\n", bandwidth)
	fmt.Printf("Bandwidth1: %.2f Mbps\n", bandwidth1)

	bandwidth2 := float64(totalBytes*8) / (elapsed)

	fmt.Printf("Bandwidth2: %.2f Mbps\n", bandwidth2)
}
