package gapi

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	"grpc_client_test/internal/logger"
	"grpc_client_test/pb"

	"google.golang.org/grpc/connectivity"
)

func (c *GrpcClient) startManageConnect(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	c.txrxwg.Add(1)
	go func() {
		defer c.txrxwg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()
		log.Printf("Starting Manage Connect routine for client")
		c.manageConnect(ctx)
	}()
}

func (c *GrpcClient) startTxRoutine(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	c.txrxwg.Add(1)
	go func() {
		defer c.txrxwg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()
		log.Printf("Starting Tx routine for client")
		c.txRoutine(ctx)
	}()
}

func (c *GrpcClient) startRxRoutine(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	c.txrxwg.Add(1)
	go func() {
		defer c.txrxwg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()
		log.Printf("Starting Rx routine for client %d")
		c.rxRoutine(ctx)
	}()
}

func (c *GrpcClient) manageConnect(ctx context.Context) {
	for {
		logger.Log.Print(2, "manageconnect...#0")
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "[ManageConn] routine exiting..")
			return
		default:
			logger.Log.Print(2, "manageconnect...#1 %v", c.conn.GetState())
			// if c.conn.GetState() != connectivity.Shutdown && c.conn.GetState() != connectivity.TransientFailure {
			// 	time.Sleep(time.Second * 1)
			// 	continue

			// }

			if c.conn.GetState() == connectivity.Shutdown || c.conn.GetState() == connectivity.TransientFailure || c.conn.GetState() == connectivity.Idle {
				logger.Log.Print(2, "manageconnect...#2")

				err1 := c.Connect()
				err2 := c.CreateStream()

				if err1 == nil && err2 == nil {
					logger.Log.Print(2, "Reconnect Successfully")
				} else {
					logger.Log.Error("gRPC Reconnect failed..[%v][%v]", err1, err2)
				}
			}

			time.Sleep(time.Second * 1)
		}
	}
}

func (c *GrpcClient) txRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "quit tx routine...")
			if c.stream == nil {
				logger.Log.Print(2, "#1 tx, stream is nil..")
				return
			}

			// 남은 메세지 처리

			// 종료
			c.mu.Lock()
			c.stream.CloseSend()
			c.mu.Unlock()

			return
		default:
			logger.Log.Print(2, "txRoutine... %v", c.conn.GetState())
			if c.conn.GetState() != connectivity.Ready {
				time.Sleep(time.Second * 1)
				continue
			}

			if c.stream == nil {
				logger.Log.Print(2, "#2 tx, stream is nil..")
				continue
			}

			// send.
			c.mu.Lock()
			err := c.stream.Send(&pb.Hello{Msg: "to Client..."})
			c.mu.Unlock()

			if err != nil {
				log.Printf("Send error: %v", err)
				return
			}

			time.Sleep(time.Second * 1)
		}
	}
}

func (c *GrpcClient) rxRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "quit rx routine...")
			return
		default:
			if c.conn.GetState() != connectivity.Ready {
				continue
			}

			c.mu.Lock()
			logger.Log.Print(2, "rxRoutine... %v", c.conn.GetState())
			c.mu.Unlock()

			if c.conn.GetState() != connectivity.Ready {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			if c.stream == nil {
				logger.Log.Print(1, "rx, stream is nil..")
				continue
			}

			// recv
			c.mu.Lock()
			resp, err := c.stream.Recv()
			c.mu.Unlock()

			if err == io.EOF {
				logger.Log.Error("Server closed stream")
				return
			}
			if err != nil {
				logger.Log.Error("Recv error: %v", err)
				return
			}
			if err == nil {
				log.Printf("From server: %s", resp.Msg)
			}
		}
	}
}
