package app

import (
	"context"
	"encoding/json"
	"fmt"
	"go_redis/general"
	"go_redis/mpserver/container"
	"go_redis/mpserver/logger"
	"strconv"
	"sync"
	"time"
)

var Mlog *general.OLog2

const PROCESS_CNT = 2 // main, child

type Application struct {
	wg     *sync.WaitGroup // app group 관리 wg
	ctx    context.Context
	cancel context.CancelFunc

	ctxsub    context.Context
	cancelSub context.CancelFunc
	Ct        *container.Container

	subRunning bool
	mngRunning bool
	closed     bool
}

func NewApplication(ct *container.Container) *Application {
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	ctxsub, cancelsub := context.WithCancel(context.Background())

	Mlog = logger.Mlog

	wg.Add(1)

	return &Application{
		wg:        wg,
		ctx:       ctx,
		cancel:    cancel,
		ctxsub:    ctxsub,
		cancelSub: cancelsub,
		Ct:        ct, // container
	}
}

func (app *Application) StartProcess() bool {
	// child process start
	Mlog.Print(2, "start process..")
	var isok bool = true

	for idx := 1; idx < PROCESS_CNT; idx++ {
		// ptrPrc := &process[idx]
		ptrPrc := app.Ct.Process[idx]

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
			Mlog.Error("child process start error... state(%X, %X)", prcstate, state)
			isok = false
			break
		}

		idx := (index % PRC_CNT) + 1
		// ptrPrc := &process[idx]
		ptrPrc := app.Ct.Process[idx]

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

func (app *Application) Shutdown() {
	Mlog.Print(5, "Shutdown manage app#1")
	app.cancel()
	<-app.ctx.Done()
	Mlog.Print(5, "Shutdown manage app#2")
}

func (app *Application) ShutDownSubscribe() {
	Mlog.Print(5, "Shutdown rdb sub..#1")
	app.cancelSub()
	<-app.ctxsub.Done()
	Mlog.Print(5, "Shutdown rdb sub..#2")
}

func (app *Application) manageProcess() {
	for idx := 1; idx < PROCESS_CNT; idx++ {

		// ptrPrc := &process[idx]
		ptrPrc := app.Ct.Process[idx]

		prctm := time.Unix(ptrPrc.LastTm, 0)

		// Mlog.Print(2, "child process id %d, %v", ptrPrc.ID, prctm)

		if !ptrPrc.IsExist() {
			Mlog.Warn("UNEXIST [%s %d]", ptrPrc.PrcName, ptrPrc.GetPid())
			Mlog.Warn("Process start [%s]", ptrPrc.PrcName)
			ptrPrc.Start2("-mode=release")
		} else if time.Since(prctm) > app.Ct.Config.PROCESS_INTERVAL {
			Mlog.Warn("RST_ABNOMAL [%s %d]", ptrPrc.PrcName, ptrPrc.GetPid())
			// Mlog.Warn("Processkill [%s %d]", ptrPrc.PrcName, ptrPrc.GetPid())
			// ptrPrc.Kill()
		}
	}
}

func (app *Application) Manage(ctx context.Context) {
	// manage process
	for {
		select {
		case <-ctx.Done():
			Mlog.Print(5, "Done. Mange..")
			return
		default:
			//....
			app.manageProcess()
			time.Sleep(time.Millisecond * 1000)
		}
	}
}

// subscribe app
func (app *Application) subscribeProcessState(ctx context.Context) {
	// subscribe redis value
	for {
		select {
		case <-ctx.Done():
			Mlog.Print(5, "Done. subscribe..")
			return
		default:
			//.... get process
			Mlog.Print(2, "subscribe..")

			raw, err := app.Ct.Rdb.Get(ctx, "prc_beserver")
			if err != nil {
				Mlog.Error("Redis Get Process error: %v", err)
			}

			var result general.Process
			err = json.Unmarshal([]byte(raw), &result)
			if err != nil {
				Mlog.Error("JSON unmarshal error: %v", err)
			}

			fmt.Printf("process Active(%v), ID:%v, name:%v, tm:%v \n", result.Active, result.ID, result.PrcName, result.LastTm)
			app.Ct.Process[1] = result

			//get terminate
			valStr, err2 := app.Ct.Rdb.Get(ctx, "terminate")
			if err2 != nil {
				// log.Fatal(err2)
				Mlog.Error("%v", err2)
			}

			sttTerminate, err2 := strconv.ParseBool(valStr)
			if err2 != nil {
				// log.Fatal("failed to parse bool:", err2)
				Mlog.Error("failed to parse bool:", err2)
			}
			app.Ct.SysInfo.Terminate = sttTerminate
			if app.Ct.SysInfo.Terminate {
				Mlog.Print(4, "Request terminate..")
			}

			time.Sleep(time.Millisecond * 1000)
		}
	}
}

func (app *Application) IsRunningSubscribe() bool {
	return app.subRunning
}

func (app *Application) IsAppClosed() bool {
	return app.closed
}

func (app *Application) clearAppEnv() {

}

func (app *Application) StartManageRoutine(ctx context.Context, mu *sync.Mutex, done chan struct{}) {
	Mlog.Print(2, "StartManageRoutine#1")
	mu.Lock()
	if app.mngRunning {
		mu.Unlock()
		return
	}
	Mlog.Print(2, "StartManageRoutine#2")
	app.mngRunning = true
	mu.Unlock()

	go func() {
		defer func() {
			mu.Lock()
			app.mngRunning = false
			mu.Unlock()
			done <- struct{}{}
		}()
		app.Manage(ctx)
	}()
}

func (app *Application) StartSubscribeRoutine(ctx context.Context, mu *sync.Mutex, done chan struct{}) {
	Mlog.Print(2, "StartSubscribeRoutine#1")
	mu.Lock()
	if app.subRunning {
		mu.Unlock()
		return
	}
	app.subRunning = true
	mu.Unlock()
	Mlog.Print(2, "StartSubscribeRoutine#2")

	go func() {
		defer func() {
			Mlog.Print(2, "Quti StartSubscribeRoutine..#1")
			mu.Lock()
			app.subRunning = false
			mu.Unlock()
			done <- struct{}{}
			Mlog.Print(2, "Quti StartSubscribeRoutine..#2")
		}()
		app.subscribeProcessState(ctx)
	}()
}

func (app *Application) Start() {
	Mlog.Print(2, "Start #1")
	defer app.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			Mlog.Error("Main func panic.. :%v", r)
		}
	}()

	defer func() {
		app.closed = true
	}()
	defer app.clearAppEnv()

	childAppCtx, childAppCancel := context.WithCancel(app.ctx)

	var mu sync.Mutex
	mngDone := make(chan struct{}, 1)
	subDone := make(chan struct{}, 1)

	app.StartSubscribeRoutine(app.ctxsub, &mu, subDone)
	// start child process
	ok := app.StartProcess()
	if !ok {
		Mlog.Error("child process start failed..")
		return
	}

	app.StartManageRoutine(childAppCtx, &mu, mngDone)

	Mlog.Print(2, "Start #2")

	for {
		select {
		case <-app.ctx.Done():
			Mlog.Print(5, "App canceled..")
			childAppCancel() // request cancel to child app
			goto WAIT
		case <-mngDone:
			Mlog.Error("Manage routine ended.. ")
			if app.ctx.Err() == nil {
				Mlog.Print(5, "Restart ManageRoutine..")
				app.StartManageRoutine(childAppCtx, &mu, mngDone)
			}
		case <-subDone:
			Mlog.Error("subscribe routine ended..")
			if app.ctx.Err() == nil {
				Mlog.Print(5, "Restart SubscribeRoutine..")
				app.StartSubscribeRoutine(app.ctxsub, &mu, subDone)
			}
		default:
			time.Sleep(time.Millisecond * 1000)
		}
	}
WAIT:
	done := make(chan struct{})
	go func() {
		for {
			mu.Lock()
			if !app.mngRunning && !app.subRunning {
				mu.Unlock()
				break
			}
			mu.Unlock()
			time.Sleep(time.Millisecond * 100)
		}
		close(done)
	}()

	select {
	case <-done:
		Mlog.Print(2, "Shutdonw main function(1)")
	case <-time.After(time.Second * 5):
		Mlog.Print(2, "Shutdonw main function(2)")
	}
	Mlog.Print(5, "Main server quit..")
}
