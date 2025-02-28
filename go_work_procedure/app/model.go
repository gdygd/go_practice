package app

type Result struct {
	Res  int // 0: 성공, 1: 실패
	Msg  string
	Data any
}
