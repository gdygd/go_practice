package db

import (
	"context"
	"database/sql"
)

type DbHandler interface {
	Init() error
	Close(*sql.DB)
	ReadSysdate(ctx context.Context) (string, error)
	ReadUser(ctx context.Context, name string) (USER, error)
}
