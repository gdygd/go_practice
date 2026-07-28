package main

import (
	"dbapp"
	"fmt"
	"log"
	"sync"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
)

func main() {
	// 운영 이력
	// 1) 그룹
	// 2) 교차로
	// 3) 교차로 상태

	// - 데이터 500건 쌓아서 한 번에 insert
	// - state = 0
	// - seq = 1

	// index
	// 1) LOC_ID, COLL_DT
	// 2) COLL_DT

	var wg sync.WaitGroup

	// 총 2개의 고루틴 추가
	wg.Add(5)

	// 쿼리문 실행
	// ※ 7월을 기준으로 insert하고 있어서 1일 ~ 31일까지 insert
	// goroutine -1
	go func() {
		defer wg.Done()
		insert2(1, 5) // 1일 ~ 10일

	}()

	// goroutine -2
	go func() {
		defer wg.Done()
		insert2(6, 10) // 11일 ~ 20일
	}()

	go func() {
		defer wg.Done()
		insert2(11, 15) // 11일 ~ 20일
	}()

	go func() {
		defer wg.Done()
		insert2(16, 20) // 11일 ~ 20일
	}()

	go func() {
		defer wg.Done()
		insert2(21, 25) // 11일 ~ 20일
	}()

	// main thread
	// insert(21, 30) // 21일 ~ 30일
	insert2(26, 30) // 21일 ~ 31일

	// insert(5, 6) // 5일 ~ 6일
	// insert(11, 12) // 11일 ~ 12일
	// insert(17, 18) // 17일 ~ 18일
	// insert(23, 24) // 23일 ~ 24일
	// insert(29, 30) // 29일 ~ 30일

	wg.Wait()
	fmt.Println("모든 작업이 완료되었습니다.")
}

func insert(startDay, endDay int) {
	var err error = nil
	var collDt string

	// for day := 1; day <= 31; day++ { // day
	for day := startDay; day <= endDay; day++ { // day
		for hh := 0; hh < 24; hh++ { // hour
			var collDtArr []string

			for mm := 0; mm < 60; mm++ { // minute
				for ss := 0; ss < 60; ss++ { // second
					// collDt = fmt.Sprintf("2024-07-%02d %02d:%02d:%02d", day, hh, mm, ss) // 202407
					collDt = fmt.Sprintf("2024-10-%02d %02d:%02d:%02d", day, hh, mm, ss) // 202410
					// fmt.Println(collDt)
					collDtArr = append(collDtArr, collDt)
				}

				// if mm%5 == 0 { // 5분 데이터량 단위로 입력 (300개) // OUT OF MEMORY 에러 발생
				// if mm%2 == 0 && len(collDtArr) > 0 { // 메인 쓰레드로만 진행하면 에러 X / 고루틴 사용시 OUT OF MEMORY 에러 발생
				if mm%1 == 0 && len(collDtArr) > 0 {
					// log.Println(collDt)
					log.Println("Inserting data:", collDtArr[0], "to", collDtArr[len(collDtArr)-1])

					// #######################################
					// Oracle
					// #######################################
					// err = dbapp.OraInsGrpOprstt(collDtArr)
					// err = dbapp.OraInsLocOprReqRst(collDtArr)
					// err = dbapp.OraInsLocColl(collDtArr)

					// #######################################
					// Tibero
					// #######################################
					// err = dbapp.TdbInsGrpOprstt(collDtArr)
					// err = dbapp.TdbInsLocOprReqRst(collDtArr)
					err = dbapp.TdbInsLocColl(collDtArr)

					// #######################################
					// MariaDB
					// #######################################
					// err = dbapp.MdbInsGrpOprstt(collDtArr)
					// err = dbapp.MdbInsLocOprReqRst(collDtArr)
					// err = dbapp.MdbInsLocColl(collDtArr)
					if err != nil {
						log.Println("ERROR !!!!!!!!!!! ", err.Error())
						panic(err)
					}

					collDtArr = nil // 5분 데이터 저장 후 배열 초기화
				}
			}
		}
	}
}

func insert2(startDay, endDay int) {
	var err error = nil
	var collDt string

	for mon := 9; mon <= 12; mon++ {
		for day := startDay; day <= endDay; day++ { // day
			//for hh := 0; hh < 12; hh++ { // hour
			for hh := 0; hh < 24; hh++ { // hour
				var collDtArr []string

				var strInsQrys []string = []string{}
				for mm := 0; mm < 40; mm++ { // minute
					//for mm := 31; mm < 60; mm++ { // minute

					for ss := 0; ss < 30; ss++ { // second

						collDt = fmt.Sprintf("2024-%02d-%02d %02d:%02d:%02d", mon, day, hh, mm, ss) // 202410

						collDtArr = append(collDtArr, collDt)

						for id := 1; id <= 50; id++ {
							strInsQry := fmt.Sprintf(`SELECT %d AS LOC_ID, to_date('%s', 'YYYY-MM-DD HH24:MI:SS') AS COLL_DT, 1 AS SEQ, '014A4A000500000000004A5050010000000000000000000000' AS STATE FROM DUAL`, id, collDt)
							strInsQrys = append(strInsQrys, strInsQry) // id별 쿼리문 추가
						}
					}

					if mm%2 == 0 && len(collDtArr) > 0 {
						// log.Println(collDt)
						log.Println("Inserting data:", collDtArr[0], "to", collDtArr[len(collDtArr)-1])

						// #######################################
						// Oracle
						// #######################################
						// err = dbapp.OraInsGrpOprstt(collDtArr)
						// err = dbapp.OraInsLocOprReqRst(collDtArr)
						//err = dbapp.OraInsLocColl(collDtArr)
						// err = dbapp.OraInsLocColl2(strInsQrys)

						// #######################################
						// Tibero
						// #######################################
						// err = dbapp.TdbInsGrpOprstt(collDtArr)
						// err = dbapp.TdbInsLocOprReqRst(collDtArr)
						// err = dbapp.TdbInsLocColl(collDtArr)
						err = dbapp.TdbInsLocColl2(strInsQrys)
						// err1 := dbapp.TdbInsLocColl3(strInsQrys)
						// err2 := dbapp.TdbInsLocColl4(strInsQrys)

						// #######################################
						// MariaDB
						// #######################################
						// err = dbapp.MdbInsGrpOprstt(collDtArr)
						// err = dbapp.MdbInsLocOprReqRst(collDtArr)
						// err = dbapp.MdbInsLocColl(collDtArr)
						if err != nil {
							//log.Println("ERROR !!!!!!!!!!! ", err.Error(), err1.Error(), err2.Error())
							log.Println("ERROR !!!!!!!!!!! ", err.Error())
							panic(err)
						}
						strInsQrys = nil
						strInsQrys = []string{}
						collDtArr = nil // 5분 데이터 저장 후 배열 초기화
					}
				}
			}
		}

	}

}
