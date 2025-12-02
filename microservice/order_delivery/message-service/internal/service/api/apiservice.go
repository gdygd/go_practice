package service

import (
	"fmt"

	"message_service/internal/db"
	"message_service/internal/memory"
	"message_service/internal/service"
)

type ApiService struct {
	dbHnd db.DbHandler
	objdb *memory.RedisDb
}

func NewApiService(dbHnd db.DbHandler, objdb *memory.RedisDb) service.ServiceInterface {
	return &ApiService{
		dbHnd: dbHnd,
		objdb: objdb,
	}
}

func (s *ApiService) Test() {
	fmt.Printf("test service")
}
