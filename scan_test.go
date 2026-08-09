package xsql

import (
	"testing"
	"time"

	gocmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestScanRows(t *testing.T) {
	db, mock := mustSQLMockNew()
	now := time.Now()

	type T struct {
		Date time.Time `db:"date"`
	}

	type S struct {
		ID           int64  `db:"id"`
		Name         string `db:"name"`
		shouldIgnore int    `db:"ignore"`
		T                   // embed
	}

	mockRows := mock.NewRows([]string{"id", "name", "date", "ignore"}).
		AddRow(1, "name_1", now, 5).
		AddRow(2, "name_2", now, 5)

	mock.ExpectQuery("select").WillReturnRows(mockRows)

	rows, _ := db.Query("select")

	got, err := ScanRows[S](rows)
	if err != nil {
		t.Error(err)
	}
	diff := gocmp.Diff(got, []*S{
		{ID: 1, Name: "name_1", T: T{Date: now}},
		{ID: 2, Name: "name_2", T: T{Date: now}},
	}, cmpopts.IgnoreUnexported(S{}))

	if diff != "" {
		t.Error(diff)
	}
}
