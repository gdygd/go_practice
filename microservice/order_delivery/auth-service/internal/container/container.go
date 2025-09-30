package container

import (
	"auth-service/internal/config"
	"auth-service/internal/db"
	"auth-service/internal/memory"
)

type Container struct {
	Config *config.Config
	DbHnd  db.DbHandler
	ObjDb  *memory.RedisDb
}

var container *Container

func NewContainer() (*Container, error) {

	return nil, nil
}

func initConfig() (config.Config, error) {

	return config.LoadConfig(".")
}
