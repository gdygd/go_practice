package app

import (
	"sync"

	gapi "grpc_client_test/internal/client/grpc"
	"grpc_client_test/internal/container"
	"grpc_client_test/internal/logger"
	"grpc_client_test/internal/msgproc"
	"grpc_client_test/internal/server/api"
)

type Application struct {
	wg         *sync.WaitGroup
	ApiServer  *api.Server
	Gclient    *gapi.GrpcClient
	MsgHandler *msgproc.MsgProcHandler
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

	gclient, _ := gapi.NewClient(wg, ct, ch_terminate)
	// if err != nil {
	// 	logger.Log.Error("Api server initialization fail.. %v", err)
	// 	return nil
	// }

	msgHandler, _ := msgproc.NewMsgProcHandler(wg, ct)

	return &Application{
		wg:         wg,
		ApiServer:  apisvr,
		Gclient:    gclient,
		MsgHandler: msgHandler,
	}
}

func (app Application) Start() {
	app.wg.Add(1)
	logger.Log.Print(3, "Start API server.. #1")
	go app.ApiServer.Start()

	app.wg.Add(1)
	logger.Log.Print(3, "Start gRPC client.. #1")
	go app.Gclient.Start()

	app.wg.Add(1)
	logger.Log.Print(3, "Start MsgHandler.. #1")
	go app.MsgHandler.Start()
}

func (app Application) Shutdown() {
	logger.Log.Print(3, "Shutdown Rest server#1")
	go app.ApiServer.Shutdown()
	logger.Log.Print(3, "Shutdown Rest server#2")

	logger.Log.Print(3, "Shutdown grpc client#1")
	go app.Gclient.Shutdown()
	logger.Log.Print(3, "Shutdown grpc client#2")

	logger.Log.Print(3, "Shutdown MsgHandler#1")
	go app.MsgHandler.Shutdown()
	logger.Log.Print(3, "Shutdown MsgHandler#2")

	app.wg.Wait()
}
