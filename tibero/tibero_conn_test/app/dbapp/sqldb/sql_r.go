package sqldb

import (
	"context"
	"log"
	"strings"
	"tct/app/am"
	"tct/app/dbapp"
)

// ------------------------------------------------------------------------------
// ReadEventLogExceptOprt
// ------------------------------------------------------------------------------
func (m *SQLHandler) ReadEventLogExceptOprt(evtLogParam am.EventLogParam) ([]am.DBEventLog, int, error) {
	db, dbErr := m.Open()
	ctx := context.Background()
	// db, ctx, cancelcxt, donectx, dbErr := m.OpenCtx()
	// var count int = 0
	var rowcnt int = 0
	eventLogList := make([]am.DBEventLog, 0)
	var err error = nil

	// go cancelcxt()
	defer func() {
		// donectx()
		m.Close(db)
	}()
	if dbErr != nil {
		log.Printf("[ReadEventLogExceptOprt DB open error] : %s \n", dbErr.Error())
		return nil, rowcnt, dbErr
	}

	// 파라미터 설정
	// 1) 날짜
	repl := strings.NewReplacer("-", "", ":", "", " ", "")
	evtLogParam.SDt = repl.Replace(evtLogParam.SDt)
	evtLogParam.EDt = repl.Replace(evtLogParam.EDt)

	// 2) order by create_dt, seq
	evtLogParam.StrSort = "asc"
	if evtLogParam.SearchTp == 1 {
		evtLogParam.StrSort = "desc"
	}

	// 3) 다른 필드로 sort
	var strSortfields []string
	for _, sortfield := range evtLogParam.SortFields {
		var field []string

		switch sortfield.Field {
		case "uType": // 장치유형
			field = append(field, "a.UNIT_TYPE")
		case "code": // 중분류
			field = append(field, "a.EVENT_CODE")
		case "lowCode": // 세분류
			field = append(field, "a.EVENT_LOW_CODE")
		case "ctor": // 발생
			field = append(field, "a.CTOR")
		}

		switch sortfield.Order {
		case 0: // 오름차순(asc)
			field = append(field, "asc")
		case 1: // 내림차순(desc)
			field = append(field, "desc")
		}

		strSortfields = append(strSortfields, strings.Join(field, " "))
	}
	evtLogParam.Sortfield = strings.Join(strSortfields, ", ") // order by에서 다른 필드로 조회하는 구문만 완성
	if len(evtLogParam.Sortfield) > 0 {
		evtLogParam.Sortfield += ", "
	}

	// 4) event low code
	evtLogParam.EvLowCode = strings.Join(dbapp.ConvertIntArrToStrArr(evtLogParam.EvLowCodes), ", ")

	// DB 버전별 함수 생성
	if m.dbVer == "1.0" {
		eventLogList, rowcnt, err = m.readEventLogExceptOprtV1(db, ctx, evtLogParam)
		if err != nil {
			log.Printf("[ReadEventLogExceptOprt, readEventLogExceptOprtV1 error] : %s \n", err.Error())
			return nil, rowcnt, err
		}

		// remainrows 설정
		rowcnt = rowcnt - (evtLogParam.Offset + evtLogParam.ReqRows)
		if rowcnt < 0 {
			rowcnt = 0
		}

	} else if m.dbVer == "2.0" {
		eventLogList, rowcnt, err = m.readEventLogExceptOprtV2(db, ctx, evtLogParam)
		if err != nil {
			log.Printf("[ReadEventLogExceptOprt, readEventLogExceptOprtV2 error] : %s \n", err.Error())
			return nil, rowcnt, err
		}

		// SEARCH_MORE_ROW로 검색된 결과는 제외하고 응답
		var finalcnt int = rowcnt - am.SEARCH_MORE_ROW // SEARCH_MORE_ROW 개수만큼 제외하고 응답
		if rowcnt <= evtLogParam.ReqRows {             // 조회된 개수가 reqRow 이하면, 조회된 개수로 응답
			finalcnt = rowcnt
		}

		if finalcnt < 0 {
			finalcnt = 0
		}
		eventLogList = eventLogList[:finalcnt]

		// remainrows 설정
		// 다음 데이터 있는 경우 1, 없는 경우 0 설정
		if rowcnt > evtLogParam.ReqRows {
			rowcnt = 1
		} else {
			rowcnt = 0
		}
	}

	log.Println("[ReadEventLogExceptOprt ok]")
	return eventLogList, rowcnt, nil
}
