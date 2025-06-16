package memory

import "go_redis/general"

const (
	STATE_OFF = 0
	STATE_ON  = 1
)

type SystemInfo struct {
	SvrUtc    int64           `json:"svrutc"`
	DbSvrComm int             `json:"dbstate"`
	Process   general.Process `json:"process"` // Process array, index 0:MAIN process, ...
}
