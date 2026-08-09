package xsql

// ToAnys converts s to a []any, copying each element. It is useful for
// passing a slice of a concrete element type, such as the []sql.NamedArg
// returned by [NamedArgs], to a variadic ...any parameter such as
// [database/sql.DB.QueryContext]'s args.
func ToAnys[S ~[]E, E any](s S) []any {
	as := make([]any, len(s))
	for i := range s {
		as[i] = s[i]
	}
	return as
}
