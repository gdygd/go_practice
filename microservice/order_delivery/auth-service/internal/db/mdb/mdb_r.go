package mdb

import (
	"auth-service/internal/db"
	"auth-service/internal/logger"
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
	logger.Log.Print(2, "ReadUser...#1")
	ado := q.GetDB()

	logger.Log.Print(2, "ReadUser...#2")

	var u db.USER

	query := `
	select USER_NM, PASSWD, EMAIL, CHG_DT, CREATE_DT from USERS where USER_NM = ?
	`
	rows, err := ado.QueryContext(ctx, query, name)
	if err != nil {
		return u, err
	}
	defer rows.Close()

	logger.Log.Print(2, "ReadUser...#3")
	if rows.Next() {
		if err := rows.Scan(
			&u.USER_NM, &u.PASSWD, &u.EMAIL, &u.CHG_DT, &u.CREATE_DT,
		); err != nil {
			return u, err
		}
	}
	logger.Log.Print(2, "ReadUser...#4")
	if err := rows.Close(); err != nil {
		return u, err
	}
	if err := rows.Err(); err != nil {
		return u, err
	}

	logger.Log.Print(2, "ReadUser...#5")
	return u, nil
}

func (q *MariaDbHandler) ReadSession(ctx context.Context, id string) (db.SESSIONS, error) {
	ado := q.GetDB()

	var se db.SESSIONS

	query := `
	SELECT ID, USER_NM, REF_TOKEN, USER_AGENT, CLIENT_IP, BLOCK_YN, EXP_DT, CREATE_DT FROM sessions a
	WHERE ID = ?
	`

	rows, err := ado.QueryContext(ctx, query, id)
	if err != nil {
		return se, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(
			&se.ID, &se.USER_NM, &se.REF_TOKEN, &se.USER_AGENT, &se.CLIENT_IP, &se.BLOCK_YN, &se.EXP_DT, &se.CREATE_DT,
		); err != nil {
			return se, err
		}
	}
	if err := rows.Close(); err != nil {
		return se, err
	}
	if err := rows.Err(); err != nil {
		return se, err
	}
	return se, nil
}
