package db

// 전체 데이터를 보낸다고 일단 가정
func (dbHand *DBHandler) ReadInput() ([]InsertInput, error) {
	db, err := dbHand.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT str, coll_dt FROM insert_input"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	sliceInput := make([]InsertInput, 0)
	for rows.Next() {
		insertInput := InsertInput{}
		err := rows.Scan(&insertInput.Str, &insertInput.Date)
		if err != nil {
			return nil, err
		}

		sliceInput = append(sliceInput, insertInput)
	}

	return sliceInput, nil
}
