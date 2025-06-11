package general

import (
	"log"
	"sort"
	"time"
)

//------------------------------------------------------------------------------
// Constant
//------------------------------------------------------------------------------
const (
	THR_START_ID = 1
	THR_MAX_ID   = 10000
)
const def_thr_check_interval = 1000 // 1sec

// ID Pool
var THR_ID_POOL []int = []int{}

type ThreadFunc func(t *Thread, ch chan bool, arg1, arg2, arg3 interface{})

type Thread struct {
	Interval int // millisecond
	Arg1     interface{}
	Arg2     interface{}
	Arg3     interface{}
	ChKill   chan bool

	StartFunc ThreadFunc
	Timer     time.Time // Thread Check timer

	RunBase // Run Base

}

func checkKillChannel(ch chan bool) {
	if len(ch) >= 1 {
		for len(ch) > 0 {
			<-ch
		}
	}
}

func getThreadID() int {
	id := THR_START_ID

	for idx := 1; idx < len(THR_ID_POOL); idx++ {
		prevId := THR_ID_POOL[idx-1]
		curId := THR_ID_POOL[idx]

		if (prevId + 1) == curId {
			if len(THR_ID_POOL)-1 == idx {
				// 마지막 ID값 도달시
				id = curId + 1
			}
			continue
		}

		if (curId - prevId) > 1 {
			id = prevId + 1
			break
		}
	}

	return id

}

func NewThread() *Thread {
	var t Thread = Thread{}
	t.RunBase.ID = getThreadID()
	sort.Ints(THR_ID_POOL)

	// 여기서 threadID 발급

	return &t
}

func (t *Thread) Init(f ThreadFunc, interval int, args ...interface{}) {

	t.StartFunc = f
	t.Interval = interval
	t.ChKill = make(chan bool, 1)

	for idx, arg := range args {
		if idx > 2 {
			break
		}
		if idx == 0 {
			t.Arg1 = arg
		} else if idx == 1 {
			t.Arg2 = arg
		} else if idx == 2 {
			t.Arg3 = arg
		}
	}
}

func (t *Thread) Start() {
	t.RunBase.Active = true
	time.Sleep(time.Millisecond * 50)
	go t.StartFunc(t, t.ChKill, t.Arg1, t.Arg2, t.Arg3)
}

func (t *Thread) Kill() {
	// 체크 채널
	// 체널 버퍼 1, 체널에 데이터가 있는지 검사
	checkKillChannel(t.ChKill)

	// send kill command
	t.RunBase.Active = false
	t.ChKill <- true
}

func (t *Thread) IsRunning(state *int) bool {
	*state = RST_OK
	if !CheckElapsedTime(&t.Timer, def_thr_check_interval) {
		return true
	}

	// thread 실행 상태 확인
	if !t.RunBase.checkRunInfo() {
		log.Println("[IsRunning] CheckRunInfo Abnomal : ", *state)
		*state = RST_ABNOMAL
		return false
	}

	return true

}

//------------------------------------------------------------------------------
// threadRegister (routine)
//------------------------------------------------------------------------------
func (t *Thread) threadRegister(id int64) {
	t.RunBase.threadRegister(id)
}

//------------------------------------------------------------------------------
// threadDeregister (routine)
//------------------------------------------------------------------------------
func (t *Thread) threadDeregister(id int64) {
	t.RunBase.threadDeregister(id)
}
