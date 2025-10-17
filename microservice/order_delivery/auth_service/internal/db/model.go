package db

import "time"

type USER struct {
	USER_NM   string
	PASSWD    string
	EMAIL     string
	CHG_DT    time.Time
	CREATE_DT time.Time
}

type SESSIONS struct {
	ID         string
	USER_NM    string
	REF_TOKEN  string
	USER_AGENT string
	CLIENT_IP  string
	BLOCK_YN   int
	EXP_DT     time.Time
	CREATE_DT  time.Time
}
