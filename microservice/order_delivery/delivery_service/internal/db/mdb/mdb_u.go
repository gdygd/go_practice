package mdb

import (
	"context"
	"delivery_service/internal/db"
	"delivery_service/internal/logger"
)

func (q *MariaDbHandler) CancelDelivery(ctx context.Context, orderId int) error {
	ado := q.GetDB()

	query := `
	UPDATE deliveries SET STATUS = ? where ORDER_ID = ?
	`
	_, err := ado.ExecContext(ctx, query, db.CANCELLED, orderId)
	if err != nil {
		logger.Log.Error("CancelOrder error. %v", err)
		return err
	}

	return nil
}

func (q *MariaDbHandler) ConfirmDelivery(ctx context.Context, orderId int) error {
	ado := q.GetDB()

	query := `
	UPDATE deliveries SET STATUS = ? where ORDER_ID = ?	
	`

	_, err := ado.ExecContext(ctx, query, db.CONFIRMED, orderId)
	if err != nil {
		logger.Log.Error("CancelOrder error. %v", err)
		return err
	}

	return nil
}
