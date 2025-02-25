package main

import (
	"fmt"

	"github.com/open-traffic-generator/snappi/gosnappi"
)

func main() {
	// 새로운 API 클라이언트 생성
	api := gosnappi.NewApi()
	config := api.NewConfig()

	// 포트 추가
	port := config.Ports().Add()
	port.SetName("port1").SetLocation("localhost:5555")

	// 트래픽 플로우 생성
	flow := config.Flows().Add()
	flow.SetName("HTTP Traffic Flow")
	flow.TxRx().Device().SetTxNames([]string{"port1"}).SetRxNames([]string{"port1"})

	// 트래픽 속도 설정
	flow.Rate().SetPps(10000) // 초당 패킷 수 설정

	// HTTP 패킷 추가
	httpPacket := flow.Packet().Add().Ethernet().IPv4().TCP().Http()
	httpPacket.Ethernet().Src().SetValue("00:0c:29:6b:72:02")
	httpPacket.Ethernet().Dst().SetValue("00:0c:29:6b:72:01")
	httpPacket.Ipv4().Src().SetValue("192.168.1.1")
	httpPacket.Ipv4().Dst().SetValue("192.168.1.2")
	httpPacket.Tcp().SrcPort().SetValue(12345)
	httpPacket.Tcp().DstPort().SetValue(80)
	httpPacket.Http().RequestMethod().SetValue("GET")
	httpPacket.Http().RequestUri().SetValue("/index.html")

	// 트래픽 발생기 시작
	if err := api.Start(config); err != nil {
		fmt.Println("Error starting traffic generator:", err)
	}
}
