package container

import (
	"fmt"
	"go_redis/config"
	"go_redis/general"
	"go_redis/general/cache"
	"os"
)

var PRC_DESC = []string{"mpserver", "./beserver"}

// Main app 구동 구조체
type Container struct {
	Process []general.Process
	Rdb     *cache.RedisClient
	Config  *config.Config

	Mlog *general.OLog2
}

var ct *Container

func NewContainer() (*Container, error) {

	config, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("config loading error..%v \n", err)
	}
	ct.Config = &config

	ct.Rdb = cache.NewRedisClient("127.0.0.1:6379")

	ct.Mlog = general.InitLogEnv("./log", "mpserver", 0)

	initProcess()

	return ct, nil
}

func initConfig() (config.Config, error) {

	return config.LoadConfig(".")
}

func initProcess() {

	ct.Mlog.Print(2, "initProcess..")

	ct.Process = make([]general.Process, len(PRC_DESC))
	for idx, _ := range PRC_DESC {
		// ct.Process[idx] = general.Process{PrcName: PRC_DESC[idx]}
		ct.Process[idx] = general.Process{PrcName: PRC_DESC[idx]}
	}

	ct.Process[0].RegisterPid(os.Getpid())
	ct.Process[0].RunBase.Active = true

}
