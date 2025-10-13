package db

import (
	"context"
	"database/sql"
)

const (
	PENDING   = 1 // 대기
	CONFIRMED = 2 // 확정
	CANCELLED = 3 // 취소
)

type DbHandler interface {
	Init() error
	Close(*sql.DB)
	ReadSysdate(ctx context.Context) (string, error)

	ReadDeliveries(ctx context.Context, username string) ([]DELIVERIES, error)

	RequestDelivery(ctx context.Context, deli DELIVERIES) (DELIVERIES, error)
	CancelDelivery(ctx context.Context, orderId int) error
	ConfirmDelivery(ctx context.Context, orderId int) error
}
