package service

import (
	"fmt"
	"order-delivery-saga/internal/db"
	"order-delivery-saga/internal/memory"
	"order-delivery-saga/internal/service"
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
