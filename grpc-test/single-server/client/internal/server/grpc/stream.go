package gapi

import (
	"fmt"
	"io"
	"log"

	"grpc_svr_test/pb"
)

func (s *Server) ConnMessage(stream pb.HelloService_ConnMessageServer) error {
	// 채널 생성
	clientMsgs := make(chan *pb.Hello)

	// recv
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
			log.Printf("From client: %s", req.Msg)
			clientMsgs <- req
		}
	}()

	// send
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
