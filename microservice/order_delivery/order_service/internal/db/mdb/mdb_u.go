package mdb

import (
	"context"
	"order-service/internal/db"
	"order-service/internal/logger"
)

func (q *MariaDbHandler) CancelOrder(ctx context.Context, orderId int) error {
	ado := q.GetDB()

	query := `
	UPDATE orders SET STATE = ? where ORDER_ID = ?
	`
	_, err := ado.ExecContext(ctx, query, db.CANCELLED, orderId)
	if err != nil {
		logger.Log.Error("CancelOrder error. %v", err)
		return err
	}

	return nil
}

func (q *MariaDbHandler) ConfirmOrder(ctx context.Context, orderId int) error {
	ado := q.GetDB()

	query := `
	UPDATE orders SET STATE = ? where ORDER_ID = ?	
	`

	_, err := ado.ExecContext(ctx, query, db.CONFIRMED, orderId)
	if err != nil {
		logger.Log.Error("CancelOrder error. %v", err)
		return err
	}

	return nil
}
