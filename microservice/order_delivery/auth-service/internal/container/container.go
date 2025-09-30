package container

import (
	"auth-service/internal/config"
	"auth-service/internal/db"
	"auth-service/internal/db/mdb"
	"auth-service/internal/memory"
	"fmt"
)

type Container struct {
	Config *config.Config
	DbHnd  db.DbHandler
	ObjDb  *memory.RedisDb
}

var container *Container

func NewContainer() (*Container, error) {

	container = &Container{}
	// load config
	config, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("config loading error..%v \n", err)
	}
	container.Config = &config

	// init database
	dbhnd := initDatabase(config)
	container.DbHnd = dbhnd

	// init object db
	obj := memory.InitRedisDb(config.RedisAddr)
	container.ObjDb = obj

	return container, nil
}

func initConfig() (config.Config, error) {

	return config.LoadConfig(".")
}

func initDatabase(config config.Config) db.DbHandler {
	mdb := mdb.NewMdbHandler(config.DBUser, config.DBPasswd, config.DBSName, config.DBAddress, config.DBPort)
	return mdb
}
