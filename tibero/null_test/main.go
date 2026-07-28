package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/godror/godror"

	_ "github.com/go-sql-driver/mysql"
)

var connType, dsn string

const (
	ORACLE  = "oracle"
	TIBERO  = "tibero"
	MARIADB = "mariadb"
)

func connOracle() (string, string) {
	return "godror", fmt.Sprintf("user=\"%s\" password=\"%s\" connectString=\"%s/%s\"", "UWSIG_TEST", "UWSIG_TEST", "192.168.2.161:1521", "theroadora11")
}

func connTibero() (string, string) {
	return "odbc", fmt.Sprintf(`Driver={%s};Server=%s;Port=%d;Database=%s;UID=%s;PWD=%s;`, "Tibero 7 ODBC driver", "10.1.0.120", 8629, "tibero", "TB_TEST", "TB_TEST")
}

func connMariaDBByToad() (string, string) {
	return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "moon", "moon", "10.1.1.164", 3306, "uwpis_test2")
}

func connMariaDBByHeidiSQL() (string, string) {
	// return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "dev", "dev", "10.1.0.115", 3306, "dbhandle_test") // blob_test
	return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "dev", "dev", "10.1.0.115", 3306, "sig_test") // signalmap
}

func marshalTest(db *sql.DB) {
	qryobj := QRY_DBGroup{
		GRP_NUM:         handleNull("GRP_NUM", 0),
		GRP_DEFCTRLMODE: handleNull("GRP_DEFCTRLMODE", 0),
		AUTO_ONLINEYN:   handleNull("AUTO_ONLINEYN", 0),
		GRP_LAT:         handleNull("GRP_LAT", 0.0),
		GRP_LON:         handleNull("GRP_LON", 0.0),
		GRP_ZOMMLV:      handleNull("GRP_ZOMMLV", 0.0),
		MIN_CYCLE:       handleNull("MIN_CYCLE", 0),
		MAX_CYCLE:       handleNull("MAX_CYCLE", 0),
		TRANS_MODE:      handleNull("TRANS_MODE", 0),
		DEL_YN:          handleNull("DEL_YN", 0),
	}

	query := `SELECT GRP_ID
						, GRP_NM
						, {{.GRP_NUM}} GRP_NUM
						, {{.GRP_DEFCTRLMODE}} GRP_DEFCTRLMODE
						, {{.AUTO_ONLINEYN}} AUTO_ONLINEYN
						, {{.GRP_LAT}} GRP_LAT
						, {{.GRP_LON}} GRP_LON
						, {{.GRP_ZOMMLV}} GRP_ZOMMLV
						, {{.MIN_CYCLE}} MIN_CYCLE
						, {{.MAX_CYCLE}} MAX_CYCLE
						, {{.TRANS_MODE}} TRANS_MODE
						, USE_YN
						, {{.DEL_YN}} DEL_YN
					from TCS_M_GROUP a
					where (a.DEL_YN = 0 or a.DEL_YN is NULL)
					order by a.GRP_NUM asc`

	finalQuery, err := MakeQuery("readGroupV1", query, qryobj)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(finalQuery)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	grps := make([]DBGroup, 0)

	for rows.Next() {
		var grp DBGroup = DBGroup{}
		err := rows.Scan(&grp.GRP_ID, &grp.GRP_NM, &grp.GRP_NUM, &grp.GRP_DEFCTRLMODE, &grp.AUTO_ONLINEYN, &grp.GRP_LAT, &grp.GRP_LON, &grp.GRP_ZOMMLV, &grp.MIN_CYCLE, &grp.MAX_CYCLE, &grp.TRANS_MODE, &grp.USE_YN, &grp.DEL_YN)

		if err != nil {
			log.Fatal(err)
		}

		grps = append(grps, grp)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	result, err := json.Marshal(grps)
	if err != nil {
		log.Fatal(err)
	}

	// mpGrp := make(map[int]DBGroup)
	// for _, data := range grps {
	// 	mpGrp[data.GRP_ID] = data
	// }

	// result, err := json.Marshal(mpGrp)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("success!!")
	log.Printf(`[string(result)] :::: %v`, string(result))
	// log.Printf(`[string(result)] :::: %v`, grps)
}

// func marshalAndUnmarshalTest() {
// 	// Creating instances with valid values
// 	nb := NullBool{sql.NullBool{Bool: true, Valid: true}}
// 	nf := NullFloat64{sql.NullFloat64{Float64: 3.14, Valid: true}}
// 	ni := NullInt64{sql.NullInt64{Int64: 42, Valid: true}}
// 	ns := NullString{sql.NullString{String: "Hello, World!", Valid: true}}

// 	// Marshaling to JSON
// 	nbJson, err := json.Marshal(nb)
// 	if err != nil {
// 		log.Fatalf("Error marshaling NullBool: %v", err)
// 	}
// 	nfJson, err := json.Marshal(nf)
// 	if err != nil {
// 		log.Fatalf("Error marshaling NullFloat64: %v", err)
// 	}
// 	niJson, err := json.Marshal(ni)
// 	if err != nil {
// 		log.Fatalf("Error marshaling NullInt64: %v", err)
// 	}
// 	nsJson, err := json.Marshal(ns)
// 	if err != nil {
// 		log.Fatalf("Error marshaling NullString: %v", err)
// 	}

// 	log.Printf("NullBool JSON: %s\n", nbJson)
// 	log.Printf("NullFloat64 JSON: %s\n", nfJson)
// 	log.Printf("NullInt64 JSON: %s\n", niJson)
// 	log.Printf("NullString JSON: %s\n", nsJson)
// 	fmt.Println()
// 	fmt.Println()

// 	// Creating instances with null values
// 	nb = NullBool{sql.NullBool{Valid: false}}
// 	nf = NullFloat64{sql.NullFloat64{Valid: false}}
// 	ni = NullInt64{sql.NullInt64{Valid: false}}
// 	ns = NullString{sql.NullString{Valid: false}}

// 	// Marshaling to JSON
// 	nbJson, err = json.Marshal(nb)
// 	if err != nil {
// 		log.Fatalf("Error marshaling NullBool: %v", err)
// 	}
// 	nfJson, err = json.Marshal(nf)
// 	if err != nil {
// 		log.Fatalf("Error marshaling NullFloat64: %v", err)
// 	}
// 	niJson, err = json.Marshal(ni)
// 	if err != nil {
// 		log.Fatalf("Error marshaling NullInt64: %v", err)
// 	}
// 	nsJson, err = json.Marshal(ns)
// 	if err != nil {
// 		log.Fatalf("Error marshaling NullString: %v", err)
// 	}

// 	log.Printf("NullBool JSON with null: %s\n", nbJson)
// 	log.Printf("NullFloat64 JSON with null: %s\n", nfJson)
// 	log.Printf("NullInt64 JSON with null: %s\n", niJson)
// 	log.Printf("NullString JSON with null: %s\n", nsJson)
// 	fmt.Println()
// 	fmt.Println()

// 	// Unmarshaling from JSON
// 	nbJson = []byte("true")
// 	nfJson = []byte("3.14")
// 	niJson = []byte("42")
// 	nsJson = []byte("\"Hello, World!\"")

// 	err = json.Unmarshal(nbJson, &nb)
// 	if err != nil {
// 		log.Fatalf("Error unmarshaling NullBool: %v", err)
// 	}
// 	err = json.Unmarshal(nfJson, &nf)
// 	if err != nil {
// 		log.Fatalf("Error unmarshaling NullFloat64: %v", err)
// 	}
// 	err = json.Unmarshal(niJson, &ni)
// 	if err != nil {
// 		log.Fatalf("Error unmarshaling NullInt64: %v", err)
// 	}
// 	err = json.Unmarshal(nsJson, &ns)
// 	if err != nil {
// 		log.Fatalf("Error unmarshaling NullString: %v", err)
// 	}

// 	log.Printf("Unmarshaled NullBool: %v, Valid: %v\n", nb.Bool, nb.Valid)
// 	log.Printf("Unmarshaled NullFloat64: %v, Valid: %v\n", nf.Float64, nf.Valid)
// 	log.Printf("Unmarshaled NullInt64: %v, Valid: %v\n", ni.Int64, ni.Valid)
// 	log.Printf("Unmarshaled NullString: %v, Valid: %v\n", ns.String, ns.Valid)
// 	fmt.Println()
// 	fmt.Println()

// 	// Unmarshaling null values from JSON
// 	nbJson = []byte("")
// 	nfJson = []byte("")
// 	niJson = []byte("")
// 	nsJson = []byte("")

// 	err = json.Unmarshal(nbJson, &nb)
// 	if err != nil {
// 		log.Fatalf("Error unmarshaling NullBool: %v", err)
// 	}
// 	err = json.Unmarshal(nfJson, &nf)
// 	if err != nil {
// 		log.Fatalf("Error unmarshaling NullFloat64: %v", err)
// 	}
// 	err = json.Unmarshal(niJson, &ni)
// 	if err != nil {
// 		log.Fatalf("Error unmarshaling NullInt64: %v", err)
// 	}
// 	err = json.Unmarshal(nsJson, &ns)
// 	if err != nil {
// 		log.Fatalf("Error unmarshaling NullString: %v", err)
// 	}

// 	log.Printf("Unmarshaled NullBool with null: Valid: %+v\n", nb)
// 	log.Printf("Unmarshaled NullFloat64 with null: Valid: %+v\n", nf)
// 	log.Printf("Unmarshaled NullInt64 with null: Valid: %+v\n", ni)
// 	log.Printf("Unmarshaled NullString with null: Valid: %+v\n", ns)
// 	fmt.Println()
// 	fmt.Println()
// }

// //////////////////////////
// ////     main     ////////
// //////////////////////////
func main() {
	// connType, dsn := connOracle()
	connType, dsn = connTibero()
	// connType, dsn := connMariaDBByToad()
	// connType, dsn := connMariaDBByHeidiSQL()

	db, err := sql.Open(connType, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// marshalTest(db)
	// marshalAndUnmarshalTest()

	http.HandleFunc("/nulltest", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/nulltest")

		body, _ := io.ReadAll(r.Body)
		str := string((body))
		log.Printf("AddGroup, /group-add (1): %v", str)

		var u DBGroup = DBGroup{}
		json.Unmarshal([]byte(str), &u)

		json.NewEncoder(w).Encode(u)
	})

	http.HandleFunc("/longinttest", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/longinttest")

		// rows, err := db.Query(`select 1234567 from dual`)
		rows, err := db.Query(`select to_number(123456) from dual`)
		// rows, err := db.Query(`select cast(1234567 as number) from dual`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var num int
		if rows.Next() {
			err = rows.Scan(&num)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Value: %v, Type: %T\n", num, num)
		}

		fmt.Fprintf(w, `Value: %v, Type: %T`, num, num)
	})

	// 웹 서버 실행
	log.Println("Starting server on :8888")
	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Nullable Bool that overrides sql.NullBool
type NullBool struct {
	sql.NullBool
}

func (nb NullBool) MarshalJSON() ([]byte, error) {
	if nb.Valid {
		return json.Marshal(nb.Bool)
	}
	return json.Marshal("")
}

func (nb *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		nb.Valid = true
		nb.Bool = *b
	} else {
		nb.Valid = false
	}
	return nil
}

// Nullable Float64 that overrides sql.NullFloat64
type NullFloat64 struct {
	sql.NullFloat64
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Float64)
	}
	return json.Marshal(0)
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	} else {
		nf.Valid = false
	}
	return nil
}

// Nullable Int64 that overrides sql.NullInt64
type NullInt64 struct {
	sql.NullInt64
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int64)
	}
	return json.Marshal(0)
}

func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	} else {
		ni.Valid = false
	}
	return nil
}

// Nullable String that overrides sql.NullString
type NullString struct {
	str string
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	log.Println("마샬 호출")
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal("")
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

// ---------------------------------------------------------------------------
// DBGroup
// ---------------------------------------------------------------------------
type DBGroup struct {
	GRP_ID int `json:"id"`
	// GRP_NM string `json:"name,string"`
	GRP_NM          NullString  `json:"name"`
	GRP_NUM         int         `json:"gNum"`
	GRP_DEFCTRLMODE int         `json:"defMode"`
	AUTO_ONLINEYN   int         `json:"autoOnlineYN"`
	GRP_LAT         NullFloat64 `json:"lat"`
	GRP_LON         NullFloat64 `json:"lon"`
	GRP_ZOMMLV      NullFloat64 `json:"zoomlvl"`
	// GRP_LAT         float64 `json:"lat,float64"`
	// GRP_LON         float64 `json:"lon,float64"`
	// GRP_ZOMMLV      float64 `json:"zoomlvl,float64"`
	MIN_CYCLE  int `json:"minCyc"`
	MAX_CYCLE  int `json:"maxCyc"`
	TRANS_MODE int `json:"tranMode"`
	USE_YN     int `json:"useYn"`
	DEL_YN     int `json:"delYn"`
}

type QRY_DBGroup struct {
	GRP_ID          string `json:"id"`
	GRP_NM          string `json:"name"`
	GRP_NUM         string `json:"gNum"`
	GRP_DEFCTRLMODE string `json:"defMode"`
	AUTO_ONLINEYN   string `json:"autoOnlineYN"`
	GRP_LAT         string `json:"lat"`
	GRP_LON         string `json:"lon"`
	GRP_ZOMMLV      string `json:"zoomlvl"`
	MIN_CYCLE       string `json:"minCyc"`
	MAX_CYCLE       string `json:"maxCyc"`
	TRANS_MODE      string `json:"tranMode"`
	USE_YN          string `json:"useYn"`
	DEL_YN          string `json:"delYn"`
}

func handleNull(colnm string, value any) string {
	var nullQry string

	switch connType {
	case ORACLE, TIBERO:
		nullQry = fmt.Sprintf(`nvl(%s, %v)`, colnm, value)
	case MARIADB:
		nullQry = fmt.Sprintf(`ifnull(%s, %v)`, colnm, value)
	}

	return nullQry
}

func MakeQuery(tmplname, query string, data any) (string, error) {
	t, err := template.New(tmplname).Parse(query)
	if err != nil {
		return "", err
	}

	var queryBuilder strings.Builder
	err = t.Execute(&queryBuilder, data)
	if err != nil {
		return "", err
	}

	return queryBuilder.String(), nil
}
