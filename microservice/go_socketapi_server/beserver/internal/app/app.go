package app

import (
	"server/beserver/internal/container"
	"server/beserver/internal/logger"
	"server/beserver/internal/server/api"
	"server/beserver/internal/server/tcpserver"
	"server/general/comm"
	"sync"
)

type Application struct {
	wg        *sync.WaitGroup
	ApiServer *api.Server
	SstServer *tcpserver.SStServer
	// MsgApp
	// container

}

func NewApplication(ct *container.Container) *Application {
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	// new httpserver
	// newserver(container)
	apisvr, err := api.NewServer(wg, ct)
	if err != nil {
		logger.Mlog.Error("Api server initialization fail.. %v", err)
		return nil
	}

	// new socket server
	tcpHnd := comm.NewTcpHandler("client1", 9092, "127.0.0.1")
	sstserver := tcpserver.NewSStServer(wg, 123, &tcpHnd)

	// newserver(container)

	// new msg server
	// newserver(container)

	return &Application{
		wg:        wg,
		ApiServer: apisvr,
		SstServer: sstserver,
	}
}

func (app Application) Start() {
	// start app (api server)
	app.wg.Add(1)
	logger.Mlog.Print(2, "Start API server .. #1")
	go app.ApiServer.Start()

	// start app (sst server)
	app.wg.Add(1)
	logger.Mlog.Print(2, "Start SST server .. #1")
	go app.SstServer.Start()
}

func (app Application) ShutDown() {
	// AppHandler별로
	logger.Mlog.Print(2, "Shutdown Http server#1")
	app.ApiServer.Shutdown()
	logger.Mlog.Print(2, "Shutdown Http server#2")

	logger.Mlog.Print(2, "Shutdown SST server#1")
	app.SstServer.Shutdown()
	logger.Mlog.Print(2, "Shutdown SST server#2")
}

func (app Application) CheckAppState() {

}
