package db

import (
	"database/sql"
	"time"
)

type USER struct {
	USER_NM   string
	PASSWD    string
	EMAIL     string
	CHG_DT    time.Time
	CREATE_DT time.Time
}

type ORDER struct {
	ORDER_ID   int
	USER_NM    string
	STATE      int
	ORDER_DT   time.Time
	TOT_AMOUNT int
}

type DELIVERIES struct {
	DELIVERY_ID int
	ORDER_ID    int
	STATUS      int
	ADDRESS     string
	REQ_DT      time.Time
	COMPL_DT    sql.NullTime
}
