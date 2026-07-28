package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"tct/app/dbapp"
	"time"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
)

const (
	QRY_TIMEOUT    = 5    // 5sec
	QRY_TIMEOUT_MS = 5000 // 5sec
)

const (
	ORACLE  = "oracle"
	TIBERO  = "tibero"
	MARIADB = "mariadb"
)

// var dblog *general.OLog2 = general.InitLogEnv("./log", "dblog", 2) // level 1~9

type SQLHandler struct {
	user string
	pw   string
	dbNm string
	host string
	port int

	dbVer string
	dbms  string

	driverNm string
	connType string
	dsn      string

	Connected bool
}

// type CancelContext func()
// type DoneContext func()

func Cancelctx(cancel context.CancelFunc) {
	if cancel == nil {
		return
	}
	// QRY_TIMEOUT 후 cancel
	log.Println("Cancelctx (1)")
	// time.Sleep(QRY_TIMEOUT * time.Second)
	// time.Sleep(500 * time.Millisecond)
	// time.Sleep(1 * time.Nanosecond)
	time.Sleep(2 * time.Second)
	cancel()
	log.Println("Cancelctx (2)")
}

func Donectx(cancel context.CancelFunc) {
	if cancel == nil {
		return
	}
	log.Println("Donectx (1)")
	//<-ctx.Done()
	cancel()
	log.Println("Donectx (2)")
}

// var am.Applog *general.OLog2 = general.InitLogEnv("./am.Applog", "am.Applog", 2) // level 1~9

// ------------------------------------------------------------------------------
// Open
// ------------------------------------------------------------------------------
func (m *SQLHandler) Open() (*sql.DB, error) {
	var strConn string
	var dbTimeout time.Duration = time.Second * 3 // Oracle, MariaDB, Tibero 3sec

	// DBMS별 접속 정보 생성
	switch m.dbms {
	case ORACLE:
		// strConn = fmt.Sprintf("user=\"%s\" password=\"%s\" connectString=\"%s/%s\"", m.user, m.pw, m.host, m.dbNm)
		strConn = fmt.Sprintf("user=\"%s\" password=\"%s\" connectString=\"%s:%d/%s\"", m.user, m.pw, m.host, m.port, m.dbNm)
		//strConn = `user="UWITS_TEST" password="UWITS_TEST" connectString="192.168.2.161:1521/theroadora11"`

	case TIBERO:
		strConn = fmt.Sprintf(`Driver={%s};Server=%s;Port=%d;Database=%s;UID=%s;PWD=%s;`, m.driverNm, m.host, m.port, m.dbNm, m.user, m.pw)
		// strConn = fmt.Sprintf(`DSN=%s`, m.dsn)

	case MARIADB:
		strConn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.user, m.pw, m.host, m.port, m.dbNm)
	}

	// log.Println(4, "database open..(1) : %s", strConn)

	var db *sql.DB
	var err error

	db, err = sql.Open(m.connType, strConn)

	if err != nil {
		log.Printf("database open err : %v\n", err.Error())
		m.Connected = false
		return nil, err
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// defer cancel()

	// if err = db.PingContext(ctx); err != nil {
	// 	log.Println(("database open ping context err : %v", err)
	// 	m.Connected = false
	// 	return nil, err
	// }

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		log.Println(1, "try DB Connect...")
		err := db.PingContext(ctx)
		done <- err
	}()

	select {
	case <-ctx.Done():
		log.Printf("Connection timeout: %v\n", ctx.Err().Error())
		m.Connected = false
		return db, ctx.Err()

	case err := <-done:
		if err != nil {
			log.Printf("Connection Error : %v\n", err.Error())
			m.Connected = false
			return db, err
		}
	}

	m.Connected = true

	return db, nil
}

// ------------------------------------------------------------------------------
// Open2
// ------------------------------------------------------------------------------
func (m *SQLHandler) Open2(dbNm string) (*sql.DB, error) {

	return nil, nil

}

// ------------------------------------------------------------------------------
// OpenCtx
// ------------------------------------------------------------------------------
func (m *SQLHandler) OpenCtx() (*sql.DB, context.Context, dbapp.CancelContext, dbapp.DoneContext, error) {
	log.Printf("database open..(1) : %v, %v \n", m.connType, m.dsn)

	var db *sql.DB
	var err error
	var strConn string

	// DBMS별 접속 정보 생성
	switch m.dbms {
	case ORACLE:
		strConn = fmt.Sprintf("user=\"%s\" password=\"%s\" connectString=\"%s/%s\"", m.user, m.pw, m.host, m.dbNm)
		//strConn = `user="UWITS_TEST" password="UWITS_TEST" connectString="192.168.2.161:1521/theroadora11"`

	case TIBERO:
		strConn = fmt.Sprintf(`Driver={%s};Server=%s;Port=%d;Database=%s;UID=%s;PWD=%s;`, m.driverNm, m.host, m.port, m.dbNm, m.user, m.pw)
		// strConn = fmt.Sprintf(`DSN=%s`, m.dsn)

	case MARIADB:
		strConn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.user, m.pw, m.host, m.port, m.dbNm)
	}

	log.Printf("database open..(1) : %v, %v\n", m.connType, strConn)
	db, err = sql.Open(m.connType, strConn)

	ctx := context.Background()
	if err != nil {
		log.Printf("database open err : %v\n", err.Error())
		m.Connected = false
		return nil, ctx, func() { Cancelctx(nil) }, func() { Donectx(nil) }, err
	}

	log.Printf("database open..(2) : %v, %v\n", m.connType, m.dsn)
	m.Connected = true

	// ctx2, cancel := context.WithCancel(context.Background())
	ctx2, cancel := context.WithCancel(ctx)

	return db, ctx2, func() { Cancelctx(cancel) }, func() { Donectx(cancel) }, nil
}

// ------------------------------------------------------------------------------
// OpenCtx2
// ------------------------------------------------------------------------------
func (m *SQLHandler) OpenCtx2() (*sql.DB, context.Context, context.CancelFunc, error) {

	return nil, nil, nil, nil
}

// ------------------------------------------------------------------------------
// OpenPing
// ------------------------------------------------------------------------------
func (m *SQLHandler) OpenPing(host string) (*sql.DB, error) {

	var strDsn = "DSN=" + m.dsn
	db, err := sql.Open(m.connType, strDsn)

	if err != nil {
		log.Printf("database open err : %v\n", err.Error())
		m.Connected = false
		return nil, err
	}

	// // var ctx context.Context
	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// defer cancel()

	// if err = db.PingContext(ctx); err != nil {
	// 	log.Println(("(OpenPing)database open ping context err : %v", err)
	// 	m.Connected = false
	// 	return nil, err
	// }

	m.Connected = true

	return db, nil
}

// ------------------------------------------------------------------------------
// Close
// ------------------------------------------------------------------------------
func (m *SQLHandler) Close(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

// ------------------------------------------------------------------------------
// Ping
// ------------------------------------------------------------------------------
func (m *SQLHandler) Ping() (bool, error) {

	var err error
	var rows *sql.Rows = new(sql.Rows)
	db, dbErr := m.Open()

	defer func() {
		log.Println(4, "[Ping close..]")
		m.Close(db)
		rows.Close()
	}()

	if dbErr != nil {
		log.Printf("[Ping db open error] : %s\n" + dbErr.Error())
		return false, dbErr
	}

	log.Println(4, "[Ping qeury...]")
	var val int = 0
	rows, err = db.Query(`select 1 VAL from dual`)
	if err != nil {
		log.Println(4, "[Ping qeury...] err %v", err.Error())
		return false, err
	}

	if rows.Next() {
		err := rows.Scan(&val)

		if err != nil {
			log.Printf("[Ping Query error(1)] : %s\n" + err.Error())
			return false, err
		}
	}

	log.Printf("[Ping qeury...](%d)\n", val)

	return true, nil
}

// ------------------------------------------------------------------------------
// GetConnected
// ------------------------------------------------------------------------------
func (m *SQLHandler) GetConnected() bool {
	return m.Connected
}

// ------------------------------------------------------------------------------
// ChangeHostAddress
// ------------------------------------------------------------------------------
func (m *SQLHandler) ChangeHostAddress(host string) {
	m.host = host
}

// ------------------------------------------------------------------------------
// NewSQLHandler
// ------------------------------------------------------------------------------
func NewSQLHandler(user, pw, dbNm, host string, port int, dbver, dbms, driver, dsn string) *SQLHandler {
	// connection 타입 지정
	var connType string
	switch dbms {
	case ORACLE:
		connType = "godror"
	case TIBERO:
		connType = "odbc"
	case MARIADB:
		connType = "mysql"
	}

	return &SQLHandler{
		user: user,
		pw:   pw,
		dbNm: dbNm,
		host: host,
		port: port,

		dbVer: dbver,
		dbms:  dbms,

		driverNm: driver,
		connType: connType,
		dsn:      dsn,

		Connected: false,
	}
}
