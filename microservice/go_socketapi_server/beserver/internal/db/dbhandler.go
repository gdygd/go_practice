package db

import (
	"context"
	"database/sql"
)

type DbHandler interface {
	Open() (*sql.DB, error)
	Close(*sql.DB)
	ReadTest(ctx context.Context) (int, error)
}
