package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"tct/app/am"
	"tct/app/dbapp"
	"time"
)

// ------------------------------------------------------------------------------
// readEventLogExceptOprtV1
// ------------------------------------------------------------------------------
func (m *SQLHandler) readEventLogExceptOprtV1(db *sql.DB, ctx context.Context, evtLogParam am.EventLogParam) ([]am.DBEventLog, int, error) {
	var count int = 0
	var qryUtype, qryEvCode, qryEvLowCode, qryUId, qryDate, qryCtor, qryOem, qryOrder string

	log.Printf("readEventLogExceptOprtV1 evtLogParam : %+v \n", evtLogParam)

	// 날짜
	qryDate += fmt.Sprintf(`and a.CREATE_DT BETWEEN %s AND %s`,
		m.toDate(true, evtLogParam.SDt, "", "", ""), m.toDate(true, evtLogParam.EDt, "", "", ""))

	// OEM
	// 하이패스 이벤트 제외하고 조회 (uw)
	// not uw => (3 and 4 and 1 and 21)필요없음
	// if am.AppVar.Oem == "uw" { // readUnitEventLog, readUnitAllEventLog, readAllEventLog에서만 사용 -> 이벤트 이력(전체/시스템/그룹/교차로)만 사용
	qryOem = `and NOT(a.UNIT_TYPE = 3 and a.EVENT_CODE = 4 and a.EVENT_LOW_CODE = 1 and a.PARAM_1 = 21)`
	// }

	// 전체, 시스템/그룹/교차로 분기 -> uType값 확인
	if evtLogParam.UTypeOk {
		// [시스템/그룹/교차로]
		qryUtype += fmt.Sprintf(`and a.UNIT_TYPE = %d`, evtLogParam.UType) // uType
		qryUId += fmt.Sprintf(`and a.UNIT_ID = %d`, evtLogParam.UId)       // uId

		if evtLogParam.EvCodeOk { // 1) 이벤트 유형 선택
			qryEvCode += fmt.Sprintf(`and a.EVENT_CODE = %d`, evtLogParam.EvCode)              // evCode
			qryEvLowCode += fmt.Sprintf("and a.EVENT_LOW_CODE in (%s)", evtLogParam.EvLowCode) // evLowCode

		} else { // 2) 이벤트 유형 전체
			// evCode
			qryEvCode += fmt.Sprintf(`and (a.EVENT_CODE, a.EVENT_LOW_CODE) in (select s.EVENT_CODE, s.EVENT_LOW_CODE
																				from UNIT_EVENT_LOW_CODE s
																				where s.UNIT_TYPE = %d
																				and s.USE_YN = 1 )`, evtLogParam.UType)
		}
	}

	// order by
	qryOrder = fmt.Sprintf(`order by %s a.CREATE_DT %s, a.SEQ %s`, evtLogParam.Sortfield, evtLogParam.StrSort, evtLogParam.StrSort)

	// total count
	strQry1 := fmt.Sprintf(`
	select count(*)
	from unit_event_state_log a
	where 1=1
	%s -- qryUtype
	%s -- qryEvCode
	%s -- qryEvLowCode
	%s -- qryUId
	%s -- qryDate
	%s -- qryCtor
	%s -- qryOem`,
		qryUtype,
		qryEvCode,
		qryEvLowCode,
		qryUId,
		qryDate,
		qryCtor,
		qryOem)

	// 장치 사용자 선택 이벤트
	qryobj := am.QRY_DBEventLog{
		CREATE_DATE: m.toChar(false, "a.CREATE_DT", true, true, "-", " ", ":"),
		PARAM_1:     m.handleNull("a.PARAM_1", 0),
		PARAM_2:     m.handleNull("a.PARAM_2", 0),
		PARAM_3:     m.handleNull("a.PARAM_3", 0),
		PARAM_4:     m.handleNull("a.PARAM_4", 0),
		PARAM_5:     m.handleNull("a.PARAM_5", 0),
	}

	// rownum에서 사용
	colnmArr := []string{
		"CREATE_DT", "UNIT_TYPE", "EVENT_CODE", "EVENT_LOW_CODE", "UNIT_ID", "CTOR", "TARGET",
		"PARAM_1", "PARAM_2", "PARAM_3", "PARAM_4", "PARAM_5"}

	strQry := fmt.Sprintf(`
			-- // rownum start
			%s
				select {{.CREATE_DATE}} CREATE_DT
					, a.UNIT_TYPE
					, a.EVENT_CODE
					, a.EVENT_LOW_CODE
					, a.UNIT_ID
					, a.CTOR
					, a.TARGET
					, {{.PARAM_1}} PARAM_1
					, {{.PARAM_2}} PARAM_2
					, {{.PARAM_3}} PARAM_3
					, {{.PARAM_4}} PARAM_4
					, {{.PARAM_5}} PARAM_5
				from unit_event_state_log a
				where 1=1
				%s -- qryUtype
				%s -- qryEvCode
				%s -- qryEvLowCode
				%s -- qryUId
				%s -- qryDate
				%s -- qryCtor
				%s -- qryOem
				%s -- qryOrder
			%s
			-- // rownum end
			%s -- limit, offset`,
		m.handleRownumStart(strings.Join(colnmArr, ", ")),
		qryUtype,
		qryEvCode,
		qryEvLowCode,
		qryUId,
		qryDate,
		qryCtor,
		qryOem,
		qryOrder,
		m.handleRownumEnd(),
		m.handleGetNRows(evtLogParam.ReqRows, evtLogParam.Offset))

	finalQuery, err := dbapp.MakeQuery("readEventLogExceptOprtV1", strQry, qryobj)
	if err != nil {
		log.Printf("[readEventLogExceptOprtV1 makeQuery error] : %s \n", err.Error())
		return nil, count, err
	}

	log.Printf("readEventLogExceptOprtV1 CNT qry : %v \n", strQry1)
	log.Printf("readEventLogExceptOprtV1 qry : %v \n", finalQuery)

	rowscnt, err := db.QueryContext(ctx, strQry1)
	if err != nil {
		return nil, count, err
	}
	defer rowscnt.Close()

	if rowscnt.Next() {
		// float64로 Scan()할 임시 변수 (count(*)값이 7자리 이상이면 float64로 Scan()됨)
		var tempCnt float64
		err := rowscnt.Scan(&tempCnt)
		if err != nil {
			log.Printf("[readEventLogExceptOprtV1 Query error(1)] : %s \n", err.Error())
			return nil, count, err
		}

		count = int(tempCnt)
	}

	rows, err := db.QueryContext(ctx, finalQuery)

	if err != nil {
		log.Printf("[readEventLogExceptOprtV1 QueryContext error] : %s \n", err.Error())
		return nil, count, err
	}
	defer rows.Close()
	eventLogList := make([]am.DBEventLog, 0)

	var idx int = 0
	var lastIdx int = 0
	lastIdx = evtLogParam.ReqRows
	if evtLogParam.ReqRows >= (count - evtLogParam.Offset) {
		lastIdx = count - evtLogParam.Offset
	}

	for rows.Next() {
		eventLog := am.DBEventLog{}
		var rnum int = 0
		var param1, param2, param3, param4, param5 float64

		// Oracle, Tibero는 rownum값을 Scan()해야함
		switch m.dbms {
		case ORACLE, TIBERO:
			err = rows.Scan(&eventLog.CREATE_DATE, &eventLog.UNIT_TYPE, &eventLog.EVENT_CODE, &eventLog.EVENT_LOW_CODE, &eventLog.UNIT_ID, &eventLog.CTOR, &eventLog.TARGET,
				&param1, &param2, &param3, &param4, &param5, &rnum)
		case MARIADB:
			err = rows.Scan(&eventLog.CREATE_DATE, &eventLog.UNIT_TYPE, &eventLog.EVENT_CODE, &eventLog.EVENT_LOW_CODE, &eventLog.UNIT_ID, &eventLog.CTOR, &eventLog.TARGET,
				&param1, &param2, &param3, &param4, &param5)
		}

		eventLog.PARAM_1 = int(param1)
		eventLog.PARAM_2 = int(param2)
		eventLog.PARAM_3 = int(param3)
		eventLog.PARAM_4 = int(param4)
		eventLog.PARAM_5 = int(param5)

		if err != nil {
			log.Printf("[readEventLogExceptOprtV1 Query error(1)] : %s \n", err.Error())
			return nil, count, err
		}
		eventLogList = append(eventLogList, eventLog)

		if idx+1 == lastIdx {
			break
		}

		idx += 1

		if err := rows.Err(); err != nil {
			log.Printf("[readEventLogExceptOprtV1 Query error(2)] : %s \n", err.Error())
			return nil, count, err
		}

	}

	if err = rows.Err(); err != nil {
		log.Printf("[readEventLogExceptOprtV1 Query error(3)] : %s \n", err.Error())
		return nil, count, err
	}

	log.Println("[readEventLogExceptOprtV1 ok]")
	return eventLogList, count, nil
}

// ------------------------------------------------------------------------------
// readEventLogExceptOprtV2
// ------------------------------------------------------------------------------
func (m *SQLHandler) readEventLogExceptOprtV2(db *sql.DB, ctx context.Context, evtLogParam am.EventLogParam) ([]am.DBEventLog, int, error) {
	// var count int = 0
	var rowcnt int = 0
	var qryUtype, qryEvCode, qryEvLowCode, qryUId, qryDate, qryCtor, qryOem, qryOrder string

	log.Printf("readEventLogExceptOprtV2 evtLogParam : %+v \n", evtLogParam)

	// 날짜
	qryDate += fmt.Sprintf(`and a.CREATE_DT BETWEEN %s AND %s`,
		m.toDate(true, evtLogParam.SDt, "", "", ""), m.toDate(true, evtLogParam.EDt, "", "", ""))

	// OEM
	// 하이패스 이벤트 제외하고 조회 (uw)
	// not uw => (3 and 4 and 1 and 21)필요없음
	// if am.AppVar.Oem == "uw" { // readUnitEventLog, readUnitAllEventLog, readAllEventLog에서만 사용 -> 이벤트 이력(전체/시스템/그룹/교차로)만 사용
	// qryOem = `and NOT(a.UNIT_TYPE = 3 and a.EVENT_CODE = 4 and a.EVENT_LOW_CODE = 1 and a.PARAM_1 = 21)`
	// }

	// 전체, 시스템/그룹/교차로 분기 -> uType값 확인
	if evtLogParam.UTypeOk {
		// [시스템/그룹/교차로]
		qryUtype += fmt.Sprintf(`and a.UNIT_TYPE = %d`, evtLogParam.UType) // uType
		qryUId += fmt.Sprintf(`and a.UNIT_ID = %d`, evtLogParam.UId)       // uId

		if evtLogParam.EvCodeOk { // 1) 이벤트 유형 선택
			qryEvCode += fmt.Sprintf(`and a.EVENT_CODE = %d`, evtLogParam.EvCode)              // evCode
			qryEvLowCode += fmt.Sprintf("and a.EVENT_LOW_CODE in (%s)", evtLogParam.EvLowCode) // evLowCode

		} else { // 2) 이벤트 유형 전체
			// evCode
			qryEvCode += fmt.Sprintf(`and (a.EVENT_CODE, a.EVENT_LOW_CODE) in (select s.EVENT_CODE, s.EVENT_LOW_CODE
																				from TCS_M_UNIT_EVENT_LOW_CODE s
																				where s.UNIT_TYPE = %d
																				and s.USE_YN = 1 )`, evtLogParam.UType)
		}
	}

	// order by
	qryOrder = fmt.Sprintf(`order by %s a.CREATE_DT %s, a.SEQ %s`, evtLogParam.Sortfield, evtLogParam.StrSort, evtLogParam.StrSort)

	// total count
	// strQry1 := fmt.Sprintf(`
	// select count(*)
	// from TCS_L_UNIT_EVENT a
	// where 1=1
	// %s -- qryUtype
	// %s -- qryEvCode
	// %s -- qryEvLowCode
	// %s -- qryUId
	// %s -- qryDate
	// %s -- qryCtor
	// %s -- qryOem`,
	// 	qryUtype,
	// 	qryEvCode,
	// 	qryEvLowCode,
	// 	qryUId,
	// 	qryDate,
	// 	qryCtor,
	// 	qryOem)

	// 장치 사용자 선택 이벤트
	qryobj := am.QRY_DBEventLog{
		CREATE_DATE: m.toChar(false, "a.CREATE_DT", true, true, "-", " ", ":"),
		PARAM_1:     m.handleNull("a.PARAM_1", 0),
		PARAM_2:     m.handleNull("a.PARAM_2", 0),
		PARAM_3:     m.handleNull("a.PARAM_3", 0),
		PARAM_4:     m.handleNull("a.PARAM_4", 0),
		PARAM_5:     m.handleNull("a.PARAM_5", 0),
	}

	// rownum에서 사용
	colnmArr := []string{
		"CREATE_DT", "UNIT_TYPE", "EVENT_CODE", "EVENT_LOW_CODE", "UNIT_ID", "CTOR", "TARGET",
		"PARAM_1", "PARAM_2", "PARAM_3", "PARAM_4", "PARAM_5"}

	strQry := fmt.Sprintf(`
			-- // rownum start
			%s
				select {{.CREATE_DATE}} CREATE_DT
					, a.UNIT_TYPE
					, a.EVENT_CODE
					, a.EVENT_LOW_CODE
					, a.UNIT_ID
					, a.CTOR
					, a.TARGET
					, {{.PARAM_1}} PARAM_1
					, {{.PARAM_2}} PARAM_2
					, {{.PARAM_3}} PARAM_3
					, {{.PARAM_4}} PARAM_4
					, {{.PARAM_5}} PARAM_5
				from TCS_L_UNIT_EVENT a
				where 1=1
				%s -- qryUtype
				%s -- qryEvCode
				%s -- qryEvLowCode
				%s -- qryUId
				%s -- qryDate
				%s -- qryCtor
				%s -- qryOem
				%s -- qryOrder
			%s
			-- // rownum end
			%s -- limit, offset`,
		m.handleRownumStart(strings.Join(colnmArr, ", ")),
		qryUtype,
		qryEvCode,
		qryEvLowCode,
		qryUId,
		qryDate,
		qryCtor,
		qryOem,
		qryOrder,
		m.handleRownumEnd(),
		m.handleGetNRows(evtLogParam.ReqRows+am.SEARCH_MORE_ROW, evtLogParam.Offset))

	finalQuery, err := dbapp.MakeQuery("readEventLogExceptOprtV2", strQry, qryobj)
	if err != nil {
		log.Printf("[readEventLogExceptOprtV2 makeQuery error] : %s \n", err.Error())
		return nil, rowcnt, err
	}

	// log.Printf(5, "readEventLogExceptOprtV2 CNT qry : %v", strQry1)
	log.Printf("readEventLogExceptOprtV2 qry : %v \n", finalQuery)

	// rowscnt, err := db.QueryContext(ctx, strQry1)
	// if err != nil {
	// 	return nil, count, err
	// }
	// defer rowscnt.Close()

	// if rowscnt.Next() {
	// 	// float64로 Scan()할 임시 변수 (count(*)값이 7자리 이상이면 float64로 Scan()됨)
	// 	var tempCnt float64
	// 	err := rowscnt.Scan(&tempCnt)
	// 	if err != nil {
	// 		log.Printf("[readEventLogExceptOprtV2 Query error(1)] : %s", err.Error())
	// 		return nil, count, err
	// 	}

	// 	count = int(tempCnt)
	// }

	// go func() {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			log.Println("취소 신호 감지:", ctx.Err())
	// 			return
	// 		default:
	// 			log.Println("작업 진행 중...")
	// 			time.Sleep(500 * time.Millisecond)
	// 		}
	// 	}
	// }()

	// time1
	log.Printf("time1 !\n")

	// rows, err = db.QueryContext(ctx, finalQuery)

	var startTm time.Time = time.Now()
	var isok bool = false
	var rows *sql.Rows = new(sql.Rows)

	log.Println("readEventLogExceptOprtV2, QueryContext start !!!")
	go func() {
		rows, err = db.QueryContext(ctx, finalQuery)
		isok = true
	}()

	for {
		if isok {
			break
		} else {
			// 타임아웃 발생
			if dbapp.CheckElapsedTime(&startTm, QRY_TIMEOUT_MS) {
				err = context.Canceled // context canceled
				return nil, rowcnt, err
			}
		}

		time.Sleep(time.Millisecond * 50)
	}

	// time2
	log.Printf("time2 !!\n")

	if err != nil {
		log.Printf("[readEventLogExceptOprtV2 QueryContext error] : %s \n", err.Error())
		return nil, rowcnt, err
	}
	defer rows.Close()
	eventLogList := make([]am.DBEventLog, 0)

	// var idx int = 0
	// var lastIdx int = 0
	// lastIdx = evtLogParam.ReqRows
	// if evtLogParam.ReqRows >= (count - evtLogParam.Offset) {
	// 	lastIdx = count - evtLogParam.Offset
	// }

	// time3
	log.Printf("time3 !!!\n")

	for rows.Next() {
		eventLog := am.DBEventLog{}
		var rnum int = 0
		var param1, param2, param3, param4, param5 float64

		// Oracle, Tibero는 rownum값을 Scan()해야함
		switch m.dbms {
		case ORACLE, TIBERO:
			err = rows.Scan(&eventLog.CREATE_DATE, &eventLog.UNIT_TYPE, &eventLog.EVENT_CODE, &eventLog.EVENT_LOW_CODE, &eventLog.UNIT_ID, &eventLog.CTOR, &eventLog.TARGET,
				&param1, &param2, &param3, &param4, &param5, &rnum)
		case MARIADB:
			err = rows.Scan(&eventLog.CREATE_DATE, &eventLog.UNIT_TYPE, &eventLog.EVENT_CODE, &eventLog.EVENT_LOW_CODE, &eventLog.UNIT_ID, &eventLog.CTOR, &eventLog.TARGET,
				&param1, &param2, &param3, &param4, &param5)
		}

		eventLog.PARAM_1 = int(param1)
		eventLog.PARAM_2 = int(param2)
		eventLog.PARAM_3 = int(param3)
		eventLog.PARAM_4 = int(param4)
		eventLog.PARAM_5 = int(param5)

		if err != nil {
			log.Printf("[readEventLogExceptOprtV2 Query error(1)] : %s \n", err.Error())
			return nil, rowcnt, err
		}
		eventLogList = append(eventLogList, eventLog)

		// if idx+1 == lastIdx {
		// 	break
		// }

		// idx += 1

		rowcnt++ // 조회되는 row 개수 확인

		if err = rows.Err(); err != nil {
			log.Printf("[readEventLogExceptOprtV2 Query error(2)] : %s \n", err.Error())
			return nil, rowcnt, err
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("[readEventLogExceptOprtV2 Query error(3)] : %s \n", err.Error())
		return nil, rowcnt, err
	}

	log.Println("[readEventLogExceptOprtV2 ok]")
	return eventLogList, rowcnt, nil
}
