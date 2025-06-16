package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go_redis/config"
	"go_redis/general"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

// ------------------------------------------------------------------------------
// Local
// ------------------------------------------------------------------------------
var process []general.Process
var PRC_DESC = []string{"mpserver", "./beserver"}

var ctx = context.Background()
var Mlog *general.OLog2 = general.InitLogEnv("./log", "mpserver", 0)
var appenv config.Config
var isAleayProcess bool = false
var terminate bool = false
var rdb *redis.Client = nil

// ------------------------------------------------------------------------------
// const
// ------------------------------------------------------------------------------
const systemini = "./system.ini"
const PROCESS_CNT = 2 // main, child

func subscribeProcessState() {
	// rdb.Get

	for {

		raw, err := rdb.Get(ctx, "prc_beserver").Result()
		if err != nil {
			log.Fatalf("Redis Get error: %v", err)
		}

		var result general.Process
		err = json.Unmarshal([]byte(raw), &result)
		if err != nil {
			log.Fatalf("JSON unmarshal error: %v", err)
		}

		time.Sleep(time.Millisecond * 100)
	}

}

// ------------------------------------------------------------------------------
// sigHandler
// ------------------------------------------------------------------------------
func sigHandler(chSig chan os.Signal) {
	for {
		signal := <-chSig
		str := fmt.Sprintf("[main] Accept Signal : %d", signal)
		Mlog.Info("%s", str)
		switch signal {
		case syscall.SIGHUP:
			Mlog.Info("[main]SIGHUP(%d)\n", signal)
		case syscall.SIGINT:
			Mlog.Info("[main]SIGINT(%d)\n", signal)
			terminate = true
		case syscall.SIGTERM:
			Mlog.Info("[main]SIGTERM(%d)\n", signal)
			terminate = true
		case syscall.SIGKILL:
			Mlog.Info("[main]SIGKILL(%d)\n", signal)
			terminate = true
		default:
			Mlog.Info("[main]Unknown signal(%d)\n", signal)
			terminate = true
		}
	}
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// ------------------------------------------------------------------------------
// initProcessDesc
// ------------------------------------------------------------------------------
func initProcessDesc() {
	Mlog.Always("initProcessDesc..")

	process = make([]general.Process, len(PRC_DESC))
	for idx, _ := range PRC_DESC {
		process[idx] = general.Process{PrcName: PRC_DESC[idx]}
	}
}

// ------------------------------------------------------------------------------
// initEnvVaiable
// ------------------------------------------------------------------------------
func initEnvVaiable() bool {
	var err error = nil
	appenv, err = config.LoadConfig(".")
	if err != nil {
		Mlog.Error("initEnvVaiable, LoadConfig error..%v", err)
		return false
	}

	return true
}

// ------------------------------------------------------------------------------
// initProcess
// ------------------------------------------------------------------------------
func initProcess() bool {
	// process initialize & start

	Mlog.Info("[main]initProcess")

	// Register PID
	process[0].RegisterPid(os.Getpid())
	process[0].RunBase.Active = true

	updateProcessInfo()

	return true
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
// isAlreadyProcess
// ------------------------------------------------------------------------------
func isAlreadyProcess() bool {
	// 프로세스 러닝상태 확인
	// 이미 프로세스가 동작중이면 false

	Mlog.Info("[main]isAlreadyProcess...")

	var isrunning bool = false

	Mlog.Print(3, "prcdesc : %v", PRC_DESC[1])

	prcnm := PRC_DESC[1] // 코어프로세스
	cmdstr := fmt.Sprintf(`ps -ef | grep %s | grep -v grep`, prcnm)
	cmd := exec.Command("bash", "-c", cmdstr)
	output, _ := cmd.CombinedOutput()
	stroutput := string(output)

	Mlog.Print(3, "cmd output : %v %d", stroutput, len(stroutput))

	if len(stroutput) == 0 {
		isrunning = false
	} else {
		isrunning = true
	}

	return isrunning
}

func SetDebugLv() {
	Mlog.Info("[main]SetDebugLv...")
	process[0].SetDebugLv(appenv.DebugLv)
	updateProcessInfo()
}

// ------------------------------------------------------------------------------
// initEnv
// ------------------------------------------------------------------------------
func initEnv() bool {
	// check process running...
	if isAlreadyProcess() {
		Mlog.Error("[main]프로세스 is running ...")
		isAleayProcess = true
		return false
	}

	Mlog.Info("[main]initEnv ...")

	if !initEnvVaiable() {
		Mlog.Error("[main] 환경변수 초기화FAIL...")
		return false
	}

	initSignal()

	initRedis()

	initProcessDesc() // process name and process command
	if !initProcess() {
		Mlog.Info("Process initialize fail..")
		return false
	}

	if !startProcess() {
		return false
	}

	go subscribeProcessState()

	Mlog.Info("[main]initEnv ok")

	SetDebugLv()

	return true
}

func startProcess() bool {
	Mlog.Print(2, "start process..")
	var isok bool = true
	var wg sync.WaitGroup
	wg.Add(1)

	for idx := 1; idx < PROCESS_CNT; idx++ {
		ptrPrc := &process[idx]
		ptrPrc.Start2("-mode=release")
	}

	//check running state

	var PRC_CNT = PROCESS_CNT - 1
	var index int = 1
	var prcstate int = 1
	var state int = 0
	for i := 0; i < PROCESS_CNT; i++ {
		state |= 1 << i
	}

	curtm := time.Now()
	for {
		if time.Since(curtm) > time.Second*3 {
			Mlog.Error("child process start error... state(%X)", state)
			isok = false
			break
		}

		idx := (index % PRC_CNT) + 1
		ptrPrc := &process[idx]
		if ptrPrc == nil {
			index++
			continue
		}

		if ptrPrc.RunBase.Active {
			prcstate |= 1 << idx

		} else {
			continue
		}

		if state == prcstate {
			Mlog.Print(2, "All child process start")
			break
		}
		time.Sleep(time.Millisecond * 100)

		index++
	}

	return isok

}

func manageProcess() {

	for idx := 1; idx < PROCESS_CNT; idx++ {

		ptrPrc := &process[idx]

		prctm := time.Unix(ptrPrc.LastTm, 0)

		// Mlog.Print(2, "child process id %d, %v", ptrPrc.ID, prctm)

		if !ptrPrc.IsExist() {
			Mlog.Warn("UNEXIST [%s %d]", ptrPrc.PrcName, ptrPrc.GetPid())
			Mlog.Warn("Process start [%s]", ptrPrc.PrcName)
			ptrPrc.Start2("-mode=release")
		} else if time.Since(prctm) > appenv.PROCESS_INTERVAL {
			Mlog.Warn("RST_ABNOMAL [%s %d]", ptrPrc.PrcName, ptrPrc.GetPid())
			// Mlog.Warn("Processkill [%s %d]", ptrPrc.PrcName, ptrPrc.GetPid())
			// ptrPrc.Kill()
		}
	}
}

func updateProcessInfo() {
	// ct.ObjDb.SetProcess(*process)
	// rdb.Set(ctx, "mpserver", process[0], 0)

	data, err := json.Marshal(process[0])
	if err != nil {
		log.Fatalf("JSON marshal error: %v", err)
	}

	err = rdb.Set(ctx, "mpserver", data, time.Hour).Err()
	if err != nil {
		log.Fatalf("[main]Redis Set (mpserver) error: %v", err)
	}
}

// ------------------------------------------------------------------------------
// clearEnv
// ------------------------------------------------------------------------------
func clearEnv() {
	if isAleayProcess {
		Mlog.Print(6, "[main] Is Aleady process")
		Mlog.Print(6, "[main] Process quit, byebye~ :)")
		return
	}

	var state int = 0
	for i := 0; i < PROCESS_CNT; i++ {
		state |= 1 << i
	}

	// All sub process quit
	var PRC_CNT = PROCESS_CNT - 1
	var index int = 1
	var prcstate int = 1

	for {
		idx := (index % PRC_CNT) + 1
		ptrPrc := &process[idx]
		if ptrPrc == nil {
			index++
			continue
		}

		// Mlog.Always("clear env  처리중.. (%d)%d, %v", idx, ptrPrc.ID, ptrPrc.RunBase.Active)

		if ptrPrc.RunBase.Active {
			// Mlog.Always("작업 처리중.. (%d)%d", idx, ptrPrc.ID)
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

	// mpserver clear pid
	ptrPrc := &process[0]
	if ptrPrc != nil {
		ptrPrc.Deregister(os.Getpid())
		Mlog.Always("clearEnv pid (%d) (%d) (%v)", os.Getpid(), ptrPrc.RunBase.ID, ptrPrc.RunBase.Active)
		ptrPrc.RunBase.Active = false
	}
	updateProcessInfo()

	Mlog.Print(2, "[main] Process quit, byebye~ :)")

	// 로그파일 close
	//pmlog := &Mlog
	Mlog.Fileclose()
}

func checkLogDir() {
	// log directory 존재 확인
	if _, err := os.Stat("./log"); os.IsNotExist(err) {
		// root dir 생성
		os.Mkdir("./log", os.ModePerm)
	}
}

func main() {
	var initOk bool = false
	checkLogDir()
	Mlog.Info("%s", "[main] Process start")

	initOk = initEnv()
	defer clearEnv()

	Mlog.Print(4, "Mng.. %v, %v", initOk, terminate)

	for {
		if !initOk || terminate {
			break
		}
		// manage process
		manageProcess()

		process[0].MarkTime()
		updateProcessInfo()
		time.Sleep(time.Millisecond * 100)
	}

	Mlog.Info("[main] Process end..")

}
