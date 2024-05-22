package database

import "github.com/jackc/pgx/v5"

type PGSQLRows struct {
	Rows pgx.Rows
}

func (t *PGSQLRows) Next() bool {
	return t.Rows.Next()
}

func (t *PGSQLRows) Scan(dest ...any) error {
	return t.Rows.Scan(dest...)
}

func (t *PGSQLRows) Close() error {
	t.Rows.Close()
	return nil
}
