package service

import (
	"fmt"
	"server/restserver/internal/db"
	"server/restserver/internal/obj"
	"server/restserver/internal/service"
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
