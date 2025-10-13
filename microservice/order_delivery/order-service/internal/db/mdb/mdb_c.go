package mdb

import (
	"context"
	"order-service/internal/db"
	"order-service/internal/logger"
)

func (q *MariaDbHandler) RequestOrder(ctx context.Context, ord db.ORDER) (db.ORDER, error) {
	ado := q.GetDB()

	var order db.ORDER

	query := `
	INSERT INTO orders (USER_NM, STATE, ORDER_DT, TOT_AMOUNT)
	VALUES(?, ?, now(), ?) RETURNING ORDER_ID, USER_NM, STATE, ORDER_DT, TOT_AMOUNT
	`

	row := ado.QueryRow(query, ord.USER_NM, db.PENDING, ord.TOT_AMOUNT)

	err := row.Scan(&order.ORDER_ID, &order.USER_NM, &order.STATE, &order.ORDER_DT, &order.TOT_AMOUNT)
	if err != nil {
		logger.Log.Error("RequestOrder error. %v", err)
		return order, err
	}

	return order, nil
}
