package sqldb

import (
	"fmt"
	"strings"
)

const (
	INT    = "int"
	STRING = "string"
)

const (
	TDBYEAR  = "YYYY"
	TDBMONTH = "MM"
	TDBDAY   = "DD"
	TDBHOUR  = "HH24"
	TDBMIN   = "MI"
	TDBSEC   = "SS"

	MDBYEAR  = "%Y"
	MDBMONTH = "%m"
	MDBDAY   = string("%d")
	MDBHOUR  = "%H"
	MDBMIN   = "%i"
	MDBSEC   = string("%s")
)

// ------------------------------------------------------------------------------
//
//	handleNull
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) handleNull(colnm string, value any) string {
	var nullQry string

	switch m.dbms {
	case ORACLE, TIBERO:
		nullQry = fmt.Sprintf(`nvl(%s, %v)`, colnm, value)
	case MARIADB:
		nullQry = fmt.Sprintf(`ifnull(%s, %v)`, colnm, value)
	}

	return nullQry
}

// ------------------------------------------------------------------------------
//
//	toChar
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) toChar(useSingleQuotes bool, colnm string, isDate, isTime bool, dateDelim, dateTimeDelim, timeDelim string) string {
	// dateDelim: 날짜 구분자           // Ex) YYYY/MM/DD(/), YYYY-MM-DD(-)
	// dateTimeDelim: 날짜-시간 구분자  // Ex) YYYYMMDD HH24MISS(" "), YYYYMMDDHH24MISS("")
	// timeDelim: 시간 구분자          // Ex) HH24:MI:SS(:), HH24MISS("")

	var toCharQry string
	var dateFormat, timeFormat string

	switch m.dbms {
	case ORACLE, TIBERO:
		dateFormat = fmt.Sprintf(`%s%s%s%s%s`, TDBYEAR, dateDelim, TDBMONTH, dateDelim, TDBDAY)
		timeFormat = fmt.Sprintf(`%s%s%s%s%s`, TDBHOUR, timeDelim, TDBMIN, timeDelim, TDBSEC)

		if useSingleQuotes { // 따옴표 사용 → 문자열로 입력시
			if isDate && isTime {
				toCharQry = fmt.Sprintf(`TO_CHAR('%s', '%s%s%s')`, colnm, dateFormat, dateTimeDelim, timeFormat)
			} else if isDate {
				toCharQry = fmt.Sprintf(`TO_CHAR('%s', '%s')`, colnm, dateFormat)
			} else if isTime {
				toCharQry = fmt.Sprintf(`TO_CHAR('%s', '%s')`, colnm, timeFormat)
			}

		} else { // 따옴표 미사용 → 컬럼명 입력시
			if isDate && isTime {
				toCharQry = fmt.Sprintf(`TO_CHAR(%s, '%s%s%s')`, colnm, dateFormat, dateTimeDelim, timeFormat)
			} else if isDate {
				toCharQry = fmt.Sprintf(`TO_CHAR(%s, '%s')`, colnm, dateFormat)
			} else if isTime {
				toCharQry = fmt.Sprintf(`TO_CHAR(%s, '%s')`, colnm, timeFormat)
			}
		}

	case MARIADB:
		dateFormat = fmt.Sprintf(`%s%s%s%s%s`, MDBYEAR, dateDelim, MDBMONTH, dateDelim, MDBDAY)
		timeFormat = fmt.Sprintf(`%s%s%s%s%s`, MDBHOUR, timeDelim, MDBMIN, timeDelim, MDBSEC)

		if useSingleQuotes { // 따옴표 사용 → 문자열로 입력시
			if isDate && isTime {
				toCharQry = fmt.Sprintf(`DATE_FORMAT('%s', '%s%s%s')`, colnm, dateFormat, dateTimeDelim, timeFormat)
			} else if isDate {
				toCharQry = fmt.Sprintf(`DATE_FORMAT('%s', '%s')`, colnm, dateFormat)
			} else if isTime {
				toCharQry = fmt.Sprintf(`DATE_FORMAT('%s', '%s')`, colnm, timeFormat)
			}

		} else { // 따옴표 미사용 → 컬럼명 입력시
			if isDate && isTime {
				toCharQry = fmt.Sprintf(`DATE_FORMAT(%s, '%s%s%s')`, colnm, dateFormat, dateTimeDelim, timeFormat)
			} else if isDate {
				toCharQry = fmt.Sprintf(`DATE_FORMAT(%s, '%s')`, colnm, dateFormat)
			} else if isTime {
				toCharQry = fmt.Sprintf(`DATE_FORMAT(%s, '%s')`, colnm, timeFormat)
			}
		}

	}

	return toCharQry
}

// ------------------------------------------------------------------------------
//
//	toDate
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) toDate(useSingleQuotes bool, colnm, dateDelim, dateTimeDelim, timeDelim string) string {
	// useSingleQuotes: 따옴표 사용 여부
	// dateDelim: 날짜 구분자           // Ex) YYYY/MM/DD(/), YYYY-MM-DD(-)
	// dateTimeDelim: 날짜-시간 구분자  // Ex) YYYYMMDD HH24MISS(" "), YYYYMMDDHH24MISS("")
	// timeDelim: 시간 구분자          // Ex) HH24:MI:SS(:), HH24MISS("")

	var toDateQry string
	var dateFormat, timeFormat string // 날짜,시간형식

	switch m.dbms {
	case ORACLE, TIBERO:
		dateFormat = fmt.Sprintf(`%s%s%s%s%s`, TDBYEAR, dateDelim, TDBMONTH, dateDelim, TDBDAY)
		timeFormat = fmt.Sprintf(`%s%s%s%s%s`, TDBHOUR, timeDelim, TDBMIN, timeDelim, TDBSEC)

		if useSingleQuotes {
			toDateQry = fmt.Sprintf(`TO_DATE('%s', '%s%s%s')`, colnm, dateFormat, dateTimeDelim, timeFormat)
		} else {
			toDateQry = fmt.Sprintf(`TO_DATE(%s, '%s%s%s')`, colnm, dateFormat, dateTimeDelim, timeFormat)
		}

	case MARIADB:
		dateFormat = fmt.Sprintf(`%s%s%s%s%s`, MDBYEAR, dateDelim, MDBMONTH, dateDelim, MDBDAY)
		timeFormat = fmt.Sprintf(`%s%s%s%s%s`, MDBHOUR, timeDelim, MDBMIN, timeDelim, MDBSEC)

		if useSingleQuotes {
			toDateQry = fmt.Sprintf(`STR_TO_DATE('%s', '%s%s%s')`, colnm, dateFormat, dateTimeDelim, timeFormat)
		} else {
			toDateQry = fmt.Sprintf(`STR_TO_DATE(%s, '%s%s%s')`, colnm, dateFormat, dateTimeDelim, timeFormat)
		}
	}

	return toDateQry
}

// ------------------------------------------------------------------------------
//
//	handleAddTime // interval
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) handleAddTime(colnm string, hour, minute, sec int) string {
	var wheTimeQry string

	switch m.dbms {
	case ORACLE, TIBERO:
		wheTimeQry = fmt.Sprintf(`%s + INTERVAL '%d' hour + INTERVAL '%d' minute + INTERVAL '%d' second`, colnm, hour, minute, sec)
	case MARIADB:
		wheTimeQry = fmt.Sprintf(`%s + INTERVAL %d hour + INTERVAL %d minute + INTERVAL %d second`, colnm, hour, minute, sec)
	}

	return wheTimeQry
}

// ------------------------------------------------------------------------------
//
//	handleSubTime // interval
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) handleSubTime(colnm string, hour, minute, sec int) string {
	var subTimeQry string

	switch m.dbms {
	case ORACLE, TIBERO:
		subTimeQry = fmt.Sprintf(`%s - INTERVAL '%d' hour - INTERVAL '%d' minute - INTERVAL '%d' second`, colnm, hour, minute, sec)
	case MARIADB:
		subTimeQry = fmt.Sprintf(`%s - INTERVAL %d hour - INTERVAL %d minute - INTERVAL %d second`, colnm, hour, minute, sec)
	}

	return subTimeQry
}

// ------------------------------------------------------------------------------
//
//	handleRownumStart
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) handleRownumStart(colnms string) string {
	var rnumQry string

	switch m.dbms {
	case ORACLE, TIBERO:
		rnumQry = fmt.Sprintf(`SELECT %s, rnum
								FROM (
										SELECT %s, rownum as rnum
										FROM
										(
								`, colnms, colnms)
	case MARIADB:
		rnumQry = ""
	}

	return rnumQry
}

// ------------------------------------------------------------------------------
//
//	handleRownumEnd
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) handleRownumEnd() string {
	var rnumQry string

	switch m.dbms {
	case ORACLE, TIBERO:
		rnumQry = `		)
					)`
	case MARIADB:
		rnumQry = ""
	}

	return rnumQry
}

// ------------------------------------------------------------------------------
//
//	handleGetOneRow
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) handleGetOneRow(colnm string) string {
	var getOneRowQry string

	switch m.dbms {
	case ORACLE, TIBERO:
		getOneRowQry = fmt.Sprintf(`where %s = 1`, colnm)
	case MARIADB:
		getOneRowQry = fmt.Sprintf(`LIMIT 1`)
	}

	return getOneRowQry
}

// ------------------------------------------------------------------------------
//
//	handleGetNRows
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) handleGetNRows(row_limit, row_offset int) string {
	var getNrowsQry string

	switch m.dbms {
	case ORACLE, TIBERO:
		getNrowsQry = fmt.Sprintf(`where rnum > %d and rnum <= %d`, row_offset, (row_limit + row_offset))
	case MARIADB:
		getNrowsQry = fmt.Sprintf(`LIMIT %d OFFSET %d`, row_limit, row_offset)
	}

	return getNrowsQry
}

// ------------------------------------------------------------------------------
//
//	handleByteLength
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) handleByteLength(colnm string) string {
	var bytelenQry string

	switch m.dbms {
	case ORACLE, TIBERO:
		bytelenQry = fmt.Sprintf(`dbms_lob.getlength(%s)`, colnm)
	case MARIADB:
		bytelenQry = fmt.Sprintf(`octet_length(%s)`, colnm)
	}

	return bytelenQry
}

// ------------------------------------------------------------------------------
//
//	convertType
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) convertType(colnm, datatype string) string {
	// 입력되는 datatype은 Oracle/Tibero를 기준으로 진행
	// Ex) 숫자 : NUMBER
	//     문자 : VARCHAR2
	//     binary data : RAW

	var ctQry string
	if m.dbms == MARIADB {
		switch strings.ToUpper(datatype) {
		case "NUMBER": // NUMBER -> DECIMAL
			datatype = "DECIMAL"
		case "VARCHAR2":
			datatype = "VARCHAR"
		case "RAW":
			datatype = "BLOB"
		}
	}

	ctQry = fmt.Sprintf(`CAST(%s AS %s)`, colnm, datatype)
	return ctQry
}

// ------------------------------------------------------------------------------
//
//	getNowDate
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) getNowDate() string {
	var curTime string

	switch m.dbms {
	case ORACLE, TIBERO:
		curTime = "sysdate"
	case MARIADB:
		curTime = "now()"
	}

	return curTime
}

// ------------------------------------------------------------------------------
//
//	getCurTime
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) getCurTime() string {
	var curTime string

	switch m.dbms {
	case ORACLE, TIBERO:
		curTime = "TO_CHAR(SYSDATE, 'HH24:MI:SS') AS CURRENT_TIME"
	case MARIADB:
		curTime = "curtime()"
	}

	return curTime
}

// ------------------------------------------------------------------------------
//
//	makeInsertEachRow
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) makeInsertEachRow(datas ...any) string {
	// datas: insert할 값들
	var insertEachRowQry string

	var strData []string
	for _, data := range datas {
		strData = append(strData, fmt.Sprintf("%v", data))
	}

	switch m.dbms {
	case ORACLE, TIBERO:
		insertEachRowQry = fmt.Sprintf(`select %s from dual`, strings.Join(strData, ", "))
	case MARIADB:
		insertEachRowQry = fmt.Sprintf(`(%s)`, strings.Join(strData, ", "))
	}

	return insertEachRowQry
}

// ------------------------------------------------------------------------------
//
//	makeInsValQuery
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) makeInsValQuery(insertMulRows []string) string {
	// dbms별로 insert쿼리문에서 VALUES 구문을 추가적으로 작성
	var insValQuery string

	switch m.dbms {
	case ORACLE, TIBERO:
		// Ex)
		// INSERT INTO table_name (column1, column2, ...)
		// SELECT value1, value2, ...
		// UNION ALL
		// SELECT value1, value2, ...;
		insValQuery = strings.Join(insertMulRows, "\n union all \n")

	case MARIADB:
		// Ex)
		// INSERT INTO table_name (column1, column2, ...)
		// VALUES (value1, value2, ...),
		//        (value1, value2, ...),
		//        ...;
		insValQuery = "VALUES " + strings.Join(insertMulRows, ", \n")
	}

	return insValQuery
}

// ------------------------------------------------------------------------------
//
//	nextval
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) nextval(seqNm string) string {
	var nextvalQuery string

	switch m.dbms {
	case ORACLE, TIBERO:
		nextvalQuery = fmt.Sprintf(`%s.nextval`, seqNm)
	case MARIADB:
		nextvalQuery = fmt.Sprintf(`nextval(%s)`, seqNm)
	}

	return nextvalQuery
}

// ------------------------------------------------------------------------------
//
//	currval
//
// ------------------------------------------------------------------------------
func (m *SQLHandler) currval(seqNm string) string {
	var currvalQuery string

	switch m.dbms {
	case ORACLE, TIBERO:
		currvalQuery = fmt.Sprintf(`%s.currval`, seqNm)
	case MARIADB:
		currvalQuery = fmt.Sprintf(`lastval(%s)`, seqNm)
	}

	return currvalQuery
}
