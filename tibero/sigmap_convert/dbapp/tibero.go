package dbapp

import (
	"database/sql"
	"log"
)

// ODBC 데이터베이스 접속 정보
var tiberoDataSourceName = "driver={Tibero 7 ODBC driver};server=10.1.0.120;port=8629;uid=SH_SIGTEST;pwd=SH_SIGTEST;database=tibero"

func ReadSignalMap() ([]Signalmap, error) {
	var err error = nil
	db, err := sql.Open("odbc", tiberoDataSourceName)
	if err != nil {
		log.Println("ERROR !! ", err)
		return nil, err
	}
	defer db.Close()

	strQry := `
		select IT_NO      
			, SI_MAPNO   
			, SI_RINGNO  
			, SI_STEP    
			, SI_SSR1_VHE
			, SI_SSR2_VHE
			, SI_SSR3_VHE
			, SI_SSR4_VHE
			, SI_SSR5_VHE
			, SI_SSR6_VHE
			, SI_SSR7_VHE
			, SI_SSR8_VHE
			, SI_SSR1_PED
			, SI_SSR2_PED
			, SI_SSR3_PED
			, SI_SSR4_PED
			, SI_SSR5_PED
			, SI_SSR6_PED
			, SI_SSR7_PED
			, SI_SSR8_PED
			, SI_MINTIME 
			, SI_MAXTIME 
			, SI_EOP     
		from SIGNAL_MAP_T_TEMP
		order by IT_NO, SI_MAPNO, SI_RINGNO, SI_STEP;
	`

	log.Printf("qry : %v\n", strQry)
	rows, err := db.Query(strQry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sigmaps := []Signalmap{}

	for rows.Next() {
		var r Signalmap
		err := rows.Scan(&r.IT_NO, &r.SI_MAPNO, &r.SI_RINGNO, &r.SI_STEP, &r.SI_SSR1_VHE, &r.SI_SSR2_VHE, &r.SI_SSR3_VHE, &r.SI_SSR4_VHE, &r.SI_SSR5_VHE, &r.SI_SSR6_VHE, &r.SI_SSR7_VHE, &r.SI_SSR8_VHE, &r.SI_SSR1_PED, &r.SI_SSR2_PED, &r.SI_SSR3_PED, &r.SI_SSR4_PED, &r.SI_SSR5_PED, &r.SI_SSR6_PED, &r.SI_SSR7_PED, &r.SI_SSR8_PED, &r.SI_MINTIME, &r.SI_MAXTIME, &r.SI_EOP)
		if err != nil {
			return nil, err
		}
		sigmaps = append(sigmaps, r)
	}

	return sigmaps, nil
}

func CheckMintime(id, bank int) (int, error) {
	var err error = nil
	db, err := sql.Open("odbc", tiberoDataSourceName)
	if err != nil {
		log.Println("ERROR !! ", err)
		return 0, err
	}
	defer db.Close()

	strQry := `
		select nvl(sum(SI_MINTIME), 0) as mintm from SIGNAL_MAP_T
		where IT_NO = :ID
		and SI_MAPNO = :BANK;
	`

	// log.Printf("qry : %v, %d, %d\n", strQry, id, bank)
	rows, err := db.Query(strQry, id, bank)
	if err != nil {
		log.Printf("CheckMintime err qry : %v, %d, %d\n", strQry, id, bank)
		return 0, err
	}
	defer rows.Close()

	var mintm int = 0
	if rows.Next() {

		err := rows.Scan(&mintm)
		if err != nil {
			return 0, err
		}
	}

	return mintm, nil
}

func ReadSignalMap2(id, bank int) ([]Signalmap, error) {
	var err error = nil
	db, err := sql.Open("odbc", tiberoDataSourceName)
	if err != nil {
		log.Println("ERROR !! ", err)
		return nil, err
	}
	defer db.Close()

	strQry := `
		select IT_NO      
			, SI_MAPNO   
			, SI_RINGNO  
			, SI_STEP    
			, NVL(SI_SSR1_VHE, 0)
			, NVL(SI_SSR1_PED,0)
			, NVL(SI_SSR2_VHE,0)
			, NVL(SI_SSR2_PED,0)
			, NVL(SI_SSR3_VHE,0)
			, NVL(SI_SSR3_PED,0)
			, NVL(SI_SSR4_VHE,0)
			, NVL(SI_SSR4_PED,0)
			, NVL(SI_SSR5_VHE,0)
			, NVL(SI_SSR5_PED,0)
			, NVL(SI_SSR6_VHE,0)
			, NVL(SI_SSR6_PED,0)
			, NVL(SI_SSR7_VHE,0)
			, NVL(SI_SSR7_PED,0)
			, NVL(SI_SSR8_VHE,0)
			, NVL(SI_SSR8_PED,0)
			, NVL(SI_MINTIME ,0)
			, NVL(SI_MAXTIME ,0)
			, NVL(SI_EOP ,0)
		from SIGNAL_MAP_T_TEMP
		where IT_NO = :ID
		and SI_MAPNO = :MAP
		order by IT_NO, SI_MAPNO, SI_RINGNO, SI_STEP;
	`

	// log.Printf("qry : %v, %d, %d\n", strQry, id, bank)
	rows, err := db.Query(strQry, id, bank)
	if err != nil {
		log.Printf("err qry : %v, %d, %d\n", strQry, id, bank)
		return nil, err
	}
	defer rows.Close()

	sigmaps := []Signalmap{}

	for rows.Next() {
		var r Signalmap
		err := rows.Scan(&r.IT_NO, &r.SI_MAPNO, &r.SI_RINGNO, &r.SI_STEP, &r.SI_SSR1_VHE, &r.SI_SSR2_VHE, &r.SI_SSR3_VHE, &r.SI_SSR4_VHE, &r.SI_SSR5_VHE, &r.SI_SSR6_VHE, &r.SI_SSR7_VHE, &r.SI_SSR8_VHE, &r.SI_SSR1_PED, &r.SI_SSR2_PED, &r.SI_SSR3_PED, &r.SI_SSR4_PED, &r.SI_SSR5_PED, &r.SI_SSR6_PED, &r.SI_SSR7_PED, &r.SI_SSR8_PED, &r.SI_MINTIME, &r.SI_MAXTIME, &r.SI_EOP)
		if err != nil {
			return nil, err
		}
		sigmaps = append(sigmaps, r)
	}

	return sigmaps, nil
}

func InsertSignalmap(id, bank int, data []byte) error {
	var err error = nil
	db, err := sql.Open("odbc", tiberoDataSourceName)
	if err != nil {
		log.Println("ERROR !! ", err)
		return err
	}
	defer db.Close()

	strQry := `
		insert into TCS_M_LOC_SIGNAL_MAP (LOC_ID, BANK_ID, DATA)
		VALUES (:ID, :BANK, :DATA);
	`

	// log.Printf("qry : %v\n", strQry)
	stmt, preErr := db.Prepare(strQry)
	if preErr != nil {
		log.Printf("err qry#1 : %v, %d, %d\n", strQry, id, bank)
		return preErr
	}
	_, err = stmt.Exec(id, bank, data)
	if err != nil {
		log.Printf("err qry#2 : %v, %d, %d\n", strQry, id, bank)
		return err
	}

	return nil
}
