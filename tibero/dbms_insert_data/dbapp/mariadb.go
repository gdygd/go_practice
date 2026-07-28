package dbapp

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// ODBC 데이터베이스 접속 정보
var mariadb = `dev:dev@tcp(10.1.0.115:3306)/sigusr_test`

// ##############################
// 운영 이력
// ##############################
// 1) 그룹
func MdbInsGrpOprstt(collDtArr []string) error {
	var err error = nil

	return err
}

// 2) 교차로
func MdbInsLocOprReqRst(collDtArr []string) error {
	var err error = nil

	return err
}

// 3) 교차로 상태
func MdbInsLocColl(collDtArr []string) error {
	var err error = nil

	db, err := sql.Open("mysql", mariadb)
	if err != nil {
		log.Println("ERROR !! ", err)
		return err
	}
	defer db.Close()

	var strInsQrys []string
	var time_format string = string("%Y-%m-%d %H:%i:%s")
	// 교차로 여러개 - seq 고정 insert
	for _, collDt := range collDtArr {
		// for i := 1; i <= 10; i++ { // 교차로 1~10
		for i := 11; i <= 50; i++ { // 교차로 11~50
			strInsQry := fmt.Sprintf(`SELECT %d AS LOC_ID
			                               , str_to_date('%s', '%s') AS COLL_DT
										   , 1 AS SEQ
										   , '0' AS STATE FROM DUAL`, i, collDt, time_format)
			strInsQrys = append(strInsQrys, strInsQry) // id별 쿼리문 추가
		}
	}

	// 교차로 고정 - seq 여러개 insert
	// for _, collDt := range collDtArr {
	// 	for i := 2; i <= 10; i++ { // seq 2~10
	// 		strInsQry := fmt.Sprintf(`SELECT 1 AS LOC_ID
	// 		                               , str_to_date('%s', '%s') AS COLL_DT
	// 									   , %d AS SEQ
	// 									   , '0' AS STATE FROM DUAL`, collDt, time_format, i)
	// 		strInsQrys = append(strInsQrys, strInsQry) // id별 쿼리문 추가
	// 	}
	// }

	// 쿼리문 완성
	// index1 = TCS_L_LOC_COLL_PART_IDX1 (loc_id, coll_dt)
	// index2 = TCS_L_LOC_COLL_PART_IDX2 (coll_dt)
	strQry := fmt.Sprintf(`
		INSERT INTO TCS_L_LOC_COLL_PART_IDX1
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
