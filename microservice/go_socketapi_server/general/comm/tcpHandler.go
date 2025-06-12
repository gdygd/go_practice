package comm

import (
	"fmt"
	"log"
	"net"
	"strconv"
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
	connMu    sync.RWMutex // 연결 변경용 락
}

// var tcpMutex = &sync.Mutex{}

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
	t.connMu.Lock()
	defer t.connMu.Unlock()

	if t.tcp != nil {
		_ = t.tcp.Close()
		t.tcp = nil
	}

	target := net.JoinHostPort(t.addr, strconv.Itoa(t.port))
	conn, err := net.DialTimeout("tcp", target, 3*time.Second)
	if err != nil {
		t.connected.Store(false)
		return false, fmt.Errorf("connect to %s failed: %w", target, err)
	}

	t.tcp = conn
	t.connected.Store(true)
	log.Printf("[TcpHandler] Connected to %s", target)
	return true, nil
}

// ---------------------------------------------------------------------------
// SendMessage
// ---------------------------------------------------------------------------
func (t *TcpHandler) Send(data []byte) (int, error) {
	if !t.connected.Load() {
		return 0, fmt.Errorf("not connected")
	}

	t.connMu.RLock()
	defer t.connMu.RUnlock()

	if t.tcp == nil {
		t.connected.Store(false)
		return 0, fmt.Errorf("connection is nil")
	}

	n, err := t.tcp.Write(data)
	if err != nil {
		t.connected.Store(false)
	}
	return n, err
}

// ---------------------------------------------------------------------------
// Read
// ---------------------------------------------------------------------------
func (t *TcpHandler) Read(data []byte) (int, error) {
	if !t.connected.Load() {
		return 0, fmt.Errorf("not connected")
	}

	t.connMu.RLock()
	defer t.connMu.RUnlock()

	if t.tcp == nil {
		t.connected.Store(false)
		return 0, fmt.Errorf("connection is nil")
	}

	n, err := t.tcp.Read(data)
	if err != nil {
		t.connected.Store(false)
	}
	return n, err
}

func (t *TcpHandler) IsConnected() bool {
	return t.connected.Load()
}

// ---------------------------------------------------------------------------
// ClearEnv
// ---------------------------------------------------------------------------
func (t *TcpHandler) Close() error {
	t.connMu.Lock()
	defer t.connMu.Unlock()

	if t.tcp != nil {
		err := t.tcp.Close()
		t.tcp = nil
		t.connected.Store(false)

		if err != nil {
			return fmt.Errorf("tcp close error: %w", err)
		}
		log.Printf("[TcpHandler] Closed connection to %s:%d", t.addr, t.port)
	}
	return nil
}
