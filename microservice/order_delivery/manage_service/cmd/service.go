package main

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"manage-service/internal/logger"

	"github.com/gdygd/goglib"
)

type Service struct {
	Name        string
	Url         string
	ServicePath string
	RestartCmd  string
	FailCount   int
	MaxFail     int
	LastRestart time.Time
	IsManaged   bool
}

const MAX_FAIL = 3

var services []*Service

const (
	SERVICE_BASE_PATH = "/home/gildong/work/go/src/go_practice/microservice/order_delivery"
	BASE_URL          = "http://10.1.0.119"
)

func initServices() {
	services = []*Service{
		{"api-gateway", BASE_URL + ":9080", SERVICE_BASE_PATH + "/api_gateway/bin/", "./api-gateway", 0, MAX_FAIL, time.Time{}, true},
		{"auth-service", BASE_URL + ":9081", SERVICE_BASE_PATH + "/auth_service/bin/", "./auth-service", 0, MAX_FAIL, time.Time{}, true},
		{"order-service", BASE_URL + ":9082", SERVICE_BASE_PATH + "/order_service/bin/", "./order-service", 0, MAX_FAIL, time.Time{}, true},
		{"delivery_service", BASE_URL + ":9083", SERVICE_BASE_PATH + "/delivery_service/bin/", "./delivery-service", 0, MAX_FAIL, time.Time{}, true},
		{"saga_service", BASE_URL + ":9084", SERVICE_BASE_PATH + "/orchestrator/bin/", "./saga-service", 0, MAX_FAIL, time.Time{}, true},
	}
}

func (s *Service) CheckHeartbeat() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// HttpRequest안에서 timeout처리가 됨.
	statuscode, _, err := goglib.HttpRequest(ctx, http.Header{}, nil, "GET", s.Url+"/heartbeat")
	if err != nil || statuscode != http.StatusOK {
		logger.Log.Error("Check Heartbeat fail.. [%s]", s.Name)
		s.FailCount++
		return false
	}

	logger.Log.Print(2, "CheckHearbeat ok..[%s]", s.Name)

	s.FailCount = 0
	return true
}

func (s *Service) Restart() {
	logger.Log.Print(2, "Restart service [%s] path :[%s] cmd : %s", s.Name, s.ServicePath, s.RestartCmd)

	cmd := exec.Command(s.RestartCmd, "&")

	// 세션 리더로 만들어 부모 프로세스와 분리
	// if runtime.GOOS == "windows" {
	// 	cmd.SysProcAttr = &syscall.SysProcAttr{
	// 		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	// 	}
	// } else {
	// 	cmd.SysProcAttr = &syscall.SysProcAttr{
	// 		Setsid: true,
	// 	}
	// }
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = s.ServicePath
	err := cmd.Start()
	if err != nil {
		logger.Log.Error("Restart Service failed.. %s", s.Name)
	}

	// go func() {
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		logger.Log.Error("Restart Service failed.. %s", s.Name)
	// 	}
	// }()

	s.LastRestart = time.Now()
	s.FailCount = 0
}

func (s *Service) Terminate() {
	logger.Log.Print(2, "Terminate service [%s]", s.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// HttpRequest안에서 timeout처리가 됨.
	statuscode, _, err := goglib.HttpRequest(ctx, http.Header{}, nil, "GET", s.Url+"/terminate")
	if err != nil || statuscode != http.StatusOK {
		logger.Log.Error("Terminate fail.. [%s] (%v)", s.Name, err)
	} else {
		logger.Log.Print(2, "Terminate ok..[%s]", s.Name)
	}
}

func TerminateServices() {
	logger.Log.Print(2, "TerminateServices..")
	for _, serv := range services {
		logger.Log.Print(2, "TerminateService.. [%s]", serv.Name)
		serv.Terminate()
	}
}

func CheckServices() {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		for _, serv := range services {
			if !serv.IsManaged {
				continue
			}

			logger.Log.Print(1, "check service : %s [%d/%d]", serv.Name, serv.FailCount, serv.MaxFail)
			if !serv.CheckHeartbeat() {
				logger.Log.Warn("[%s] health check failed (%d/%d)", serv.Name, serv.FailCount, serv.MaxFail)

				if serv.FailCount >= serv.MaxFail {
					serv.Restart()
				}

			}
		}
	}
}
