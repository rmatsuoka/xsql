// Package xsql extends the standard library's [database/sql] package with
// reflection-based helpers for scanning rows into structs and for building
// named query arguments from structs.
//
// Struct fields are mapped to SQL columns and named parameters using the
// "db" struct tag. Unexported fields, and exported fields without a "db"
// tag, are ignored.
//
//	type User struct {
//		ID   int64  `db:"id"`
//		Name string `db:"name"`
//	}
//
// [ScanRows] scans a [database/sql.Rows] into a slice of struct pointers:
//
//	rows, err := db.Query("select id, name from users")
//	...
//	users, err := xsql.ScanRows[User](rows)
//
// [NamedArgs] builds [database/sql.NamedArg] values from a struct for use as
// named query parameters. [ToAnys] converts them to []any for a variadic
// args parameter:
//
//	args := xsql.ToAnys(xsql.NamedArgs(user)) // {Name: "id", ...}, {Name: "name", ...}
//	_, err := db.Exec("insert into users (id, name) values (@id, @name)", args...)
package xsql
