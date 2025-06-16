package service

import (
	"fmt"
	"go_redis/beserver/internal/db"
	"go_redis/beserver/internal/memory"
	"go_redis/beserver/internal/service"
)

type ApiService struct {
	dbHnd db.DbHandler
	objdb *memory.RedisDb
	// socHnd
	// obj
}

func NewApiService(dbHnd db.DbHandler, objdb *memory.RedisDb) service.ApiServiceInterface {
	return &ApiService{
		dbHnd: dbHnd,
		objdb: objdb,
	}
}

func (s *ApiService) Test() {
	fmt.Printf("test service")
}
