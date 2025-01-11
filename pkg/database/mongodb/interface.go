package mongodb

type NoSQLRow interface {
	Scan(dest any) error
}

type NoSQLRows interface {
	Next() bool
	Scan(dest any) error
	Close() error
}
