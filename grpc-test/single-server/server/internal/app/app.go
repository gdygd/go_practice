package app

import (
	"sync"

	"grpc_svr_test/internal/container"
	"grpc_svr_test/internal/logger"
	"grpc_svr_test/internal/server/api"
	gapi "grpc_svr_test/internal/server/grpc"
)

type Application struct {
	wg        *sync.WaitGroup
	ApiServer *api.Server
	GServer   *gapi.Server
}

func NewApplication(ct *container.Container, ch_terminate chan bool) *Application {
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	// new httpserver
	// newserver(container)
	apisvr, err := api.NewServer(wg, ct, ch_terminate)
	if err != nil {
		logger.Log.Error("Api server initialization fail.. %v", err)
		return nil
	}

	gserver, err := gapi.NewServer(wg, ct, ch_terminate)
	if err != nil {
		logger.Log.Error("Api server initialization fail.. %v", err)
		return nil
	}

	return &Application{
		wg:        wg,
		ApiServer: apisvr,
		GServer:   gserver,
	}
}

func (app Application) Start() {
	app.wg.Add(1)
	logger.Log.Print(3, "Start API server.. #1")
	go app.ApiServer.Start()

	app.wg.Add(1)
	logger.Log.Print(3, "Start gRPC server.. #1")
	go app.GServer.StartgPRC()
}

func (app Application) Shutdown() {
	logger.Log.Print(3, "Shutdown Rest server#1")
	app.ApiServer.Shutdown()
	logger.Log.Print(3, "Shutdown Rest server#2")

	logger.Log.Print(3, "Shutdown gRPC server#1")
	app.GServer.ShutdowngRPC()
	logger.Log.Print(3, "Shutdown gRPC server#2")
}
