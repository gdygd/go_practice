package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"server/config"
	"server/shmobj"
	"strconv"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"server/general"

	"github.com/ghetzel/shmtool/shm"
	"gopkg.in/ini.v1"
)

// ------------------------------------------------------------------------------
// Glocal
// ------------------------------------------------------------------------------
var SharedMem *shmobj.SharedMemory
var process *general.Process
var Mlog *general.OLog2 = general.InitLogEnv("./log", "mpserver", 0)

// ------------------------------------------------------------------------------
// const
// ------------------------------------------------------------------------------
const systemini = "./system.ini"
const sysenvini = "./sys_env.ini"

// ------------------------------------------------------------------------------
// Local
// ------------------------------------------------------------------------------
var isAleayProcess bool = false
var PRC_DESC = []string{"mpserver", "./beserver"}
var prcArgv = [][]string{{""}, {""}}

var pSegment *shm.Segment

var appenv config.Config

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
			SharedMem.System.Terminate = true
		case syscall.SIGTERM:
			Mlog.Info("[main]SIGTERM(%d)\n", signal)
			SharedMem.System.Terminate = true
		case syscall.SIGKILL:
			Mlog.Info("[main]SIGKILL(%d)\n", signal)
			SharedMem.System.Terminate = true
		default:
			Mlog.Info("[main]Unknown signal(%d)\n", signal)
			SharedMem.System.Terminate = true
		}
	}
}

// ------------------------------------------------------------------------------
// initProcessDesc
// ------------------------------------------------------------------------------
func initProcessDesc() {
	Mlog.Always("initProcessDesc..")
	for idx := 0; idx < shmobj.MAX_PROCESS; idx++ {
		SharedMem.Process[idx] = general.InitProcess(PRC_DESC[idx], prcArgv[idx])
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
	process = &SharedMem.Process[shmobj.PRC_IDX_MAIN]

	for idx := 0; idx < shmobj.MAX_PROCESS; idx++ {
		Mlog.Info(PRC_DESC[idx])
		if idx == shmobj.PRC_IDX_MAIN {
			SharedMem.Process[idx].RunBase.Active = true
			continue
		} else {
			Mlog.Info(">>%v", PRC_DESC[idx])
			// SharedMem.Process[idx].Start2("-mode=release")
		}
	}

	return true
}

// ------------------------------------------------------------------------------
// initMemory
// ------------------------------------------------------------------------------
func initMemory() bool {
	Mlog.Info("initMemory")
	// create memory
	psg, err := shm.Create(shmobj.MEM_SIZE) // shared memory size

	if err != nil {
		Mlog.Info("shared memory create fail")
		return false
	}

	pSegment = psg

	// Write shared id to system.ini
	cfg, err := ini.Load(systemini)
	if err != nil {
		Mlog.Info("fail to read sign_system.ini:%v", err)
		return false
	}

	id, _ := cfg.Section("SHARED_ID").Key("id").Int()
	Mlog.Info("[main]iniID : %d", id)

	shmId := pSegment.Id
	strShmId := strconv.Itoa(shmId)
	Mlog.Info("[main]sharedID : %s", strShmId)

	// Save ini file to shared memory id
	cfg.Section("SHARED_ID").Key("id").SetValue(strShmId)
	cfg.SaveTo(systemini)

	usPtr, err := pSegment.Attach()
	if err != nil {
		Mlog.Info("shared memory attach fail")
		return false
	}

	Mlog.Info(">>segment : %v", pSegment)
	// shared memory attach
	SharedMem = (*shmobj.SharedMemory)(unsafe.Pointer(usPtr))
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
	process.SetDebugLv(appenv.DebugLv)
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

	if !initMemory() {
		Mlog.Info("Share memory created fail..")
		return false
	}

	initProcessDesc() // process name and process command
	if !initProcess() {
		Mlog.Info("Process initialize fail..")
		return false
	}

	if !startProcess() {
		return false
	}

	Mlog.Info("[main]initEnv ok")

	// Register PID
	process.RegisterPid(os.Getpid())
	process.RunBase.Active = true

	SetDebugLv()

	return true
}

func startProcess() bool {
	Mlog.Print(2, "start process..")
	var isok bool = true
	var wg sync.WaitGroup
	wg.Add(1)

	for idx := 1; idx < shmobj.MAX_PROCESS; idx++ {
		ptrPrc := &SharedMem.Process[idx]
		ptrPrc.Start2("-mode=release")
	}

	//check running state

	var PRC_CNT = shmobj.MAX_PROCESS - 1
	var index int = 1
	var prcstate int = 1
	var state int = 0
	for i := 0; i < shmobj.MAX_PROCESS; i++ {
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
		ptrPrc := &SharedMem.Process[idx]
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

	for idx := 1; idx < shmobj.MAX_PROCESS; idx++ {

		ptrPrc := &SharedMem.Process[idx]

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
	for i := 0; i < shmobj.MAX_PROCESS; i++ {
		state |= 1 << i
	}

	// All sub process quit
	var PRC_CNT = shmobj.MAX_PROCESS - 1
	var index int = 1
	var prcstate int = 1

	for {
		idx := (index % PRC_CNT) + 1
		ptrPrc := &SharedMem.Process[idx]
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

	// detach memory
	// destroy shared memory
	if pSegment != nil {
		pSegment.Destroy()
	}

	Mlog.Print(2, "memery destroy")

	// mpserver clear pid
	ptrPrc := &SharedMem.Process[0]
	if ptrPrc != nil {
		ptrPrc.Deregister(os.Getpid())
		Mlog.Always("clearEnv pid (%d) (%d) (%v)", os.Getpid(), ptrPrc.RunBase.ID, ptrPrc.RunBase.Active)
		ptrPrc.RunBase.Active = false
	}

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

	Mlog.Print(4, "Mng.. %v, %v", initOk, SharedMem.System.Terminate)

	for {
		if !initOk || SharedMem.System.Terminate {
			break
		}
		// manage process
		manageProcess()

		process.MarkTime()
		time.Sleep(time.Millisecond * 100)
	}

	Mlog.Info("[main] Process end..")

}
