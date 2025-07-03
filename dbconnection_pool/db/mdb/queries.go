package mdb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewHandler(user, pw, dbNm, host string, port int) *Queries {

	return &Queries{user: user, pw: pw, dbNm: dbNm, host: host, port: port}
}

type Queries struct {
	//db       db.DBTX
	db       *sql.DB
	user     string
	pw       string
	dbNm     string
	sid      string
	host     string
	port     int
	connType string
	dsn      string
	mu       sync.RWMutex
}

func (q *Queries) Init() error {

	// dsn := fmt.Sprintf("%s/%s@%s:%d/%s", q.user, q.pw, q.host, q.port, q.sid)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", q.user, q.pw, q.host, q.port, q.dbNm)
	fmt.Println("dns : ", dsn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Connection pool 설정
	// db.SetMaxOpenConns(30)                  // 동시 최대 연결 수
	// db.SetMaxIdleConns(15)                  // 유휴 상태로 유지할 연결 수
	// db.SetConnMaxLifetime(10 * time.Minute) // 연결의 최대 수명
	db.SetMaxOpenConns(5)                  // 동시 최대 연결 수
	db.SetMaxIdleConns(3)                  // 유휴 상태로 유지할 연결 수
	db.SetConnMaxLifetime(1 * time.Minute) // 연결의 최대 수명

	// Ping으로 연결 확인 (타임아웃 적용)
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	fmt.Println("db open")
	q.mu.Lock()
	q.db = db
	q.mu.Unlock()

	return nil
}
func (q *Queries) GetDB() *sql.DB {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.db
}

func (q *Queries) Close(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func (q *Queries) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err : %v, rv err %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
