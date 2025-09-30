package service

import (
	"auth-service/internal/db"
	"auth-service/internal/memory"
	"auth-service/internal/service"
	"fmt"
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
