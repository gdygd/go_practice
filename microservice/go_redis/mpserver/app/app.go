package app

import (
	"go_redis/beserver/internal/container"
	"go_redis/general"
	"go_redis/general/cache"
	"sync"
)

var Mlog *general.OLog2

type Application struct {
	wg        *sync.WaitGroup
	Ct        *container.Container
	mangeDone chan struct{}
	subDone   chan struct{}
}

func NewApplication(ct *container.Container) *Application {
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	Mlog = ct.Mlog

	var mangeDone chan struct{}
	var subDone chan struct{}

	// app add
	wg.Add(2) // manageprocess, subscribe

	go subscribeProcessState(wg, ct.Rdb, subDone)

	return &Application{
		wg:        wg,
		Ct:        ct,
		mangeDone: mangeDone,
		subDone:   subDone,
	}
}

func (app *Application) StartProcess() {
	// child process start

}

func (app *Application) ShutDownManagePrcoess() {
	app.mangeDone <- struct{}{}
}

func (app *Application) ShutDownSubscribe() {
	app.subDone <- struct{}{}
}

func (app *Application) ManageProcess(wg *sync.WaitGroup) {
	// manage child process running state
	var isDone bool = false

	defer app.wg.Done()

	for {
		select {
		case <-app.mangeDone:
			isDone = true
			Mlog.Print(2, "Done MangeProces...")
			break
		default:

		}

		if isDone {
			Mlog.Print(2, "close ManageProcess..")
			break
		}
	}
}

// subscribe app
func subscribeProcessState(wg *sync.WaitGroup, rdb *cache.RedisClient, subDone chan struct{}) {
	// update child process state
	//ct.Process[1]

	var isDone bool = false
	defer wg.Done()

	for {
		select {
		case <-subDone:
			Mlog.Print(2, "Done MangeProces...")
		default:

		}
		if isDone {
			Mlog.Print(2, "close subscribeProcessState..")
			break
		}
	}

}
