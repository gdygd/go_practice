package mdb

import (
	"context"
	"delivery_service/internal/db"
	"delivery_service/internal/logger"
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

func (q *MariaDbHandler) ReadDeliveries(ctx context.Context, username string) ([]db.DELIVERIES, error) {
	ado := q.GetDB()

	query := `
	SELECT b.DELIVERY_ID, b.ORDER_ID, b.STATUS, b.ADDRESS, b.REQ_DT, b.COMPL_DT
	from orders a
	inner JOIN deliveries b
	ON a.ORDER_ID = b.ORDER_ID
	WHERE a.USER_NM = ?
	ORDER BY b.REQ_DT desc
	`

	rows, err := ado.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []db.DELIVERIES = []db.DELIVERIES{}
	for rows.Next() {
		deli := db.DELIVERIES{}

		err := rows.Scan(&deli.DELIVERY_ID, &deli.ORDER_ID, &deli.STATUS, &deli.ADDRESS, &deli.REQ_DT, &deli.COMPL_DT)
		if err != nil {
			logger.Log.Error("ReadOrderInfo Scan fail..(%v)", err)
			return nil, err
		}

		deliveries = append(deliveries, deli)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return deliveries, nil
}
