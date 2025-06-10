package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=127.0.0.1 port=5432 user=root password=secret dbname=customers sslmode=disable"

	fmt.Println("Connecting to customers DB...")
	db, err := sql.Open("postgres", connStr)
	checkErr2(err)
	defer db.Close()

	fmt.Println("Inserting customers...")

	// 10억 고객을 1천만 단위로 나누어 100번 반복 → 성능 테스트
	for i := 0; i < 100; i++ {
		fmt.Printf("Inserting 10M customers... (%d/100)\n", i+1)

		query := `
			INSERT INTO customers(name)
			SELECT md5(random()::text)
			FROM generate_series(1, 10000000)
		`

		_, err := db.Exec(query)
		checkErr2(err)
	}

	fmt.Println("Done.")
}

func checkErr2(err error) {
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
