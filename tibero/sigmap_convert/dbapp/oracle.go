package dbapp

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// ODBC 데이터베이스 접속 정보
var oracle = `user="UWSIG_TEST" password="UWSIG_TEST" connectString="192.168.2.161:1521/theroadora11"`

// ##############################
// 운영 이력
// ##############################
// 1) 그룹
func OraInsGrpOprstt(collDtArr []string) error {
	var err error = nil

	return err
}

// 2) 교차로
func OraInsLocOprReqRst(collDtArr []string) error {
	var err error = nil

	return err
}

// 3) 교차로 상태
func OraInsLocColl(collDtArr []string) error {
	var err error = nil

	db, err := sql.Open("godror", oracle)
	if err != nil {
		log.Println("ERROR !! ", err)
		return err
	}
	defer db.Close()

	var strInsQrys []string
	for _, collDt := range collDtArr {
		for i := 1; i <= 50; i++ { // 교차로 1~50
			strInsQry := fmt.Sprintf(`SELECT %d AS LOC_ID, to_date('%s', 'YYYY-MM-DD HH24:MI:SS') AS COLL_DT, 1 AS SEQ, '014A4A000500000000004A5050010000000000000000000000' AS STATE FROM DUAL`, i, collDt)
			strInsQrys = append(strInsQrys, strInsQry) // id별 쿼리문 추가
		}
	}
	// for _, collDt := range collDtArr {
	// 	for i := -1; i >= -10; i-- { // seq = -1 ~ -10
	// 		// strInsQry := fmt.Sprintf(`SELECT 86 AS LOC_ID, to_date('%s', 'YYYY-MM-DD HH24:MI:SS') AS COLL_DT, %d AS SEQ, '0' AS STATE FROM DUAL`, collDt, i)
	// 		strInsQry := fmt.Sprintf(`SELECT 86 AS LOC_ID, to_date('%s', 'YYYY-MM-DD HH24:MI:SS') AS COLL_DT, %d AS SEQ, '0' AS STATE FROM DUAL`, collDt, i)
	// 		strInsQrys = append(strInsQrys, strInsQry) // id별 쿼리문 추가
	// 	}
	// }

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

// 3) 교차로 상태
func OraInsLocColl2(strInsQrys []string) error {
	var err error = nil

	db, err := sql.Open("godror", oracle)
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
