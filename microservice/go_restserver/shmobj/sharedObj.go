package shmobj

import (
	"server/general"
)

const MEM_SIZE = 1024 * 5000 // 5000kb
const MEM_KEY = 0x1234

const PRC_IDX_MAIN = 0
const PRC_IDX_PRC01 = 1

const MAX_PROCESS = 2
const MAX_CLIENT = 100

type SharedMemory struct {
	System  SystemInfo
	Process [MAX_PROCESS]general.Process // Process array, index 0:MAIN process, ...
	Client  [MAX_CLIENT]ClientObj
}
