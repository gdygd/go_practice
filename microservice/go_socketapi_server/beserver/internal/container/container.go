package container

import (
	"fmt"
	"server/beserver/internal/db"
	"server/beserver/internal/db/mdb"
	"server/beserver/internal/obj"
	"server/config"
	"server/shmobj"
)

// config
// db
// obj
type Container struct {
	Config    *config.Config
	DbHnd     db.DbHandler
	ObjDb     *obj.ObjectDb
	SharedMem *shmobj.SharedMemory
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
	obj := obj.InitObjectDb()
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

func initObjectDb() *obj.ObjectDb {

	return obj.InitObjectDb()
}
