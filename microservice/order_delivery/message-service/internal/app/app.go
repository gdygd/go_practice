package app

import (
	"sync"

	"message_service/internal/container"
	"message_service/internal/logger"
	"message_service/internal/server/api"
	"message_service/internal/server/msgq"
)

type Application struct {
	wg           *sync.WaitGroup
	ApiServer    *api.Server
	RabbitClient *msgq.RabbitMQClient
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

	rabbitClient, err := msgq.NewClient(wg, ct)
	if err != nil {
		logger.Log.Error("RabbitMQ Client initialization fail.. %v", err)
		return nil
	}

	return &Application{
		wg:           wg,
		ApiServer:    apisvr,
		RabbitClient: rabbitClient,
	}
}

func (app Application) Start() {
	app.wg.Add(1)
	logger.Log.Print(3, "Start API server.. #1")
	go app.ApiServer.Start()

	app.wg.Add(1)
	logger.Log.Print(3, "Start Rabbitmq client.. #1")
	go app.RabbitClient.Start()
}

func (app Application) Shutdown() {
	logger.Log.Print(3, "Shutdown Rest server#1")
	app.ApiServer.Shutdown()
	logger.Log.Print(3, "Shutdown Rest server#2")

	logger.Log.Print(3, "Shutdown RabbitMQ client#1")
	app.RabbitClient.Shutdown()
	logger.Log.Print(3, "Shutdown RabbitMQ client#2")
}
