package xsql

import (
	"database/sql"
	"reflect"
)

// NamedArgs returns a [database/sql.NamedArg] for every exported field of s
// that has a "db" struct tag. Each argument's name is the tag value and its
// value is the field's value.
//
// s must be a struct or a pointer to a struct; NamedArgs panics otherwise,
// including when s is a nil pointer.
func NamedArgs(s any) []sql.NamedArg {
	return NamedArgsWithPrefix("", s)
}

// NamedArgsWithPrefix is like [NamedArgs] but prepends prefix to every
// argument name. It is useful for disambiguating parameters when a query
// binds fields from more than one struct.
//
// s must be a struct or a pointer to a struct; NamedArgsWithPrefix panics
// otherwise, including when s is a nil pointer.
func NamedArgsWithPrefix(prefix string, s any) []sql.NamedArg {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic("s must be struct or pointer to struct")
	}
	fields := reflect.VisibleFields(v.Type())
	var args []sql.NamedArg
	for _, f := range fields {
		if !f.IsExported() {
			continue
		}
		name, ok := f.Tag.Lookup(structTagKey)
		if !ok {
			continue
		}
		value := v.FieldByIndex(f.Index).Interface()
		args = append(args, sql.Named(prefix+name, value))
	}
	return args
}
