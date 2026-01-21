package service

import (
	"fmt"

	"ws_test/internal/db"
	"ws_test/internal/service"
)

type ApiService struct {
	dbHnd db.DbHandler
}

func NewApiService(dbHnd db.DbHandler) service.ServiceInterface {
	return &ApiService{
		dbHnd: dbHnd,
	}
}

func (s *ApiService) Test() {
	fmt.Printf("test service")
}
