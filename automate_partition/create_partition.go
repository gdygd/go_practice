package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// 연결 정보
	postgresConn := "host=127.0.0.1 port=5432 user=root password=secret dbname=postgres sslmode=disable"

	// 1. postgres 데이터베이스에 연결
	fmt.Println("Connecting to postgres...")
	dbPostgres, err := sql.Open("postgres", postgresConn)
	checkErr(err)
	defer dbPostgres.Close()

	// 2. customers DB 생성 (기존 삭제 주석 처리)
	fmt.Println("Dropping and creating database customers...")
	_, _ = dbPostgres.Exec("DROP DATABASE customers") // 주석: 필요 시 수동 실행
	// _, err = dbPostgres.Exec("CREATE DATABASE customers")
	checkErr(err)

	// // 3. customers DB에 연결
	// customersConn := "host=127.0.0.1 port=5432 user=root password=secret dbname=customers sslmode=disable"
	// fmt.Println("Connecting to customers DB...")
	// dbCustomers, err := sql.Open("postgres", customersConn)
	// checkErr(err)
	// defer dbCustomers.Close()

	// // 4. customers 테이블 생성 (파티션 테이블)
	// fmt.Println("Creating customers partitioned table...")
	// createTableSQL := `CREATE TABLE customers (
	// 	id SERIAL,
	// 	name TEXT
	// ) PARTITION BY RANGE (id)`
	// _, err = dbCustomers.Exec(createTableSQL)
	// checkErr(err)

	// // 5. 파티션 생성 루프
	// fmt.Println("Creating partitions...")
	// for i := 0; i < 100; i++ {
	// 	idFrom := i * 10000000
	// 	idTo := (i + 1) * 10000000
	// 	partitionName := fmt.Sprintf("customers_%d_%d", idFrom, idTo)

	// 	createPartitionSQL := fmt.Sprintf(`CREATE TABLE %s (LIKE customers INCLUDING INDEXES)`, partitionName)
	// 	attachPartitionSQL := fmt.Sprintf(`
	// 		ALTER TABLE customers
	// 		ATTACH PARTITION %s
	// 		FOR VALUES FROM (%d) TO (%d)
	// 	`, partitionName, idFrom, idTo)

	// 	fmt.Printf("Creating partition: %s\n", partitionName)
	// 	_, err = dbCustomers.Exec(createPartitionSQL)
	// 	checkErr(err)
	// 	_, err = dbCustomers.Exec(attachPartitionSQL)
	// 	checkErr(err)
	// }

	fmt.Println("All partitions created successfully.")
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
