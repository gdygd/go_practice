package memory

import (
	"context"
	"go_redis/beserver/internal/cache"
	"go_redis/beserver/internal/logger"
	"go_redis/general"
	"sync"
	"time"
)

type RedisDb struct {
	ctx    context.Context
	Rdb    *cache.RedisClient
	User   map[string]User
	SysStt SystemInfo
	usermu sync.Mutex
	sysmu  sync.Mutex
}

func InitRedisDb() *RedisDb {
	// init database

	RedisDb := &RedisDb{
		ctx:    context.Background(),
		Rdb:    cache.NewRedisClient("127.0.0.1:6379"),
		User:   make(map[string]User),
		SysStt: SystemInfo{SvrUtc: 0, DbSvrComm: STATE_OFF},
	}

	return RedisDb
}

func (r *RedisDb) setUser(id string, value User) {
	r.usermu.Lock()
	defer r.usermu.Unlock()

	r.User[id] = value
}

func (r *RedisDb) setServerUtc(utc int64) {
	r.sysmu.Lock()
	defer r.sysmu.Unlock()

	r.SysStt.SvrUtc = utc
}

func (r *RedisDb) SetUser(id string, value User) {

	if r.Rdb.Set(r.ctx, id, value, 0) != nil {
		logger.Mlog.Error("Redis SetUser error. Key(%d), v:%v", id, value)
	} else {
		r.setUser(id, value)
	}
}

func (r *RedisDb) GetUser(id string) User {
	return User{}
}

func (r *RedisDb) SetServerUtc() {
	curtm := time.Now().Unix()
	if err := r.Rdb.Set(r.ctx, "svrutc", curtm, 0); err != nil {
		logger.Mlog.Error("Redis SetServerUtc error ")

	} else {
		r.setServerUtc(curtm)
	}
}

func (r *RedisDb) SetProcess(prc general.Process) {

	if err := r.Rdb.Set(r.ctx, "prc_beserver", prc, 0); err != nil {
		logger.Mlog.Error("Redis SetProcess error ")
	}
}
