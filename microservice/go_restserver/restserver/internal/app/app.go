package app

import (
	"server/restserver/internal/container"
	"server/restserver/internal/logger"
	"server/restserver/internal/server/api"
)

type Application struct {
	ApiServer *api.Server
	// SStServerApp
	// MsgApp
	// container

}

func NewApplication(ct *container.Container) *Application {

	// new httpserver
	// newserver(container)
	apisvr, err := api.NewServer(ct)
	if err != nil {
		logger.Mlog.Error("Api server initialization fail.. %v", err)
		return nil
	}

	// new socket server
	// newserver(container)

	// new msg server
	// newserver(container)

	return &Application{
		ApiServer: apisvr,
	}
}

func (app Application) Start() {
	// start app (api server)
	app.ApiServer.Start()
}

func (app Application) ShutDown() {
	// AppHandler별로
	logger.Mlog.Print(2, "Shutdown Http server#1")
	app.ApiServer.Shutdown()
	logger.Mlog.Print(2, "Shutdown Http server#2")
}

func (app Application) CheckAppState() {

}
