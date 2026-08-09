package xsql

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

func mustSQLMockNew() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return db, mock
}
