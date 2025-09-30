package mdb

import (
	"auth-service/internal/db"
	"context"
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

func (q *MariaDbHandler) ReadUser(ctx context.Context, name string) (db.USER, error) {
	ado := q.GetDB()

	var u db.USER

	query := `
	select USER_NM, PASSWD, EMAIL, CHG_DT, CREADT_DT from USERS where USER_NM = ?
	`
	rows, err := ado.QueryContext(ctx, query, name)
	if err != nil {
		return u, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(
			&u.USER_NM, &u.PASSWD, &u.EMAIL, &u.CHG_DT, &u.CREATE_DT,
		); err != nil {
			return u, err
		}
	}
	if err := rows.Close(); err != nil {
		return u, err
	}
	if err := rows.Err(); err != nil {
		return u, err
	}
	return u, nil
}
