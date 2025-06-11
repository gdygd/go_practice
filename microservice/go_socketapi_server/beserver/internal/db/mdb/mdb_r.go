package mdb

import (
	"context"
)

func (q *MariaDbHandler) ReadTest(ctx context.Context) (int, error) {
	query := `select 1 from dual;`
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return 0, err
	}

	defer rows.Close()
	var value int = 0
	if rows.Next() {
		if err := rows.Scan(&value); err != nil {
			return 0, err
		}
	}
	if err := rows.Close(); err != nil {
		return 0, err
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return value, nil
}
