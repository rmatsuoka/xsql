package xsql

import (
	"cmp"
	"context"
	"database/sql"
	"reflect"
)

const structTagKey = "db"

// Querier is the query method required by [QueryAndScan]. *[database/sql.DB],
// *[database/sql.Conn], and *[database/sql.Tx] all satisfy Querier.
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// QueryAndScan runs query against db and scans the resulting rows into a
// slice of *T using [ScanRows]. The rows are always closed before
// QueryAndScan returns, even on error.
func QueryAndScan[T any](ctx context.Context, db Querier, query string, args ...any) ([]*T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	list, scanErr := ScanRows[T](rows)
	closeErr := rows.Close()
	return list, cmp.Or(scanErr, closeErr)
}

// ScanRows scans r into a slice of *T, matching each column to the exported
// field of T whose "db" struct tag equals the column name. Columns with no
// matching field, and fields without a "db" tag, are ignored.
//
// ScanRows does not close r; the caller remains responsible for closing it,
// for example with defer r.Close() or via [QueryAndScan].
//
// T must be a struct type; ScanRows panics otherwise.
func ScanRows[T any](r *sql.Rows) ([]*T, error) {
	typ := reflect.TypeFor[T]()
	if typ.Kind() != reflect.Struct {
		panic("T must be struct")
	}
	cs, err := r.Columns()
	if err != nil {
		return nil, err
	}
	indexes := fieldIndexesFromColumns(typ, cs)
	var ts []*T
	for r.Next() {
		t := new(T)
		dist := pointersFromFieldIndexes(t, indexes)
		if err := r.Scan(dist...); err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	return ts, r.Err()
}

type discard struct{}

func (discard) Scan(any) error {
	return nil
}

func fieldIndexesFromColumns(typ reflect.Type, columns []string) [][]int {
	fields := reflect.VisibleFields(typ)
	indexes := make([][]int, len(columns))
	for i, c := range columns {
		f, ok := findStructFieldByTag(fields, structTagKey, c)
		if !ok || !f.IsExported() {
			continue
		}
		indexes[i] = f.Index
	}
	return indexes
}

func pointersFromFieldIndexes(t any, indexes [][]int) []any {
	vp := reflect.ValueOf(t).Elem()
	dist := make([]any, len(indexes))
	for i, index := range indexes {
		if len(index) == 0 {
			dist[i] = discard{}
			continue
		}
		p := vp.FieldByIndex(index)
		dist[i] = p.Addr().Interface()
	}
	return dist
}

func findStructFieldByTag(fs []reflect.StructField, key, tag string) (reflect.StructField, bool) {
	for _, f := range fs {
		if tag == f.Tag.Get(key) {
			return f, true
		}
	}
	return reflect.StructField{}, false
}
