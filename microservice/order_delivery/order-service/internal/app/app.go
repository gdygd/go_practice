package app

import (
	"order-service/internal/container"
	"order-service/internal/logger"
	"order-service/internal/server/api"
	"sync"
)

type Application struct {
	wg        *sync.WaitGroup
	ApiServer *api.Server
}

func NewApplication(ct *container.Container) *Application {
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	// new httpserver
	// newserver(container)
	apisvr, err := api.NewServer(wg, ct)
	if err != nil {
		logger.Log.Error("Api server initialization fail.. %v", err)
		return nil
	}

	return &Application{
		wg:        wg,
		ApiServer: apisvr,
	}
}

func (app Application) Start() {
	app.wg.Add(1)
	logger.Log.Print(3, "Start API server.. #1")
	go app.ApiServer.Start()
}

func (app Application) Shutdown() {
	logger.Log.Print(3, "Shutdown Rest server#1")
	app.ApiServer.Shutdown()
	logger.Log.Print(3, "Shutdown Rest server#2")
}
