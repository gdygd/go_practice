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
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "[ManageConn] routine exiting..")
			return
		default:
			if c.conn.GetState() != connectivity.Shutdown {
				time.Sleep(time.Second * 1)
				continue

			}

			if c.conn.GetState() == connectivity.Shutdown {
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
			return
		default:
			// logger.Log.Print(2, "txRoutine...")
			if c.conn.GetState() != connectivity.Ready {
				continue
			}

			// send.
			err := c.stream.Send(&pb.Hello{Msg: "to Client..."})
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

			// logger.Log.Print(2, "rxRoutine...")

			// recv
			resp, err := c.stream.Recv()
			if err == io.EOF {
				log.Println("Server closed stream")
				return
			}
			if err != nil {
				log.Fatalf("Recv error: %v", err)
				return
			}
			if err == nil {
				log.Printf("From server: %s", resp.Msg)
			}
		}
	}
}
