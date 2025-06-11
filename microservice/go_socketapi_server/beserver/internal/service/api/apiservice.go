package service

import (
	"fmt"
	"server/beserver/internal/db"
	"server/beserver/internal/obj"
	"server/beserver/internal/service"
)

type ApiService struct {
	dbHnd db.DbHandler
	objdb *obj.ObjectDb
	// socHnd
	// obj
}

func NewApiService(dbHnd db.DbHandler, objdb *obj.ObjectDb) service.ApiServiceInterface {
	return &ApiService{
		dbHnd: dbHnd,
		objdb: objdb,
	}
}

func (s *ApiService) Test() {
	fmt.Printf("test service")
}
