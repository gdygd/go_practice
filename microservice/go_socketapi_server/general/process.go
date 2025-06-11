package general

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

//------------------------------------------------------------------------------
// Constant
//------------------------------------------------------------------------------

const def_proc_check_interval = 1000 // 1sec

type Process struct {
	PrcName  string
	PrcName2 [20]byte
	LastTm   int64 // 마지막 루프타임	(utc time)
	Cmd      []string
	DebugLv  int       // Debug level
	Timer    time.Time // Process Check timer

	RunBase // Run Base

}

// ------------------------------------------------------------------------------
// InitProcess
// ------------------------------------------------------------------------------
func InitProcess(name string, cmd []string) Process {
	b := []byte(name)
	//copy(arr[:], slice[:4])
	var arr [20]byte
	copy(arr[:], b[:])

	var proc Process = Process{PrcName: name, PrcName2: arr, Cmd: cmd, DebugLv: 3}
	return proc
}

// ------------------------------------------------------------------------------
// IsActiveProcess
// ------------------------------------------------------------------------------
func (p *Process) IsActiveProcess(pid int) bool {
	if p.RunBase.ID == pid {
		return true
	}
	return false
}

func (p *Process) RegisterPid(pid int) {
	p.RunBase.register(pid)
}

func (p *Process) Deregister(pid int) {
	p.RunBase.deregister(pid)
}

func (p *Process) DeActive() {
	p.RunBase.Active = false
}

// ------------------------------------------------------------------------------
// IsRunning
// ------------------------------------------------------------------------------
func (p *Process) IsRunning(state *int) bool {
	*state = RST_OK
	if !CheckElapsedTime(&p.Timer, def_proc_check_interval) {
		return true
	}

	// 프로세스 존재 확인
	if !p.IsExist() {
		log.Println("[IsRunning] UnExist : ", *state)
		*state = RST_UNEXIST
		return false
	}

	// 프로세스 실행 상태 확인
	if !p.RunBase.checkRunInfo() {
		log.Println("[IsRunning] CheckRunInfo Abnomal : ", *state)
		*state = RST_ABNOMAL
		return false
	}

	return true
}

func (p *Process) MarkTime() {
	p.LastTm = time.Now().Unix()
}

// ------------------------------------------------------------------------------
// IsExist
// ------------------------------------------------------------------------------
func (p *Process) IsExist() bool {

	var ok bool = true

	proc, err := os.FindProcess(p.RunBase.ID)
	err = proc.Signal(syscall.Signal(0))

	if err != nil {
		log.Println("## > unexist : ", p.PrcName, p.RunBase.ID)
		ok = false
	}

	return (ok)
}

// ------------------------------------------------------------------------------
// Kill
// ------------------------------------------------------------------------------
func (p *Process) Kill() bool {

	var ok bool = false

	process, err := os.FindProcess(p.RunBase.ID)

	err = process.Signal(syscall.Signal(0))

	if err != nil {
		return true
	} else {
		//err1 := process.Kill()
		err1 := process.Signal(syscall.SIGTERM)
		if err1 == nil {
			ok = true
		}
	}

	return (ok)
}

// ------------------------------------------------------------------------------
// Kill
// ------------------------------------------------------------------------------
func (p *Process) Signal(sig os.Signal) bool {

	var ok bool = false

	process, err := os.FindProcess(p.RunBase.ID)

	err = process.Signal(syscall.Signal(0))

	if err != nil {
		return true
	} else {
		//err1 := process.Kill()
		err1 := process.Signal(sig)
		if err1 == nil {
			ok = true
		}
	}

	return (ok)
}

// ------------------------------------------------------------------------------
// Start
// ------------------------------------------------------------------------------
func (p *Process) Start() (bool, int) { // (process name, command)

	var ok bool = false
	var pid int = -1
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	if process, err := os.StartProcess(p.PrcName, p.Cmd, &procAttr); err != nil {
		log.Println("ERROR Unable to run %s: %s\n", p.PrcName, err.Error())
	} else {
		// Active
		p.RunBase.Active = true
		ok = true

		//
		pid = process.Pid
		log.Println("%s running as pid %d\n", p.PrcName, process.Pid)
	}

	return ok, pid
}

// ------------------------------------------------------------------------------
// Start
// ------------------------------------------------------------------------------
func (p *Process) Start2(param string) (bool, int) { // (process name, command)
	var ok bool = true
	var pid int = -1
	//cmd := exec.Command(p.PrcName, p.Cmd...)
	cmd := exec.Command(p.PrcName, param)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Start2 (1)" + p.PrcName)

	go func() {
		err := cmd.Run()
		if err != nil {
			log.Println("[Start2] err", err.Error())
			ok = false
		} else {
			log.Println("[Start2] ok", cmd)
			p.RunBase.Active = true
			ok = true
		}

	}()

	// err := cmd.Run()
	// if err != nil {
	// 	log.Println("[Start2]", err.Error())
	// 	ok = false
	// } else {
	// 	log.Println(cmd)
	// 	p.RunObj.Active = true
	// 	ok = true
	// }

	log.Println("Start2 (2)")
	p.LastTm = time.Now().Unix()

	return ok, pid

}

// ------------------------------------------------------------------------------
//
// ------------------------------------------------------------------------------
func (p *Process) Start3() (bool, int) { // (process name, command)
	var pid int = -1
	binary, lookErr := exec.LookPath(p.PrcName)
	var ok bool = true

	if lookErr != nil {
		log.Println("[Start3] lookErr ", lookErr.Error())
		ok = false
	}

	env := os.Environ()

	go func() {
		execErr := syscall.Exec(binary, p.Cmd, env)
		if execErr != nil {
			log.Println("[Start3] execErr ", execErr.Error())
			ok = false
		}
	}()

	return ok, pid

}

// ------------------------------------------------------------------------------
//
// ------------------------------------------------------------------------------
func (p *Process) Start4() (bool, int) { // (process name, command)
	var ok bool = true
	var pid int = -1

	cmd := exec.Command(p.PrcName)

	log.Println("Start4 (1)" + p.PrcName)

	go func() {
		err := cmd.Start()
		if err != nil {
			log.Println("[Start4] err", err.Error())
			ok = false
		} else {
			log.Println("[Start4] ok", cmd)
			p.RunBase.Active = true
			ok = true
		}

	}()

	log.Println("Start2 (2)")

	return ok, pid

}

// ------------------------------------------------------------------------------
// Start5
// ------------------------------------------------------------------------------
func (p *Process) Start5() (bool, int) { // (process name, command)

	var ok bool = false
	var pid int = -1
	// var procAttr os.ProcAttr
	// procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	sttyArgs := syscall.ProcAttr{
		"",
		[]string{},
		[]uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		nil,
	}

	if prcid, err := syscall.ForkExec(p.PrcName, p.Cmd, &sttyArgs); err != nil {
		log.Println("ERROR Unable to run %s: %s\n", p.PrcName, err.Error())
	} else {
		// Active
		p.RunBase.Active = true
		ok = true

		//
		pid = prcid
		log.Println("%s running as pid %d\n", p.PrcName, prcid)
	}

	return ok, pid
}

func (p *Process) GetPid() int {
	return (p.RunBase.ID)
}

// ------------------------------------------------------------------------------
// GetPNameByArr
// ------------------------------------------------------------------------------
func (p *Process) GetPNameByArr() string {
	var b []byte = []byte{}
	idx := 0
	for {
		if len(p.PrcName2) == idx {
			break
		}
		if p.PrcName2[idx] == 0 {
			break
		}
		b = append(b, p.PrcName2[idx])
		idx++
	}
	return string(b[:])
}

// ------------------------------------------------------------------------------
// SetDebugLv
// ------------------------------------------------------------------------------
func (p *Process) SetDebugLv(lv int) {
	p.DebugLv = lv
}
