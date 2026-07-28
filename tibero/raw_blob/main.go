package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/godror/godror"

	_ "github.com/go-sql-driver/mysql"
)

func connOracle() (string, string) {
	return "godror", fmt.Sprintf("user=\"%s\" password=\"%s\" connectString=\"%s/%s\"", "ORATEST", "ORATEST", "192.168.2.161:1521", "theroadora11")
}

func connTibero() (string, string) {
	// return "odbc", fmt.Sprintf(`Driver={%s};Server=%s;Port=%d;Database=%s;UID=%s;PWD=%s;`, "Tibero 7 ODBC driver", "10.1.0.120", 8629, "tibero", "TB_TEST", "TB_TEST")
	return "odbc", fmt.Sprintf(`Driver={%s};Server=%s;Port=%d;Database=%s;UID=%s;PWD=%s;`, "Tibero 7 ODBC driver", "10.1.0.120", 8629, "tibero", "MOON", "MOON")
}

func connMariaDBByToad() (string, string) {
	return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "moon", "moon", "10.1.1.164", 3306, "uwpis_test2")
}

func connMariaDBByHeidiSQL() (string, string) {
	// return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "dev", "dev", "10.1.0.115", 3306, "dbhandle_test") // blob_test
	return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "dev", "dev", "10.1.0.115", 3306, "sig_test") // signalmap
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

	// connection test
	rows, err := db.Query("select '[START]' from dual")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var str string
	for rows.Next() {
		err := rows.Scan(&str)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf(`[str] : %s`, str)

	type blobTest struct {
		Id   int       `json:"id"`
		Varb [608]byte `json:"varb"`
		Blob [608]byte `json:"blob"`
	}

	// 간단한 웹 서버 핸들러
	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		str := string(body)

		result := blobTest{}
		json.Unmarshal([]byte(str), &result)

		// log.Printf("[result] :: %+v", result)

		var stepArr []byte = []byte{}
		stepArr = append(stepArr, result.Blob[:]...)

		// query := fmt.Sprintf(`
		// insert into blob_test(varb_val, blob_val)
		// values ('%X', '%X')`, stepArr, stepArr)

		var query string
		switch connType {
		case "godror":
			query = `insert into blob_test(varb_val, blob_val)
					 values(:1, :2)`
		case "odbc":

		case "mysql":
			query = `insert into blob_test(varb_val, blob_val)
					 values(?, ?)`
		}

		_, err = db.Exec(query, stepArr, stepArr)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, "result: %+v", result)
	})

	// BLOB값 확인
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/query")

		query := `select varb_val, blob_val from blob_test`
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var varbList [][]uint8 = make([][]uint8, 0, 0)
		var blobList [][]uint8 = make([][]uint8, 0, 0)

		for rows.Next() {
			var varb_val []uint8 = make([]uint8, 1216, 1216)
			var blob_val []uint8 = make([]uint8, 1216, 1216)

			err := rows.Scan(&varb_val, &blob_val)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("varb_val: %v", varb_val)
			log.Println()
			log.Printf("blob_val: %v", blob_val)

			varbList = append(varbList, varb_val)
			blobList = append(blobList, blob_val)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		log.Printf("varbList: %v", varbList)
		log.Println()
		log.Printf("blobList: %v", blobList)

		fmt.Fprintf(w, "Query result: %v", blobList)
	})

	// BLOB값 확인
	http.HandleFunc("/signalmap", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/signalmap")

		query := `select data from loc_signal_map`
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var dataList [][]uint8 = make([][]uint8, 0, 0)

		for rows.Next() {
			var data []uint8 = make([]uint8, 1216, 1216)

			err := rows.Scan(&data)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("data: %v", data)
			log.Println()

			dataList = append(dataList, data)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		log.Printf("dataList: %v", dataList)
		log.Println()

		fmt.Fprintf(w, "Query result: %v", dataList)
	})

	http.HandleFunc("/nulltest", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/nulltest")

		// query := `select name from null_test`
		query := `select nvl(id, 0) as id, nvl(name, '') as name from null_test`
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		type nullTest struct {
			id   int
			name string
			// name sql.NullString
		}

		var nulltests []nullTest

		for rows.Next() {
			var nulltest nullTest

			err = rows.Scan(&nulltest.id, &nulltest.name)
			// err = rows.Scan(&nulltest.name)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf(`[%d] :: [%s]`, nulltest.id, nulltest.name)
			// log.Printf(`[%v]`, nulltest.name)
			// log.Printf(`[Value] :: %s .... [isNull] :: %v`, nulltest.name.String, nulltest.name.Valid)

			nulltests = append(nulltests, nulltest)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, "Query result: %v", nulltests)
	})

	// 웹 서버 실행
	log.Println("Starting server on :8888")
	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
