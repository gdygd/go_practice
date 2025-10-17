package db

import "time"

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
