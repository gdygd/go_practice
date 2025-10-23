package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"grpc_test/pb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

// Chat(grpc.BidiStreamingServer[Hello, Hello]) error
// func (s *server) Chat(stream pb.HelloService_ChatServer) error {
// 	log.Println("✅ Client connected")

// 	for {
// 		req, err := stream.Recv()
// 		if err == io.EOF {
// 			log.Println("Client stream closed")
// 			return nil
// 		}
// 		if err != nil {
// 			log.Printf("Receive error: %v", err)
// 			return err
// 		}

// 		log.Printf("📩 From client: %s", req.Msg)

// 		// 클라이언트가 보낸 메시지에 응답
// 		resp := &pb.Hello{Msg: "Hello Client"}
// 		if err := stream.Send(resp); err != nil {
// 			log.Printf("Send error: %v", err)
// 			return err
// 		}

// 		// 서버가 클라이언트에게 요청(역방향 호출)
// 		serverReq := &pb.Hello{Msg: "Hello Server"}
// 		if err := stream.Send(serverReq); err != nil {
// 			log.Printf("Send error: %v", err)
// 			return err
// 		}
// 	}
// }

func (s *server) Chat(stream pb.HelloService_ChatServer) error {
	log.Println("✅ Client connected")

	// 채널 생성
	clientMsgs := make(chan *pb.Hello)

	// 1️⃣ 수신 goroutine
	go func() {
		defer close(clientMsgs)
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				log.Println("client closed stream")
				return
			}
			if err != nil {
				log.Printf("receive error: %v", err)
				return
			}
			log.Printf("📩 From client: %s", req.Msg)
			clientMsgs <- req
		}
	}()

	// 2️⃣ 송신 루프 (메인 goroutine)
	for {
		select {
		case msg, ok := <-clientMsgs:
			if !ok {
				log.Println("client message channel closed")
				return nil
			}
			// 클라이언트가 보낸 메시지에 응답
			resp := &pb.Hello{Msg: fmt.Sprintf("Hello Client, you said: %s", msg.Msg)}
			if err := stream.Send(resp); err != nil {
				log.Printf("send error: %v", err)
				return err
			}

			// 서버가 클라이언트에 별도 메시지 push (역방향)
			serverPush := &pb.Hello{Msg: "Hello Server"}
			if err := stream.Send(serverPush); err != nil {
				log.Printf("send error: %v", err)
				return err
			}
		case <-stream.Context().Done():
			log.Println("client disconnected")
			return nil
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})

	fmt.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	ch_signal := make(chan os.Signal, 2)
	signal.Notify(ch_signal, syscall.SIGHUP, syscall.SIGINT)
	<-ch_signal

	s.GracefulStop()

	time.Sleep(time.Second * 5)
}

// func main() {
// 	fmt.Printf("Hello\n")
// 	var h pb.Hello
// 	h.Msg = "123"
// 	fmt.Printf("H:%s", h.Msg)
// }
