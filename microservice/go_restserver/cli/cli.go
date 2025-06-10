package main

import (
	"fmt"
	"os"
	"os/signal"
	"server/cli/cmd"
	"strconv"
	"syscall"
	"unsafe"

	"server/shmobj"

	"gopkg.in/ini.v1"

	"github.com/ghetzel/shmtool/shm"
)

// ------------------------------------------------------------------------------
// const
// ------------------------------------------------------------------------------
const MemorySize = 1024 * 100 // 10kb
const systemini = "./system.ini"
const skey = 0x1234 // shared memory key

// ------------------------------------------------------------------------------
// local
// ------------------------------------------------------------------------------
var SharedMem *shmobj.SharedMemory
var pSegment *shm.Segment
var terminate bool = false

// ------------------------------------------------------------------------------
// sigHandler
// ------------------------------------------------------------------------------
func sigHandler(chSig chan os.Signal) {
	for {
		signal := <-chSig
		fmt.Printf("Accept Signal %d", signal)
		switch signal {
		case syscall.SIGHUP: // 터미널 연결 끊겼을경우
			fmt.Printf("SIGHUP(%d)\n", signal)
			terminate = true
			clearEnv()
		case syscall.SIGINT:
			fmt.Printf("SIGINT(%d)\n", signal)
			terminate = true
			clearEnv()
		case syscall.SIGTERM:
			fmt.Printf("SIGTERM(%d)\n", signal)
			terminate = true
			clearEnv()
		case syscall.SIGKILL:
			fmt.Printf("SIGKILL(%d)\n", signal)
			terminate = true
			clearEnv()
		default:
			fmt.Printf("Unknown signal(%d)\n", signal)
			terminate = true
			clearEnv()
		}
	}
}

// ------------------------------------------------------------------------------
// initSignal
// ------------------------------------------------------------------------------
func initSignal() {
	fmt.Printf("iniSignal...\n")
	// signal handler
	ch_signal := make(chan os.Signal, 10)
	signal.Notify(ch_signal, syscall.SIGSEGV, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGUSR1)
	go sigHandler(ch_signal)
}

// ------------------------------------------------------------------------------
// initMemory
// ------------------------------------------------------------------------------
func initMemory() bool {
	fmt.Printf("initMemory...\n")

	// Write shared id to system.ini
	cfg, err := ini.Load(systemini)
	if err != nil {
		fmt.Printf("fail to read sign_system.ini %v\n", err)
		return false
	}

	id, _ := cfg.Section("SHARED_ID").Key("id").Int()

	psg, err := shm.Open(id)
	if err != nil {
		fmt.Printf("shared memory open fail\n")
		return false
	}

	pSegment = psg

	shmId := pSegment.Id
	strShmId := strconv.Itoa(shmId)

	// Save ini file to shared memory id
	cfg.Section("SHARED_ID").Key("id").SetValue(strShmId)
	cfg.SaveTo(systemini)

	usPtr, err := pSegment.Attach()
	if err != nil {
		fmt.Printf("shared memory attach fail\n")
		return false
	}

	// shared memory attach
	SharedMem = (*shmobj.SharedMemory)(unsafe.Pointer(usPtr))

	return true
}

// ------------------------------------------------------------------------------
// initEnv
// ------------------------------------------------------------------------------
func initEnv() bool {
	initSignal()

	if !initMemory() {
		fmt.Printf("Share memory open fail..\n")
		return false
	}

	return true
}

// ------------------------------------------------------------------------------
// clearEnv
// ------------------------------------------------------------------------------
func clearEnv() {
	// detach memory
	fmt.Printf("[cli] memory detach (1)\n")
	addr := unsafe.Pointer(SharedMem)
	pSegment.Detach(addr)
	fmt.Printf("[cli] memory detach (2)\n")
	os.Exit(0)
}

func main() {
	fmt.Println("Command Line Interface :)\n")

	var initOk bool = false
	initOk = initEnv()

	// input command
	cmdStr := make([]string, 100)
	command := make([]interface{}, 100)
	for i := range cmdStr {
		command[i] = &cmdStr[i]
	}

	cli := cmd.NewCLI()
	cli.InitialMessage()

	if initOk {
		cli.SetShmMemory(SharedMem)
	}

	//defer clearEnv()
	for !cli.Exit && !terminate {

		if !initOk || terminate {
			fmt.Println("Exit Command Line Interface bye~~ :)\n")
			break
		}

		if terminate {
			fmt.Println("terminate Command Line Interface bye~~ :)\n")
			break
		}

		fmt.Printf("CLI>>")
		count, _ := fmt.Scanln(command...)

		if count > 9 {
			fmt.Println("Invalid command")
			continue
		}

		if count == 0 {
			continue
		}

		cli.SetCommand(cmdStr[0:count])
		cli.PrintCmd()

		cli.Run()

		if cli.Terminate {
			fmt.Println("cli process terminate")
			break
		}
	}
	clearEnv()
}
