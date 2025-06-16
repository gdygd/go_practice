package tcpserver

import (
	"context"
	"log"
	"sync"
	"time"
)

func (a *SStServer) clearSstServerEnv() {
	log.Printf("Clearing SST SERVER env ID: %d", a.Id)

	if a.Tcp != nil {
		if err := a.Tcp.Close(); err != nil {
			log.Printf("Error closing TcpHandler for ID %d: %v", a.Id, err)
		} else {
			log.Printf("Closed TcpHandler for ID %d (Addr: %s)", a.Id, a.Tcp.GetAddress())
		}
	} else {
		log.Printf("TcpHandler is nil for ID %d", a.Id)
	}

	a.Id = 0
	a.Tcp = nil
}

func (a *SStServer) Shutdown() {
	log.Printf("SST server Shutdown ..")
	close(a.txQueue)
	a.cancel()      // 종료 시그널
	a.txrxWg.Wait() // Tx,Rx 종료 대기
	log.Printf("Tx/Rx routines finished  ")
	<-a.ctx.Done()
	log.Printf("SST server Shutdown finished!")
}

func (a *SStServer) clientTxRoutine(ctx context.Context) {
	log.Printf("tx routine #1")
	for {
		log.Printf("tx routine #2")
		select {
		case <-ctx.Done():
			log.Printf("tx routine #3")
			for msg := range a.txQueue {
				a.Tcp.Send(msg) // 최종 처리
			}
			log.Printf("[Tx] Exiting after queue : client %d", a.Id)

			return
		default:
			log.Printf("tx routine #4")
			if a.Tcp.IsConnected() {
				a.ManageTx() // txQueue에서 값을 가져온 후 처리
			}
			time.Sleep(time.Millisecond * 200)
		}
	}
}

func (a *SStServer) clientRxRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("[Rx] Exiting for client %d", a.Id)
			return
		default:
			if a.Tcp.IsConnected() {
				a.ManageRx() // rx메세지 처리
			}
			time.Sleep(time.Millisecond * 200)
		}
	}
}

func (a *SStServer) manageConnect(ctx context.Context) {
	retryCount := 0

	for {
		select {
		case <-ctx.Done():
			log.Printf("[Conn] Reconnect routine exiting: client %d", a.Id)
			return
		default:
			if a.Tcp.IsConnected() {
				time.Sleep(1 * time.Second)
				continue
			}

			ok, err := a.Tcp.Connect()
			if err != nil {
				retryCount++
				if retryCount%3 == 0 {
					log.Printf("[Conn-%d] Reconnect failed (%d attempts): %v", a.Id, retryCount, err)
				}
				time.Sleep(3 * time.Second)
				continue
			}

			if ok {
				log.Printf("[Conn-%d] Reconnected successfully.", a.Id)
				retryCount = 0
			}
		}
	}
}

func (a *SStServer) StartTxRoutine(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	a.txrxWg.Add(1)
	go func() {
		defer a.txrxWg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()
		log.Printf("Starting Tx routine for client %d", a.Id)
		a.clientTxRoutine(ctx)
	}()
}

func (a *SStServer) StartRxRoutine(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	a.txrxWg.Add(1)
	go func() {
		defer a.txrxWg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()
		log.Printf("Starting Rx routine for client %d", a.Id)
		a.clientRxRoutine(ctx)
	}()
}

func (a *SStServer) StartManageConnect(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	a.txrxWg.Add(1)
	go func() {
		defer a.txrxWg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()
		log.Printf("Starting Manage Connect routine for client %d", a.Id)
		a.manageConnect(ctx)
	}()
}

func (a *SStServer) close() {
	if err := a.Tcp.Close(); err != nil {
		log.Printf("close error for client %d: %v", a.Id, err)
	}
}

func (a *SStServer) Start() {
	defer a.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("THRclient panic: %v \n", r)
		}
	}()

	defer a.clearSstServerEnv()

	txrxCtx, txrxCancel := context.WithCancel(a.ctx)
	// defer txrxCancel()	// 아래에서 호출함.

	var mu sync.Mutex
	var txRunning, rxRunning, connRunning bool // 실행상태 체크
	txDone := make(chan struct{}, 1)
	rxDone := make(chan struct{}, 1)
	connectDone := make(chan struct{}, 1)

	// 최초 시작
	a.StartTxRoutine(txrxCtx, &mu, &txRunning, txDone)
	a.StartRxRoutine(txrxCtx, &mu, &rxRunning, rxDone)
	a.StartManageConnect(txrxCtx, &mu, &connRunning, connectDone)

	// check routine
	for {
		select {
		case <-a.ctx.Done():
			log.Printf("THRclient context canceled for client %d", a.Id)
			txrxCancel() // TxRx 루틴 종료 요청
			goto WAIT
		case <-txDone:
			log.Printf("Tx routine ended for client %d", a.Id)
			if a.ctx.Err() == nil {
				a.StartTxRoutine(txrxCtx, &mu, &txRunning, txDone)
			}
		case <-rxDone:
			log.Printf("Rx routine ended for client %d", a.Id)
			if a.ctx.Err() == nil {
				a.StartRxRoutine(txrxCtx, &mu, &rxRunning, rxDone)
			}
		case <-connectDone:
			log.Printf("Connect routine ended for client %d", a.Id)
			if a.ctx.Err() == nil {
				a.StartManageConnect(txrxCtx, &mu, &connRunning, connectDone)
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}

WAIT:
	done := make(chan struct{})
	go func() {
		for {
			mu.Lock()
			if !txRunning && !rxRunning && !connRunning {
				mu.Unlock()
				break
			}
			mu.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		log.Printf("Graceful shutdown timeout sst server (%d)", a.Id)
	}

	log.Printf("SST server quit.. %d", a.Id)
}
