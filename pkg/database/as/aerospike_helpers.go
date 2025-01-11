package as

import (
	"github.com/aerospike/aerospike-client-go/v7"
	"gitlab.com/devpro_studio/go_utils/decode"
)

type ASRow struct {
	row *aerospike.Record
}

type ASRows struct {
	rows *aerospike.Recordset
	row  *aerospike.Result
}

func (t *ASRow) Scan(dest any) error {
	if _, ok := dest.(map[string]interface{}); ok {
		for k, v := range t.row.Bins {
			dest.(map[string]interface{})[k] = v
		}
	} else {
		err := decode.Decode(t.row.Bins, &dest, "db", 0)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *ASRows) Next() bool {
	var ok bool
	t.row, ok = <-t.rows.Results()
	return ok && t.row.Err == nil
}

func (t *ASRows) Scan(dest any) error {
	if _, ok := dest.(map[string]interface{}); ok {
		for k, v := range t.row.Record.Bins {
			dest.(map[string]interface{})[k] = v
		}
	} else {
		err := decode.Decode(t.row.Record.Bins, &dest, "db", 0)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *ASRows) Close() error {
	return t.rows.Close()
}
