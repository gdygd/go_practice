package dbapp

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// ODBC 데이터베이스 접속 정보
var tiberoDataSourceName = "driver={Tibero 7 ODBC driver};server=10.1.0.120;port=8629;uid=SIGUSR_TEST;pwd=SIGUSR_TEST;database=tibero"

// ##############################
// 운영 이력
// ##############################
// 1) 그룹
func TdbInsGrpOprstt(collDtArr []string) error {
	var err error = nil

	return err
}

// 2) 교차로
func TdbInsLocOprReqRst(collDtArr []string) error {
	var err error = nil

	// ODBC 데이터베이스에 연결
	db, err := sql.Open("odbc", tiberoDataSourceName)
	if err != nil {
		log.Println("ERROR !! ", err)
		return err
	}
	defer db.Close()

	var strInsQrys []string
	for _, collDt := range collDtArr {
		for i := 1; i <= 50; i++ { // 교차로 1~50
			strInsQry := fmt.Sprintf(`
				SELECT %d AS LOC_ID
				     , to_date('%s', 'YYYY-MM-DD HH24:MI:SS') AS CREATE_DT
					 , 1 AS SEQ
					 , sysdate AS COLL_DT
					 , 0 AS REQ_CTRLMODE
					 , 160 AS REQ_CYCLE
					 , 0 AS REQ_OFFSET_SEC
					 , '0' AS REQ_PHASE
					 , 0 AS RST_CTRLMODE
					 , 0 AS RST_CTRLSTATE
					 , '0' AS RST_PHASE
					 , 160 AS RST_CYCLE
					 , '0' AS RST_PED_PHASE
					 , 0 AS PPC_FLASHTIME
					 , 0 AS PPC_ALLREDTIME
					 , 0 AS PPC_PHASE
					 , 2 AS PLANNUM
					 , 0 AS BANK_ID
				FROM DUAL`, i, collDt)
			strInsQrys = append(strInsQrys, strInsQry) // id별 쿼리문 추가
		}
	}
	// for _, collDt := range collDtArr {
	// 	for i := 2; i <= 10; i++ { // seq = 2 ~ 10
	// 		strInsQry := fmt.Sprintf(`SELECT 1 AS LOC_ID, to_date('%s', 'YYYY-MM-DD HH24:MI:SS') AS COLL_DT, %d AS SEQ, '0' AS STATE FROM DUAL`, collDt, i)
	// 		strInsQrys = append(strInsQrys, strInsQry) // id별 쿼리문 추가
	// 	}
	// }

	// 쿼리문 완성
	// index1 = TCS_L_LOC_OPRREQ_RST_PART_IDX1 (loc_id, req_dt)
	// index2 = TCS_L_LOC_OPRREQ_RST_PART_IDX2 (req_dt)
	strQry := fmt.Sprintf(`
	INSERT INTO TCS_L_LOC_OPRREQ_RST_PART_IDX1
	%s`, strings.Join(strInsQrys, "\n UNION ALL \n"))

	// INSERT INTO TCS_L_LOC_COLL_COPY2
	// SELECT * FROM DUAL
	// UNION ALL
	// SELECT * FROM DUAL;

	// 쿼리문 확인
	// fmt.Println(strQry)

	// 쿼리 실행
	_, err = db.Exec(strQry)
	if err != nil {
		log.Println("ERROR !!!! ", err)
		return err
	}

	return err
}

// 3) 교차로 상태
func TdbInsLocColl(collDtArr []string) error {
	var err error = nil

	// ODBC 데이터베이스에 연결
	db, err := sql.Open("odbc", tiberoDataSourceName)
	if err != nil {
		log.Println("ERROR !! ", err)
		return err
	}
	defer db.Close()

	var strInsQrys []string
	for _, collDt := range collDtArr {
		for i := 1; i <= 50; i++ { // 교차로 1~50
			strInsQry := fmt.Sprintf(`SELECT %d AS LOC_ID
										   , to_date('%s', 'YYYY-MM-DD HH24:MI:SS') AS COLL_DT
										   , 1 AS SEQ
										   , '0' AS STATE FROM DUAL`, i, collDt)
			strInsQrys = append(strInsQrys, strInsQry) // id별 쿼리문 추가
		}
	}
	// for _, collDt := range collDtArr {
	// 	for i := 2; i <= 10; i++ { // seq = 2 ~ 10
	// 		strInsQry := fmt.Sprintf(`SELECT 1 AS LOC_ID
	// 									   , to_date('%s', 'YYYY-MM-DD HH24:MI:SS') AS COLL_DT
	// 									   , %d AS SEQ
	// 									   , '0' AS STATE FROM DUAL`, collDt, i)
	// 		strInsQrys = append(strInsQrys, strInsQry) // id별 쿼리문 추가
	// 	}
	// }

	// 쿼리문 완성
	// index1 = TCS_L_LOC_COLL_PART_IDX1 (loc_id, coll_dt)
	// index2 = TCS_L_LOC_COLL_PART_IDX2 (coll_dt)
	strQry := fmt.Sprintf(`
		-- INSERT INTO TCS_L_LOC_COLL_PART_IDX1
		INSERT INTO TCS_L_LOC_COLL_DETAIL_PART_IDX
		%s`, strings.Join(strInsQrys, "\n UNION ALL \n"))

	// INSERT INTO TCS_L_LOC_COLL_COPY2
	// SELECT * FROM DUAL
	// UNION ALL
	// SELECT * FROM DUAL;

	// 쿼리문 확인
	// fmt.Println(strQry)

	// 쿼리 실행
	_, err = db.Exec(strQry)
	if err != nil {
		log.Println("ERROR !!!! ", err)
		return err
	}

	return err
}

// 3) 교차로 상태2
func TdbInsLocColl2(strInsQrys []string) error {
	var err error = nil

	// ODBC 데이터베이스에 연결
	db, err := sql.Open("odbc", tiberoDataSourceName)
	if err != nil {
		log.Println("ERROR !! ", err)
		return err
	}
	defer db.Close()

	// 쿼리문 완성
	// index1 = TCS_L_LOC_COLL_PART_IDX1 (loc_id, coll_dt)
	// index2 = TCS_L_LOC_COLL_PART_IDX2 (coll_dt)
	strQry := fmt.Sprintf(`
		INSERT INTO LOC_COLL_STATE_LOG_Y1
		%s`, strings.Join(strInsQrys, "\n UNION ALL \n"))

	// INSERT INTO TCS_L_LOC_COLL_COPY2
	// SELECT * FROM DUAL
	// UNION ALL
	// SELECT * FROM DUAL;

	// 쿼리문 확인
	// fmt.Println(strQry)

	// 쿼리 실행
	_, err = db.Exec(strQry)
	if err != nil {
		log.Println("ERROR !!!! ", err)
		return err
	}

	return err
}

// 3) 교차로 상태2
func TdbInsLocColl3(strInsQrys []string) error {
	var err error = nil

	// ODBC 데이터베이스에 연결
	db, err := sql.Open("odbc", tiberoDataSourceName)
	if err != nil {
		log.Println("ERROR !! ", err)
		return err
	}
	defer db.Close()

	// 쿼리문 완성
	// index1 = TCS_L_LOC_COLL_PART_IDX1 (loc_id, coll_dt)
	// index2 = TCS_L_LOC_COLL_PART_IDX2 (coll_dt)
	strQry := fmt.Sprintf(`
		INSERT INTO LOC_COLL_STATE_LOG_Y2
		%s`, strings.Join(strInsQrys, "\n UNION ALL \n"))

	// INSERT INTO TCS_L_LOC_COLL_COPY2
	// SELECT * FROM DUAL
	// UNION ALL
	// SELECT * FROM DUAL;

	// 쿼리문 확인
	// fmt.Println(strQry)

	// 쿼리 실행
	_, err = db.Exec(strQry)
	if err != nil {
		log.Println("ERROR !!!! ", err)
		return err
	}

	return err
}

// 3) 교차로 상태3
func TdbInsLocColl4(strInsQrys []string) error {
	var err error = nil

	// ODBC 데이터베이스에 연결
	db, err := sql.Open("odbc", tiberoDataSourceName)
	if err != nil {
		log.Println("ERROR !! ", err)
		return err
	}
	defer db.Close()

	// 쿼리문 완성
	// index1 = TCS_L_LOC_COLL_PART_IDX1 (loc_id, coll_dt)
	// index2 = TCS_L_LOC_COLL_PART_IDX2 (coll_dt)
	strQry := fmt.Sprintf(`
		INSERT INTO LOC_COLL_STATE_LOG_Y3
		%s`, strings.Join(strInsQrys, "\n UNION ALL \n"))

	// INSERT INTO TCS_L_LOC_COLL_COPY2
	// SELECT * FROM DUAL
	// UNION ALL
	// SELECT * FROM DUAL;

	// 쿼리문 확인
	// fmt.Println(strQry)

	// 쿼리 실행
	_, err = db.Exec(strQry)
	if err != nil {
		log.Println("ERROR !!!! ", err)
		return err
	}

	return err
}
