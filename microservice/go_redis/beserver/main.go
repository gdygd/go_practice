package main

import (
	"flag"
	"fmt"
	"go_redis/beserver/internal/app"
	"go_redis/beserver/internal/container"
	"go_redis/beserver/internal/logger"
	"go_redis/general"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-ini/ini"
)

// ------------------------------------------------------------------------------
// Glocal
// ------------------------------------------------------------------------------
var process *general.Process

// ------------------------------------------------------------------------------
// local
// ------------------------------------------------------------------------------
var ct *container.Container
var process_mode *string
var terminate bool = false
var server *app.Application = nil

var isShutDownApp bool = false

// ------------------------------------------------------------------------------
// const
// ------------------------------------------------------------------------------
const systemini = "./system.ini"

var sharedid int = 0

// ------------------------------------------------------------------------------
// sigHandler
// ------------------------------------------------------------------------------
func sigHandler(chSig chan os.Signal) {
	logger.Mlog.Print(2, "[server]sigHandler")
	for {
		signal := <-chSig
		str := fmt.Sprintf("[server] Accept Signal : %d", signal)
		logger.Mlog.Print(2, "%s", str)
		switch signal {
		case syscall.SIGHUP:
			logger.Mlog.Print(2, "[server]SIGHUP(%d)\n", signal)
		case syscall.SIGINT:
			logger.Mlog.Print(2, "[server]SIGINT(%d)\n", signal)
			shudownApp()
			terminate = true
			// os.Exit(0)
		case syscall.SIGTERM:
			logger.Mlog.Print(2, "SIGTERM(%d)\n", signal)
			terminate = true
			// os.Exit(0)
		case syscall.SIGKILL:
			logger.Mlog.Print(2, "SIGKILL(%d)\n", signal)
			terminate = true
		case syscall.SIGUSR1:
			logger.Mlog.Print(2, "SIGUSR1(%d)\n", signal)
			go shudownApp()
			// os.Exit(0)
		default:
			logger.Mlog.Print(2, "Unknown signal(%d)\n", signal)
			//panic(signal)
		}
	}
}

// ------------------------------------------------------------------------------
// initEnvVaiable
// ------------------------------------------------------------------------------
func initEnvVaiable() bool {
	cfg, err := ini.Load(systemini)
	if err != nil {
		logger.Mlog.Error("[server]fail to read sign_system.ini %v", err)
		return false
	}

	sharedid, _ = cfg.Section("SHARED_ID").Key("id").Int()
	return true
}

// ------------------------------------------------------------------------------
// initContainer
// ------------------------------------------------------------------------------
func initContainer() bool {
	var err error = nil
	ct, err = container.NewContainer()
	if err != nil {
		logger.Mlog.Print(2, "[server]initContainer err.. %v \n", err)
		return false
	}

	return true
}

// ------------------------------------------------------------------------------
// initSignal
// ------------------------------------------------------------------------------
func initSignal() {

	logger.Mlog.Print(2, "[server]initSignal...")
	// signal handler
	ch_signal := make(chan os.Signal, 10)
	signal.Notify(ch_signal, syscall.SIGSEGV, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGUSR1)
	go sigHandler(ch_signal)
}

// ------------------------------------------------------------------------------
// initProcess
// ------------------------------------------------------------------------------
func initProcess() bool {
	logger.Mlog.Print(2, "[server]initProcess..")
	process = &general.Process{PrcName: "./beserver", DebugLv: 2}
	process.RunBase.Active = true
	process.RegisterPid(os.Getpid())
	updateProcessInfo()
	return true
}

// ------------------------------------------------------------------------------
// initEnv
// ------------------------------------------------------------------------------
func initEnv() bool {

	initEnvVaiable()

	// container
	if !initContainer() {
		return false
	}

	// signal
	initSignal()

	initProcess()

	return true

}

func shudownApp() {
	if isShutDownApp {
		return
	}
	isShutDownApp = true
	logger.Mlog.Print(2, "[server]shudownApp..")

	server.ShutDown()
	process.DeActive()

	updateProcessInfo()
}

func clearEnv() {
	logger.Mlog.Print(2, "[server]clearEnv..")
	logger.Mlog.Print(2, "server active#1 : %v ", process.Active)

	if *process_mode == "release" {
		logger.Mlog.Print(2, "Deregister process ..%d ", os.Getpid())
		process.Deregister(os.Getpid())

		updateProcessInfo()
	}

	logger.Mlog.Print(2, "Bye.. :)")

	// app shtudown
}

func updateProcessInfo() {
	ct.ObjDb.SetProcess(*process)
}

func main() {

	process_mode = flag.String("mode", "debug", "프로세스 실행 모드를 선택")
	flag.Parse()
	logger.Mlog.Print(2, "process mode : %s", *process_mode)
	logger.Mlog.Print(2, "남은 인자들:", flag.Args())

	ok := initEnv()
	defer clearEnv()

	logger.Mlog.Print(2, "init state : %v", ok)

	if !ok {
		logger.Mlog.Error("initEnv Error...")
	}

	// NewApplication()
	if ok {
		server = app.NewApplication(ct)
		go server.Start()
	}

	for ok {
		time.Sleep(time.Millisecond * 1000)
		if terminate {
			logger.Mlog.Print(2, "Quit server proces .. ")
			break
		}
		process.MarkTime()
		updateProcessInfo()
		// logger.Mlog.Print(2, "child server tm : %v", process.LastTm)
		// check app server state
	}

}
