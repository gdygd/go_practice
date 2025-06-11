package tcpserver

import (
	"context"
	"server/general/comm"
	"sync"
)

/*
client := NewSStServer(123, myTcpHandler)
go THRclient(&wg, t, client)

client.Shutdown()
*/
type SStServer struct {
	wg     *sync.WaitGroup // app group 관리 wg
	txrxWg sync.WaitGroup
	ctx    context.Context // 마스터 컨텍스트, 상위/하위 루틴 관리
	cancel context.CancelFunc
	Id     int
	Tcp    *comm.TcpHandler

	txQueue chan []byte
	rxQueue chan []byte
}

/*
tcp := comm.NewTcpHandler("client1", 8000, "127.0.0.1")
client := NewSStServer(123, &tcp)
*/
func NewSStServer(wg *sync.WaitGroup, id int, tcp *comm.TcpHandler) *SStServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &SStServer{
		wg:      wg,
		ctx:     ctx,
		cancel:  cancel,
		Id:      id,
		Tcp:     tcp,
		txQueue: make(chan []byte, 100),
	}
}
