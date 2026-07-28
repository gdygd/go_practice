package main

import (
	"log"
	"tct/app/am"
	"tct/app/dbapp/sqldb"
)

const (
	ORACLE  = "oracle"
	TIBERO  = "tibero"
	MARIADB = "mariadb"
)

func nothing(sqlHandler *sqldb.SQLHandler) {
	// 아무것도 하지 않음
	log.Println("NOTHING !!")
}

func selectDBMS(dbms string) *sqldb.SQLHandler {
	var sqlHandler = &sqldb.SQLHandler{}

	switch dbms {
	case ORACLE:
		// oracle
		ouser := "UWSIG_TEST"
		opw := "UWSIG_TEST"
		odbNm := "theroadora11"
		// ohost:= "10.1.44.121" // invalid ip address
		ohost := "192.168.2.161" // valid ip address
		oport := 1521
		odbver := "1.0"
		odbms := ORACLE
		odriverNm := ""
		odsn := ""

		sqlHandler = sqldb.NewSQLHandler(ouser, opw, odbNm, ohost, oport, odbver, odbms, odriverNm, odsn)

	case TIBERO:
		// tibero
		tdriverNm := "Tibero 7 ODBC driver"
		tuser := "SIGUSR_TEST"
		tpw := "SIGUSR_TEST"
		tdbNm := "tibero"
		// thost:=     "10.1.44.121" // invalid ip address
		thost := "10.1.0.120" // valid ip address
		tport := 8629
		tdbver := "2.0"
		tdbms := TIBERO
		tdsn := "tbodbc"

		sqlHandler = sqldb.NewSQLHandler(tuser, tpw, tdbNm, thost, tport, tdbver, tdbms, tdriverNm, tdsn)

	case MARIADB:
		// mariadb
		muser := "dev"
		mpw := "dev"
		mdbNm := "sigusr_test"
		// mhost:= "10.1.44.121" // invalid ip address
		mhost := "10.1.0.114" // valid ip address
		mport := 3306
		mdbver := "2.0"
		mdbms := MARIADB
		mdriverNm := ""
		mdsn := ""

		sqlHandler = sqldb.NewSQLHandler(muser, mpw, mdbNm, mhost, mport, mdbver, mdbms, mdriverNm, mdsn)
	}

	return sqlHandler
}

func main() {
	// 밀리초와 마이크로초까지 출력되도록 설정
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// var sqlHandler = selectDBMS(ORACLE)
	var sqlHandler = selectDBMS(TIBERO)
	// var sqlHandler = selectDBMS(MARIADB)
	log.Printf("[sqlHandler] ::: %+v\n", sqlHandler)
	// nothing(sqlHandler)

	// 파라미터 값 설정
	var evtLogParam am.EventLogParam // request 파라미터 변환해서 db에 전달
	evtLogParam.SDt = "20220801000000"
	// evtLogParam.SDt = "20240801000000"
	evtLogParam.EDt = "20241201235959"
	evtLogParam.Page = 0
	evtLogParam.ReqRows = 200
	evtLogParam.SearchTp = 1
	evtLogParam.Offset = 0

	evtLogParam.UType = 0
	evtLogParam.UId = 0
	evtLogParam.Ctor = 0

	evtLogParam.EvCode = 0
	evtLogParam.EvLowCodes = []int{}

	// bool값 설정
	evtLogParam.SDtOk = true
	evtLogParam.EDtOk = true
	evtLogParam.PageOk = true
	evtLogParam.ReqRowsOk = true
	evtLogParam.SearchTpOk = true

	evtLogParam.UTypeOk = false
	evtLogParam.UIdOk = false
	evtLogParam.CtorOk = false

	evtLogParam.EvCodeOk = false
	evtLogParam.EvLowCodeOk = false

	// 다른 필드로 sort
	// evtLogParam.SortFields = [{Field:uType Order:0} {Field:code Order:0}]

	eventLogList, rowcnt, err := sqlHandler.ReadEventLogExceptOprt(evtLogParam)
	if err != nil {
		log.Printf("[main, ReadEventLogExceptOprt error] : %s \n", err.Error())
	}

	// log.Printf("[rowcnt] ::: %v\n", rowcnt)
	// fmt.Println()
	// log.Printf("[eventLogList] ::: %v\n", eventLogList)
	log.Printf("[rowcnt / eventLogList] %v / %p", rowcnt, eventLogList)
}
