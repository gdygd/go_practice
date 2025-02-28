package db

import (
	"fmt"
)

func (dbHand *DBHandler) WriteInput(str string) error {
	db, err := dbHand.Open()
	if err != nil {
		return err
	}
	defer dbHand.Close(db)

	query := fmt.Sprintf(`CALL P_InsertInput('%s')`, str)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
