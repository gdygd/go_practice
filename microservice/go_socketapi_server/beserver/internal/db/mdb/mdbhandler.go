package mdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/godror/godror"
)

func NewMdbHandler(user, pw, dbname, host string, port int) *MariaDbHandler {

	return &MariaDbHandler{user: user, pw: pw, dbNm: dbname, host: host, port: port}
}

type MariaDbHandler struct {
	db   *sql.DB
	user string
	pw   string
	dbNm string
	host string
	port int
}

func (q *MariaDbHandler) Open() (*sql.DB, error) {

	dbSrc := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Local", q.user, q.pw, q.host, q.dbNm)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db, err := sql.Open("mysql", dbSrc)
	if err != nil {
		return nil, err
	}

	// PingContext로 연결 확인
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	fmt.Println("db open")
	q.db = db

	return db, nil
}

func (q *MariaDbHandler) Close(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
