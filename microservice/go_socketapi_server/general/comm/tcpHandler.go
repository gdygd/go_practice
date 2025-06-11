package comm

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// ---------------------------------------------------------------------------
// Implement CommHandler
// ---------------------------------------------------------------------------
type TcpHandler struct {
	name string
	port int
	addr string

	tcp       net.Conn
	connected atomic.Bool
}

var tcpMutex = &sync.Mutex{}

func newTcpHandler(name string, port int, addr string) TcpHandler {
	return TcpHandler{
		name: name,
		port: port,
		addr: addr,
	}
}

// ---------------------------------------------------------------------------
// SetCommEnv
// ---------------------------------------------------------------------------
func (t *TcpHandler) SetCommEnv(name string, port int, addr string) {
	t.name = name
	t.port = port
	t.addr = addr
}

func (t *TcpHandler) GetAddress() string {
	return fmt.Sprintf("%s:%d", t.addr, t.port)
}

// ---------------------------------------------------------------------------
// Connect
// ---------------------------------------------------------------------------
func (t *TcpHandler) Connect() (bool, error) {
	target := fmt.Sprintf("%s:%d", t.addr, t.port)
	conn, err := net.DialTimeout("tcp", target, time.Second*1)
	if err != nil {
		log.Printf("Connect error to %s: %v", target, err)
		return false, err
	}

	t.tcp = conn
	t.connected.Store(true)
	log.Printf("Connected to %s", target)
	return true, nil
}

// ---------------------------------------------------------------------------
// SendMessage
// ---------------------------------------------------------------------------
func (t *TcpHandler) Send(data []byte) (int, error) {
	if !t.connected.Load() {
		return 0, fmt.Errorf("not connected")
	}
	cnt, err := t.tcp.Write(data)
	if err != nil {
		t.connected.Store(false)
	}
	return cnt, err
}

// ---------------------------------------------------------------------------
// Read
// ---------------------------------------------------------------------------
func (t *TcpHandler) Read(data []byte) (int, error) {
	if !t.connected.Load() {
		return 0, fmt.Errorf("not connected")
	}
	cnt, err := t.tcp.Read(data)
	if err != nil {
		t.connected.Store(false)
	}
	return cnt, err
}

func (t *TcpHandler) IsConnected() bool {
	return t.connected.Load()
}

// ---------------------------------------------------------------------------
// ClearEnv
// ---------------------------------------------------------------------------
func (t *TcpHandler) Close() error {
	if !t.connected.Load() {
		return fmt.Errorf("connection already closed for %s:%d", t.addr, t.port)
	}

	t.connected.Store(false)
	if t.tcp != nil {
		if err := t.tcp.Close(); err != nil {

			return err
		} else {
			// log.Printf("Closed connection to %s:%d", t.addr, t.port)
			return nil
		}
	}

	return nil
}
