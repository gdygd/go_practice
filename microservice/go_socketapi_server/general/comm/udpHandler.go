package comm

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

// ---------------------------------------------------------------------------
// Implement CommHandler
// ---------------------------------------------------------------------------
type UdpHandler struct {
	name     string
	sendport int
	recvport int
	addr     string

	udp       net.PacketConn
	udpAddr   net.Addr
	connected atomic.Bool
	connMu    sync.RWMutex // 연결 변경용 락
}

// var udpMutex = &sync.Mutex{}

func newUdpHandler(name string, sendport, recvport int, addr string) UdpHandler {
	return UdpHandler{
		name:     name,
		sendport: sendport,
		recvport: recvport,
		addr:     addr,
	}
}

// ---------------------------------------------------------------------------
// SetCommEnv
// ---------------------------------------------------------------------------
func (u *UdpHandler) SetCommEnv(name string, sendport, recvport int, addr string) {
	u.name = name
	u.sendport = sendport
	u.recvport = recvport
	u.addr = addr
}

// ---------------------------------------------------------------------------
// Connect
// ---------------------------------------------------------------------------
func (u *UdpHandler) Connect() (bool, error) {
	u.connMu.Lock()
	defer u.connMu.Unlock()

	if u.udp != nil {
		_ = u.udp.Close()
		u.udp = nil
	}

	target := fmt.Sprintf("%s:%d", u.addr, u.sendport)
	udpAddr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		log.Printf("ResolveUDPAddr error: %v", err)
		return false, err
	}
	u.udpAddr = udpAddr

	recvAddr := fmt.Sprintf(":%d", u.recvport)
	conn, err := net.ListenPacket("udp", recvAddr)
	if err != nil {
		log.Printf("ListenPacket error on %s: %v", recvAddr, err)
		return false, err
	}

	u.udp = conn
	u.connected.Store(true)
	log.Printf("UDP connected to %s (recv: %s)", target, recvAddr)
	return true, nil
}

// ---------------------------------------------------------------------------
// SendMessage
// ---------------------------------------------------------------------------
func (u *UdpHandler) Send(data []byte) (int, error) {
	if !u.connected.Load() {
		return 0, fmt.Errorf("not connected")
	}

	u.connMu.RLock()
	defer u.connMu.RUnlock()

	if u.udpAddr == nil {
		return 0, fmt.Errorf("destination address not set")
	}
	cnt, err := u.udp.WriteTo(data, u.udpAddr)
	if err != nil {
		u.connected.Store(false)
	}
	return cnt, err
}

// ---------------------------------------------------------------------------
// Read
// ---------------------------------------------------------------------------
func (u *UdpHandler) Read(data []byte) (int, error) {
	if !u.connected.Load() {
		return 0, fmt.Errorf("not connected")
	}

	u.connMu.RLock()
	defer u.connMu.RUnlock()

	cnt, addr, err := u.udp.ReadFrom(data)
	if err != nil {
		u.connected.Store(false)
		return cnt, err
	}

	u.udpAddr = addr
	return cnt, nil
}

func (u *UdpHandler) IsConnected() bool {
	return u.connected.Load()
}

// ---------------------------------------------------------------------------
// Close
// ---------------------------------------------------------------------------
func (u *UdpHandler) Close() error {
	u.connMu.Lock()
	defer u.connMu.Unlock()

	if u.udp != nil {
		err := u.udp.Close()
		u.udp = nil
		u.connected.Store(false)

		if err != nil {
			return fmt.Errorf("UDP close error: %v", err)
		}
		log.Printf("[UdpHandler] Connection closed for %s:%d", u.addr, u.recvport)
	}
	return nil
}
