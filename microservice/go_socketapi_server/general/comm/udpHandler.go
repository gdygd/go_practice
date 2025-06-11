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
}

var udpMutex = &sync.Mutex{}

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

	cnt, addr, err := u.udp.ReadFrom(data)
	if err != nil {
		u.connected.Store(false)
		return cnt, err
	}

	u.udpAddr = addr // 최근 응답한 송신자 기준으로 주소 갱신
	return cnt, nil
}

func (u *UdpHandler) IsConnected() bool {
	return u.connected.Load()
}

// ---------------------------------------------------------------------------
// Close
// ---------------------------------------------------------------------------
func (u *UdpHandler) Close() {
	if !u.connected.Load() {
		log.Printf("UDP connection already closed for %s:%d", u.addr, u.recvport)

		return
	}

	u.connected.Store(false)
	if u.udp != nil {
		if err := u.udp.Close(); err != nil {
			log.Printf("UDP close error: %v", err)
		} else {
			log.Printf("UDP connection closed for %s:%d", u.addr, u.recvport)
		}
	}
}
