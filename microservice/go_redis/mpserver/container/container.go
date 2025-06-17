package container

import (
	"fmt"
	"go_redis/config"
	"go_redis/general"
	"go_redis/general/cache"
	"go_redis/mpserver/logger"
	"os"
)

var PRC_DESC = []string{"mpserver", "./beserver"}

type SystemInfo struct {
	Terminate bool
	SvtUtc    int64
}

// Main app 구동 구조체
type Container struct {
	Process []general.Process
	Rdb     *cache.RedisClient
	Config  *config.Config
	SysInfo SystemInfo
}

var ct *Container
var Mlog *general.OLog2

func NewContainer() (*Container, error) {
	ct = &Container{}
	Mlog = logger.Mlog
	Mlog.Print(2, "NewContainer ...#1")

	config, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("config loading error..%v \n", err)
	}
	Mlog.Print(2, "NewContainer ...#2")
	ct.Config = &config

	Mlog.Print(2, "NewContainer ...#2.1")

	ct.Rdb = cache.NewRedisClient("127.0.0.1:6379")

	Mlog.Print(2, "NewContainer ...#3")
	initProcess()

	Mlog.Print(2, "NewContainer ...#4")
	return ct, nil
}

func initConfig() (config.Config, error) {

	return config.LoadConfig(".")
}

func initProcess() {

	Mlog.Print(2, "initProcess..")

	ct.Process = make([]general.Process, len(PRC_DESC))
	for idx, _ := range PRC_DESC {
		// ct.Process[idx] = general.Process{PrcName: PRC_DESC[idx]}
		ct.Process[idx] = general.Process{PrcName: PRC_DESC[idx]}
	}

	ct.Process[0].RegisterPid(os.Getpid())
	ct.Process[0].RunBase.Active = true

}
