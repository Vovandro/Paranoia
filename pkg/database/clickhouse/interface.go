package clickhouse

type SQLRow interface {
	Scan(dest ...any) error
}

type SQLRows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}
