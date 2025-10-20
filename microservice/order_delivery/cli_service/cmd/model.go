package cmd

import "time"

type Service struct {
	Name        string
	Url         string
	ServicePath string
	RestartCmd  string
	FailCount   int
	LastRestart time.Time
	IsManaged   bool
}

var services []*Service

const (
	SERVICE_BASE_PATH = "/home/gildong/work/go/src/go_practice/microservice/order_delivery"
	BASE_URL          = "http://10.1.0.119"
)

func InitServices() {
	services = []*Service{
		{"api-gateway", BASE_URL + ":9080", SERVICE_BASE_PATH + "/api_gateway/bin/", "./api-gateway", 0, time.Time{}, true},
		{"auth-service", BASE_URL + ":9081", SERVICE_BASE_PATH + "/auth_service/bin/", "./auth-service", 0, time.Time{}, true},
		{"order-service", BASE_URL + ":9082", SERVICE_BASE_PATH + "/order_service/bin/", "./order-service", 0, time.Time{}, true},
		{"delivery_service", BASE_URL + ":9083", SERVICE_BASE_PATH + "/delivery_service/bin/", "./delivery-service", 0, time.Time{}, true},
		{"saga_service", BASE_URL + ":9084", SERVICE_BASE_PATH + "/orchestrator/bin/", "./saga-service", 0, time.Time{}, true},
	}
}

type StateCommand struct {
	serviceName string
}

type ResetCommand struct {
	serviceName string
}

type DebugCommand struct {
	serviceName string
	level       int
}
