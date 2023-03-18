package utils

import (
	"database/sql"
)

func IsEmptyRow(err error) bool {
	return sql.ErrNoRows == err
}

func IsNotEmptyRow(err error) bool {
	return sql.ErrNoRows != err
}
