package xsql_test

import (
	"context"
	"database/sql"

	"github.com/rmatsuoka/xsql"
)

func Example() {
	ctx := context.Background()
	var db *sql.DB

	type User struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}

	u := User{ID: 3, Name: "rmatsuoka"}

	args := xsql.ToAnys(xsql.NamedArgs(u))
	_, err := db.Exec("insert users (id, name) values (@id, @name)", args...)
	if err != nil {
		panic(err)
	}

	users, err := xsql.QueryAndScan[User](ctx, db, "select id, name from users where id = ?", 3)
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		println(user.Name)
	}
}

func ExampleNamedArgsWithPrefix() {
	ctx := context.Background()
	var db *sql.DB

	type User struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}

	// old and new share the same "id"/"name" columns, so each is bound
	// with its own prefix to keep the named parameters distinct.
	old := User{ID: 3}
	updated := User{ID: 3, Name: "rmatsuoka"}

	args := append(
		xsql.ToAnys(xsql.NamedArgsWithPrefix("old_", old)),
		xsql.ToAnys(xsql.NamedArgsWithPrefix("new_", updated))...,
	)

	users, err := xsql.QueryAndScan[User](ctx, db,
		`update users set name = @new_name where id = @old_id returning id, name`,
		args...,
	)
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		println(user.Name)
	}
}
