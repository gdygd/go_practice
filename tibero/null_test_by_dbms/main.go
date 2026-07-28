package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/godror/godror"

	_ "github.com/go-sql-driver/mysql"
)

func connOracle() (string, string) {
	return "godror", fmt.Sprintf("user=\"%s\" password=\"%s\" connectString=\"%s/%s\"", "UWSIG_TEST", "UWSIG_TEST", "192.168.2.161:1521", "theroadora11")
}

func connTibero() (string, string) {
	return "odbc", fmt.Sprintf(`Driver={%s};Server=%s;Port=%d;Database=%s;UID=%s;PWD=%s;`, "Tibero 7 ODBC driver", "10.1.0.120", 8629, "tibero", "MOON", "MOON")
}

func connMariaDBByToad() (string, string) {
	return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "moon", "moon", "10.1.1.164", 3306, "uwpis_test2")
}

func connMariaDBByHeidiSQL() (string, string) {
	// return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "dev", "dev", "10.1.0.115", 3306, "dbhandle_test") // blob_test
	return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "dev", "dev", "10.1.0.115", 3306, "sig_test") // signalmap
}

type DBNODE struct {
	NODE_ID string `json:"id"`
	NODE_NM string `json:"name"`
	//LOC_ID           int     `json:"locId"`
	NODE_TYPE        int     `json:"type"`
	NODE_LAT         float64 `json:"lat"`
	NODE_LON         float64 `json:"lon"`
	TURN_P           string  `json:"turnp"`
	REMARK           string  `json:"remark"`
	GEOM             string  `json:"-"`
	NODE_CREATE_TYPE int     `json:"createType"`
}

func main() {
	// connType, dsn := connOracle()
	connType, dsn := connTibero()
	// connType, dsn := connMariaDBByToad()
	// connType, dsn := connMariaDBByHeidiSQL()

	db, err := sql.Open(connType, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/////////////////////////////////
	// mariadb
	// mdbHandleNullQuery := `select a.NODE_ID, ifnull(NODE_NM, "' '"), ifnull(NODE_TYPE, 0) NODE_TYPE
	// 						, NODE_LAT, NODE_LON, ifnull(TURN_P, '0') TURN_P
	// 						, ifnull(REMARK, '0') REMARK
	// 						, ifnull(NODE_CREATE_TYPE, 0) NODE_CREATE_TYPE
	// 					from node a
	// 					left outer join (select *
	// 									from loc_mst
	// 									where use_yn = 1 and del_yn = 0) b
	// 					on a.NODE_ID = b.NODE_ID
	// 					WHERE remark IS NULL
	// 					AND turn_p IS NULL
	// 					order by node_id desc
	// 					limit 2`

	// mdbHandleNullOnlyIntQuery := `select a.NODE_ID, a.NODE_NM, ifnull(NODE_TYPE, 0) NODE_TYPE
	// 						, NODE_LAT, NODE_LON, ifnull(TURN_P, 0) TURN_P
	// 						, a.REMARK
	// 						, ifnull(NODE_CREATE_TYPE, 0) NODE_CREATE_TYPE
	// 					from node a
	// 					left outer join (select *
	// 									from loc_mst
	// 									where use_yn = 1 and del_yn = 0) b
	// 					on a.NODE_ID = b.NODE_ID
	// 					WHERE remark IS NULL
	// 					AND turn_p IS NULL
	// 					order by node_id desc
	// 					limit 2`

	// mdbHandleNullOnlyStringQuery := `select a.NODE_ID, ifnull(a.NODE_NM, '-') NODE_NM, a.NODE_TYPE
	// 									, NODE_LAT, NODE_LON, a.TURN_P
	// 									, ifnull(a.REMARK, '-') REMARK
	// 									, a.NODE_CREATE_TYPE
	// 								from node a
	// 								left outer join (select *
	// 												from loc_mst
	// 												where use_yn = 1 and del_yn = 0) b
	// 								on a.NODE_ID = b.NODE_ID
	// 								WHERE remark IS NULL
	// 								AND turn_p IS NULL
	// 								order by node_id desc
	// 								limit 2`

	// mdbIsNullQuery := `select a.NODE_ID, a.NODE_NM , a.NODE_TYPE, a.NODE_LAT, a.NODE_LON, a.TURN_P, a.REMARK, a.NODE_CREATE_TYPE
	// 					from node a
	// 					left outer join (select *
	// 									from loc_mst
	// 									where use_yn = 1 and del_yn = 0) b
	// 					on a.NODE_ID = b.NODE_ID
	// 					WHERE remark IS NULL
	// 					AND turn_p IS NULL
	// 					order by node_id desc
	// 					limit 2`

	//////////////////////////////////
	// oracle / tibero
	// oraHandleNullQuery := `
	// 			select b.node_id, b.node_nm, b.node_type, b.node_lat, b.node_lon, b.turn_p, b.remark, b.node_create_type
	// 			from(
	// 					select a.*, rownum as rnum
	// 					from (
	// 								select a.NODE_ID, NVL(NODE_NM, '-') NODE_NM , NVL(NODE_TYPE, 0) NODE_TYPE
	// 									, NODE_LAT, NODE_LON, NVL(TURN_P, 0) TURN_P
	// 									, NVL(REMARK, '-') REMARK
	// 									, NVL(NODE_CREATE_TYPE, 0) NODE_CREATE_TYPE
	// 								from node a
	// 								left outer join (select *
	// 												from loc_mst
	// 												where use_yn = 1 and del_yn = 0) b
	// 								on a.NODE_ID = b.NODE_ID
	// 								WHERE remark IS NULL
	// 								AND turn_p IS NULL
	// 								order by node_id desc
	// 						) a
	// 				) b
	// 			where rnum <= 2`

	// oraHandleNullOnlyIntQuery := `
	// 			select b.node_id, b.node_nm, b.node_type, b.node_lat, b.node_lon, b.turn_p, b.remark, b.node_create_type
	// 			from(
	// 					select a.*, rownum as rnum
	// 					from (
	// 								select a.NODE_ID, a.NODE_NM , NVL(NODE_TYPE, 0) NODE_TYPE
	// 									, NODE_LAT, NODE_LON, NVL(TURN_P, 0) TURN_P
	// 									, a.REMARK
	// 									, NVL(NODE_CREATE_TYPE, 0) NODE_CREATE_TYPE
	// 								from node a
	// 								left outer join (select *
	// 												from loc_mst
	// 												where use_yn = 1 and del_yn = 0) b
	// 								on a.NODE_ID = b.NODE_ID
	// 								WHERE remark IS NULL
	// 								AND turn_p IS NULL
	// 								order by node_id desc
	// 						) a
	// 				) b
	// 			where rnum <= 2`

	oraHandleNullOnlyStringQuery := `
				select b.node_id, b.node_nm, b.node_type, b.node_lat, b.node_lon, b.turn_p, b.remark, b.node_create_type
				from(
						select a.*, rownum as rnum
						from (
									select a.NODE_ID, nvl(a.NODE_NM, '-') NODE_NM , a.NODE_TYPE
										, NODE_LAT, NODE_LON, a.TURN_P
										, nvl(a.REMARK, '-') REMARK
										, a.NODE_CREATE_TYPE
									from node a
									left outer join (select *
													from loc_mst
													where use_yn = 1 and del_yn = 0) b
									on a.NODE_ID = b.NODE_ID
									WHERE remark IS NULL
									AND turn_p IS NULL
									order by node_id desc
							) a
					) b
				where rnum <= 2`

	// oraIsNullQuery := `
	// 			select b.node_id, b.node_nm, b.node_type, b.node_lat, b.node_lon, b.turn_p, b.remark, b.node_create_type
	// 			from(
	// 					select a.*, rownum as rnum
	// 					from (
	// 								select a.NODE_ID, a.NODE_NM , a.NODE_TYPE, a.NODE_LAT, a.NODE_LON, a.TURN_P, a.REMARK, a.NODE_CREATE_TYPE
	// 								from node a
	// 								left outer join (select *
	// 												from loc_mst
	// 												where use_yn = 1 and del_yn = 0) b
	// 								on a.NODE_ID = b.NODE_ID
	// 								WHERE remark IS NULL
	// 								AND turn_p IS NULL
	// 								order by node_id desc
	// 						) a
	// 				) b
	// 			where rnum <= 2`

	finalQuery := oraHandleNullOnlyStringQuery
	rows, err := db.Query(finalQuery)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	nodeList := make([]DBNODE, 0)

	for rows.Next() {
		nodeData := DBNODE{}
		err := rows.Scan(&nodeData.NODE_ID, &nodeData.NODE_NM, &nodeData.NODE_TYPE, &nodeData.NODE_LAT, &nodeData.NODE_LON, &nodeData.TURN_P, &nodeData.REMARK, &nodeData.NODE_CREATE_TYPE)

		if err != nil {
			log.Fatal(err)
		}
		nodeList = append(nodeList, nodeData)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for i, v := range nodeList {
		log.Printf(`nodeList[%d] :: %+v`, i, v)
	}
}
