package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DBHandler struct {
	Host     string
	Port     int
	DBname   string
	User     string
	Password string
}

func (dbHand *DBHandler) Open() (*sql.DB, error) {
	dbInfo := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		dbHand.User, dbHand.Password, dbHand.Host, dbHand.Port, dbHand.DBname,
	)

	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (dbHand *DBHandler) Close(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		return err
	}

	return nil
}
