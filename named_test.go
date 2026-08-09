package xsql

import (
	"database/sql"
	"testing"

	gocmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNamedArgsWithPrefix(t *testing.T) {
	type S struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}

	opts := []gocmp.Option{
		cmpopts.SortSlices(func(x, y sql.NamedArg) bool { return x.Name < y.Name }),
		gocmp.AllowUnexported(sql.NamedArg{}),
	}

	t.Run("struct", func(t *testing.T) {
		nameds := NamedArgsWithPrefix("p_", S{ID: 3, Name: "name_3"})
		diff := gocmp.Diff(
			[]sql.NamedArg{sql.Named("p_id", int64(3)), sql.Named("p_name", "name_3")},
			nameds,
			opts...,
		)
		if diff != "" {
			t.Error(diff)
		}
	})

	t.Run("pointer", func(t *testing.T) {
		nameds := NamedArgsWithPrefix("p_", &S{ID: 3, Name: "name_3"})
		diff := gocmp.Diff(
			[]sql.NamedArg{sql.Named("p_id", int64(3)), sql.Named("p_name", "name_3")},
			nameds,
			opts...,
		)
		if diff != "" {
			t.Error(diff)
		}
	})
}
