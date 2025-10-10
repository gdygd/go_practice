package mdb

import (
	"context"
	"order-service/internal/db"
	"order-service/internal/logger"
)

func (q *MariaDbHandler) ReadSysdate(ctx context.Context) (string, error) {
	ado := q.GetDB()

	query := `
	select now() as dt from dual
	`

	rows, err := ado.QueryContext(ctx, query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	strDateTime := ""
	if rows.Next() {
		if err := rows.Scan(
			&strDateTime,
		); err != nil {
			return "", err
		}
	}
	if err := rows.Close(); err != nil {
		return "", err
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return strDateTime, nil
}

func (q *MariaDbHandler) ReadOrderInfo(ctx context.Context, username string) ([]db.ORDER, error) {
	ado := q.GetDB()

	query := `
	SELECT a.ORDER_ID, a.USER_NM, a.STATE, a.ORDER_DT, a.TOT_AMOUNT 
	FROM orders a
	WHERE a.USER_NM = ?
	ORDER BY a.ORDER_DT DESC
	`

	rows, err := ado.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []db.ORDER = []db.ORDER{}
	for rows.Next() {
		ord := db.ORDER{}

		err := rows.Scan(&ord.ORDER_ID, &ord.USER_NM, &ord.STATE, &ord.ORDER_DT, &ord.TOT_AMOUNT)
		if err != nil {
			logger.Log.Error("ReadOrderInfo Scan fail..(%v)", err)
			return nil, err
		}

		orders = append(orders, ord)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
