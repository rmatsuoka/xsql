# xsql

`xsql` extends the standard library's [`database/sql`](https://pkg.go.dev/database/sql) package with reflection-based helpers for scanning rows into structs and for building named query arguments from structs.

```
go get github.com/rmatsuoka/xsql
```

## How it works

Struct fields are mapped to SQL columns and named parameters using the `db` struct tag. Unexported fields, and exported fields without a `db` tag, are ignored.

```go
type User struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
```

## Scanning rows into structs

[`ScanRows`](https://pkg.go.dev/github.com/rmatsuoka/xsql#ScanRows) scans a `*sql.Rows` into a slice of struct pointers, matching each column to the field whose `db` tag equals the column name:

```go
rows, err := db.Query("select id, name from users")
if err != nil {
	// ...
}
defer rows.Close()

users, err := xsql.ScanRows[User](rows)
```

[`QueryAndScan`](https://pkg.go.dev/github.com/rmatsuoka/xsql#QueryAndScan) combines running a query and scanning its result, closing the rows for you:

```go
users, err := xsql.QueryAndScan[User](ctx, db, "select id, name from users")
```

## Building named arguments

[`NamedArgs`](https://pkg.go.dev/github.com/rmatsuoka/xsql#NamedArgs) builds `[]sql.NamedArg` from a struct's (or pointer to a struct's) tagged fields. [`ToAnys`](https://pkg.go.dev/github.com/rmatsuoka/xsql#ToAnys) converts a typed slice such as `[]sql.NamedArg` to `[]any` for a variadic `args` parameter:

```go
u := User{ID: 3, Name: "rmatsuoka"}

args := xsql.ToAnys(xsql.NamedArgs(u))
_, err := db.Exec("insert into users (id, name) values (@id, @name)", args...)
```

[`NamedArgsWithPrefix`](https://pkg.go.dev/github.com/rmatsuoka/xsql#NamedArgsWithPrefix) prepends a prefix to every argument name, useful when a query binds fields from more than one struct with overlapping column names:

```go
old := User{ID: 3}
updated := User{ID: 3, Name: "rmatsuoka"}

args := append(
	xsql.ToAnys(xsql.NamedArgsWithPrefix("old_", old)),
	xsql.ToAnys(xsql.NamedArgsWithPrefix("new_", updated))...,
)

rows, err := db.QueryContext(ctx,
	`update users set name = @new_name where id = @old_id returning id, name`,
	args...,
)
```
