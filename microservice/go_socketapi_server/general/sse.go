package general

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
)

const (
	CHBUF_SSE = 100
)

//---------------------------------------------------------------------------
// Event Channel
//---------------------------------------------------------------------------
var ChEvent chan EventData = initChannelEvent()

//---------------------------------------------------------------------------
// mutex
//---------------------------------------------------------------------------
var sseMTX = &sync.Mutex{} // sse mutex

//---------------------------------------------------------------------------
// Event Data Format
//---------------------------------------------------------------------------
type EventData struct {
	Msgtype string `json:"event"`
	Data    string `json:"data"`
	Id      string `json:"id"`
}

//---------------------------------------------------------------------------
// initChannelEvent
//---------------------------------------------------------------------------
func initChannelEvent() chan EventData {
	channel := make(chan EventData, 100)
	return channel
}

func clearSSEChannel() {
	for len(ChEvent) > 0 {
		<-ChEvent
	}
}

func checkSSEChannel() {
	if len(ChEvent) >= CHBUF_SSE {
		//Applog.Warn("clearSSEChannel [%d][%d]", len(ChEvent), CHBUF_SSE)
		clearSSEChannel()
	}
}

func SendSSE(data EventData) {
	sseMTX.Lock()
	checkSSEChannel()

	ChEvent <- data
	sseMTX.Unlock()
}

func (e *EventData) PrepareMessage() []byte {
	var data bytes.Buffer

	if len(e.Id) > 0 {
		data.WriteString(fmt.Sprintf("id: %s\n", strings.Replace(e.Id, "\n", "", -1)))
	}
	if len(e.Msgtype) > 0 {
		data.WriteString(fmt.Sprintf("event: %s\n", strings.Replace(e.Msgtype, "\n", "", -1)))
	}
	if len(e.Data) > 0 {
		lines := strings.Split(e.Data, "\n")
		for _, line := range lines {
			data.WriteString(fmt.Sprintf("data: %s\n", line))
		}
	}

	data.WriteString("\n")
	return data.Bytes()
}

func (e *EventData) PrepareMessage2(id string) []byte {
	var data bytes.Buffer
	e.Id = id
	if len(e.Id) > 0 {
		data.WriteString(fmt.Sprintf("id: %s\n", strings.Replace(e.Id, "\n", "", -1)))
	}
	if len(e.Msgtype) > 0 {
		data.WriteString(fmt.Sprintf("event: %s\n", strings.Replace(e.Msgtype, "\n", "", -1)))
	}
	if len(e.Data) > 0 {
		lines := strings.Split(e.Data, "\n")
		for _, line := range lines {
			data.WriteString(fmt.Sprintf("data: %s\n", line))
		}
	}
	data.WriteString("\n")
	return data.Bytes()
}
