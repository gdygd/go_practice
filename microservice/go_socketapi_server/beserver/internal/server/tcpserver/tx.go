package tcpserver

import (
	"log"
)

// ------------------------------------------------------------------------------
// ManageTx
// ------------------------------------------------------------------------------
func (a *SStServer) ManageTx() bool {
	log.Printf("tx...#1")

	data := []byte(string("hello"))

	return a.sendMessage(0x01, len(data), data, 1)
}

func (a *SStServer) sendMessage(code byte, length int, info []byte, lv int) bool {
	log.Printf("sendMessage...#1")
	if len(info) < length {
		log.Printf("sendMessage: info too short. expected %d, got %d", length, len(info))
		return false
	}
	log.Printf("sendMessage...#2")

	packet := make([]byte, 0, length+3)
	packet = append(packet, 0x7E, 0x7E, code)
	packet = append(packet, info[:length]...)

	log.Printf("sendMessage...#3")
	_, err := a.Tcp.Send(packet)
	log.Printf("sendMessage...#4")
	if err != nil {
		log.Printf("sendMessage...#4")
		log.Printf("send err, connection close... : id=%v code=%v err=%v", a.Id, code, err)
		a.close()
		return false
	}
	log.Printf("sendMessage...#6")
	return true
}
