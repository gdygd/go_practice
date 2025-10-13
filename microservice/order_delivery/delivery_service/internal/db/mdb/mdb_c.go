package mdb

import (
	"context"
	"delivery_service/internal/db"
	"delivery_service/internal/logger"
)

func (q *MariaDbHandler) RequestDelivery(ctx context.Context, delivery db.DELIVERIES) (db.DELIVERIES, error) {
	ado := q.GetDB()

	var deli db.DELIVERIES

	query := `
	INSERT INTO deliveries (ORDER_ID, STATUS, ADDRESS, REQ_DT)
	VALUES(?, ?, ?, now()) RETURNING DELIVERY_ID, ORDER_ID, STATUS, ADDRESS, REQ_DT
	`

	row := ado.QueryRow(query, delivery.ORDER_ID, db.PENDING, delivery.ADDRESS)

	err := row.Scan(&deli.DELIVERY_ID, &deli.ORDER_ID, &deli.STATUS, &deli.ADDRESS, &deli.REQ_DT)
	if err != nil {
		logger.Log.Error("RequestOrder error. %v", err)
		return deli, err
	}

	return deli, nil
}
