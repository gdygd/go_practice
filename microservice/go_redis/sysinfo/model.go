package sysinfo

import (
	"go_redis/general"
)

//---------------------------------------------------------------------------
// SystemInfo
//---------------------------------------------------------------------------

type SystemStatus struct {
	Terminate bool  `json:"terminate"`
	SvrUtc    int64 `json:"svrutc"`
	DbSvrComm int   `json:"dbstate"`
}

type SystemInfo struct {
	System  SystemStatus      `json:"system"`
	Process []general.Process `json:"process"` // Process array, index 0:MAIN process, ...
}
