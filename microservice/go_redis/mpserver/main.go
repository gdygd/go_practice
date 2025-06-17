package main

import (
	"fmt"
	"go_redis/general"
	"go_redis/mpserver/app"
	"go_redis/mpserver/container"
	"go_redis/mpserver/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var mainApp *app.Application
var ct *container.Container
var Mlog *general.OLog2

// ------------------------------------------------------------------------------
// sigHandler
// ------------------------------------------------------------------------------
func sigHandler(chSig chan os.Signal) {
	for {
		signal := <-chSig
		str := fmt.Sprintf("[main] Accept Signal : %d", signal)
		Mlog.Print(5, "%s", str)
		switch signal {
		case syscall.SIGHUP:
			Mlog.Print(5, "[main]SIGHUP(%d)\n", signal)
		case syscall.SIGINT:
			Mlog.Print(5, "[main]SIGINT(%d)\n", signal)
			clearEnv()

		case syscall.SIGTERM:
			Mlog.Print(5, "[main]SIGTERM(%d)\n", signal)
			clearEnv()

		case syscall.SIGKILL:
			Mlog.Print(5, "[main]SIGKILL(%d)\n", signal)
			clearEnv()

		default:
			Mlog.Print(5, "[main]Unknown signal(%d)\n", signal)

		}
	}
}

// ------------------------------------------------------------------------------
// initSignal
// ------------------------------------------------------------------------------
func initSignal() {
	Mlog.Info("[main]iniSignal")
	// signal handler
	ch_signal := make(chan os.Signal, 10)
	signal.Notify(ch_signal, syscall.SIGSEGV, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGUSR1)
	go sigHandler(ch_signal)
}

// ------------------------------------------------------------------------------
// initContainer
// ------------------------------------------------------------------------------
func initContainer() bool {
	var err error = nil
	ct, err = container.NewContainer()
	if err != nil {
		Mlog.Print(3, "[server]initContainer err.. %v \n", err)
		return false
	}

	return true
}

func checkLogDir() {
	// log directory 존재 확인
	if _, err := os.Stat("./log"); os.IsNotExist(err) {
		// root dir 생성
		os.Mkdir("./log", os.ModePerm)
	}
}

func initEvn() bool {
	checkLogDir()

	// container
	if !initContainer() {
		return false
	}

	// signal
	initSignal()

	return true
}

func childProcessKill() {
	var state int = 0
	for i := 0; i < app.PROCESS_CNT; i++ {
		state |= 1 << i
	}

	// All sub process quit
	var PRC_CNT = app.PROCESS_CNT - 1
	var index int = 1
	var prcstate int = 1

	for {
		idx := (index % PRC_CNT) + 1 // 0 is main, child process starts idx 1
		ptrPrc := &ct.Process[idx]
		if ptrPrc == nil {
			index++
			continue
		}

		Mlog.Always("clear env  처리중.. (%d)%d, %v", idx, ptrPrc.ID, ptrPrc.RunBase.Active)

		if ptrPrc.RunBase.Active {
			Mlog.Always("작업 처리중.. (%d)%d", idx, ptrPrc.ID)
			// child sub app종료 시그널
			ptrPrc.Signal(syscall.SIGUSR1)

		} else {
			Mlog.Always("Kill Process:[%v][%d] [%s]", ptrPrc.RunBase.Active, ptrPrc.RunBase.ID, ptrPrc.PrcName)
			ptrPrc.Kill()
			prcstate |= 1 << idx
		}

		if state == prcstate {
			Mlog.Print(2, "All child process quit")
			break
		}
		time.Sleep(time.Millisecond * 500)

		index++
	}

}

func clearEnv() {
	mainApp.Shutdown() // close manege routine

	childProcessKill()

	mainApp.ShutDownSubscribe()

	Mlog.Print(4, "All clear main process, ")
}

func main() {
	Mlog = logger.Mlog

	initEvn()
	mainApp = app.NewApplication(ct)
	go mainApp.Start()
	Mlog.Print(2, "App Start")

	for {
		if mainApp.IsAppClosed() {
			Mlog.Print(2, "Quit main server.. ")
			break
		}
	}

	Mlog.Print(4, "Main process quit.. Bye~~:)")

}
