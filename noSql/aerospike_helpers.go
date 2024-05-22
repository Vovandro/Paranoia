package noSql

import (
	"github.com/aerospike/aerospike-client-go/v7"
	"reflect"
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
		err := asScan(t.row.Bins, dest)

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
		err := asScan(t.row.Record.Bins, dest)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *ASRows) Close() error {
	return t.rows.Close()
}

func asScan(from aerospike.BinMap, to interface{}) error {
	vv := reflect.TypeOf(to)
	vv2 := reflect.ValueOf(to)

	for i := 0; i < vv.NumField(); i++ {
		tag, ok2 := vv.Field(i).Tag.Lookup("db")
		if ok2 {
			if v, ok3 := from[tag]; ok3 {
				vv2.Field(i).Send(reflect.ValueOf(v))
			}
		}
	}

	return nil
}
