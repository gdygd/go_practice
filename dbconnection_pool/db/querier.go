package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	Init() (*sql.DB, error)
	Close(*sql.DB)
	ReadSysdate(ctx context.Context) (string, error)
}

// var _ Querier = (*Queries)(nil)
