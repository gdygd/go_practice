package mdb

import (
	"auth-service/internal/db"
	"auth-service/internal/logger"
	"context"
)

func (q *MariaDbHandler) CreateSession(ctx context.Context, ss db.SESSIONS) (SESSIONS, error) {
	ado := q.GetDB()

	var se db.SESSIONS

	query := `
	INSERT INTO sessions (ID, USER_NM, REF_TOKEN, USER_AGENT, CLIENT_IP, BLOCK_YN, EXP_DT, CREATE_DT)
	VALUES(?, ?, ?, ?, ?, ?, ?, now()) RETURNING ID, USER_NM, REF_TOKEN, USER_AGENT, CLIENT_IP, BLOCK_YN, EXP_DT, CREATE_DT
	`

	row := ado.QueryRow(query, ss.ID, ss.USER_NM, ss.REF_TOKEN, ss.USER_AGENT, ss.CLIENT_IP, ss.BLOCK_YN, ss.EXP_DT)

	err := row.Scan(&se.ID, &se.USER_NM, &se.REF_TOKEN, &se.USER_AGENT, &se.CLIENT_IP, &se.BLOCK_YN, &se.EXP_DT, &se.CREATE_DT)
	if err != nil {
		logger.Log.Error("CreateSession error. %v", err)
	}

	return se, err
}
