package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/godror/godror"
)

func main() {
	// dsn := `user="UWSIGR28_TEST" password="UWSIGR28_TEST" connectString="(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=192.168.2.161)(PORT=11521))(CONNECT_DATA=(SERVICE_NAME = theroadora11)))"`
	// dsn := `user="UWSIGR28_TEST" password="UWSIGR28_TEST" connectString="(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=192.168.2.161)(PORT=1521))(CONNECT_DATA=(SERVICE_NAME=theroadora11)))"`
	dsn := `user="UWSIGR28_TEST" password="UWSIGR28_TEST" connectString="THEROADORA11"`

	db, err := sql.Open("godror", dsn)
	if err != nil {
		log.Fatalf("failed to open connection: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	var sysdate string
	err = db.QueryRowContext(ctx, "SELECT TO_CHAR(SYSDATE, 'YYYY-MM-DD HH24:MI:SS') FROM dual").Scan(&sysdate)
	if err != nil {
		log.Fatalf("query failed: %v", err)
	}

	fmt.Println("Oracle SYSDATE:", sysdate)
}
