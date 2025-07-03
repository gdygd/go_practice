package main

import (
	"context"
	"dbconnection_pool/db/mdb"
	"fmt"
	"os"
	"time"
)

func main() {

	dbHnd := mdb.NewHandler("dev", "dev", "test_db", "10.1.0.114", 3306)
	if err := dbHnd.Init(); err != nil {
		fmt.Printf("err.. %v \n", err)
		os.Exit(0)
	}

	ctx := context.Background()
	for {
		tm, err := dbHnd.ReadSysdate(ctx)
		if err != nil {
			fmt.Printf("readsysdate.. err: %v \n", err)
			break
		}
		fmt.Printf("tm : %s \n", tm)

		time.Sleep(time.Second * 1)

	}

}
