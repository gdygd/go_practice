package msgproc

import (
	"context"
	"log"
	"sync"
	"time"

	"grpc_client_test/internal/container"
	"grpc_client_test/internal/logger"
)

// ------------------------------------------------------------------------------
// MsgProcHandler
// ------------------------------------------------------------------------------
type MsgProcHandler struct {
	wg     *sync.WaitGroup // app에서 관리하는 wg
	msgwg  sync.WaitGroup
	ctx    context.Context // master context
	cancel context.CancelFunc
	ct     *container.Container
}

func NewMsgProcHandler(wg *sync.WaitGroup, ct *container.Container) (*MsgProcHandler, error) {
	ctx, cancel := context.WithCancel(context.Background())
	msgHadler := &MsgProcHandler{
		wg:     wg,
		ctx:    ctx,
		cancel: cancel,
		ct:     ct,
	}

	return msgHadler, nil
}

func (m *MsgProcHandler) Shutdown() {
	logger.Log.Print(2, "MsgProc Shutdown #1")
	m.cancel()
	m.msgwg.Wait()
	logger.Log.Print(2, "MsgProc Shutdown #2")
}

func (m *MsgProcHandler) Start() {
	defer m.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("MsgProcHandler panic: %v \n", r)
		}
	}()

	msgCtx, msgCancel := context.WithCancel(m.ctx)

	var mu sync.Mutex
	var msgRunning bool
	msgDone := make(chan struct{}, 1)

	m.startMsgProc(msgCtx, &mu, &msgRunning, msgDone)

	for {
		select {
		case <-m.ctx.Done():
			logger.Log.Print(2, "msgproc handler master context canceled...")
			msgCancel()
			goto WAIT
		case <-msgDone:
			logger.Log.Print(2, "msg proc routine ended...")
			if m.ctx.Err() == nil {
				m.startMsgProc(msgCtx, &mu, &msgRunning, msgDone)
			}
		default:
			time.Sleep(time.Second * 1)
		}
	}
WAIT:
	done := make(chan struct{})
	go func() {
		for {
			mu.Lock()
			if !msgRunning {
				mu.Unlock()
				break
			}
			mu.Unlock()
			time.Sleep(time.Second * 1)
		}
		close(msgDone)
	}()

	select {
	case <-done:
		logger.Log.Print(2, "Msg proc Graceful shutdown!")
	case <-time.After(time.Second * 5):
		logger.Log.Print(2, "Graceful shutdown timeout.. msg handler")
	}

	logger.Log.Print(2, "Msg Proc quit..")
}

func (m *MsgProcHandler) startMsgProc(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}

	*running = true
	mu.Unlock()

	m.msgwg.Add(1)
	go func() {
		defer m.msgwg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()
		logger.Log.Print(2, "Start Msg Handler routine")
		m.msgHandler(ctx)
	}()
}

func (m *MsgProcHandler) msgHandler(ctx context.Context) {
	msg_ch := m.ct.Bus.Subscribe("test_msg")

	for {
		logger.Log.Print(2, "msgHandler...#1")
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "msgHandler routine exiting..")
			return
		case msg := <-msg_ch:
			logger.Log.Print(2, "subscribe  msg! (%v) %v", msg.Topic, msg.Data)
		default:
			// message processing
			logger.Log.Print(2, "msg handler!")
			time.Sleep(time.Second * 1)
		}
	}
}
