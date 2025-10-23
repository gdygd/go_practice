package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"grpc_test/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	// 클라이언트 → 서버 전송 고루틴
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("💬 Input message: ")
			scanner.Scan()
			text := scanner.Text()

			if text == "exit" {
				stream.CloseSend()
				return
			}

			err := stream.Send(&pb.Hello{Msg: text})
			if err != nil {
				log.Printf("Send error: %v", err)
				return
			}
		}
	}()

	// 서버 → 클라이언트 수신 루프
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("Server closed stream")
			break
		}
		if err != nil {
			log.Fatalf("Recv error: %v", err)
		}

		log.Printf("📨 From server: %s", resp.Msg)

		// 서버가 "Hello Server" 요청을 보냈을 경우 클라이언트가 응답
		if resp.Msg == "Hello Server" {
			time.Sleep(500 * time.Millisecond)
			if err := stream.Send(&pb.Hello{Msg: "Hello Server Response"}); err != nil {
				log.Printf("Response send error: %v", err)
			}
		}
	}
}
