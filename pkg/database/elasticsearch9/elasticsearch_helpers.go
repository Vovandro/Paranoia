package elasticsearch9

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"

	"github.com/elastic/go-elasticsearch/v9/esapi"
)

type ESRow struct{ res *esapi.Response }

type ESRows struct {
	res *esapi.Response
	dec *json.Decoder
	cur struct {
		hits []json.RawMessage
		idx  int
	}
}

func (r *ESRow) Scan(dest any) error {
	if dest == nil {
		return errors.New("dest is nil")
	}
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return errors.New("dest is not a pointer")
	}
	defer r.res.Body.Close()
	var payload struct {
		Source json.RawMessage `json:"_source"`
	}
	if err := json.NewDecoder(r.res.Body).Decode(&payload); err != nil {
		return err
	}
	if len(payload.Source) == 0 {
		return errors.New("empty _source")
	}
	return json.Unmarshal(payload.Source, dest)
}

func (r *ESRows) Next() bool {
	if r.dec == nil {
		r.dec = json.NewDecoder(r.res.Body)
	}
	if r.cur.hits != nil && r.cur.idx < len(r.cur.hits) {
		return true
	}
	var payload struct {
		Hits struct {
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := r.dec.Decode(&payload); err != nil {
		return false
	}
	r.cur.hits = make([]json.RawMessage, len(payload.Hits.Hits))
	for i, h := range payload.Hits.Hits {
		r.cur.hits[i] = h.Source
	}
	r.cur.idx = 0
	return len(r.cur.hits) > 0
}

func (r *ESRows) Scan(dest any) error {
	if dest == nil {
		return errors.New("dest is nil")
	}
	if reflect.TypeOf(dest).Kind() != reflect.Ptr {
		return errors.New("dest is not a pointer")
	}
	if r.cur.idx >= len(r.cur.hits) {
		return io.EOF
	}
	err := json.Unmarshal(r.cur.hits[r.cur.idx], dest)
	r.cur.idx++
	return err
}

func (r *ESRows) Close() error {
	if r.res != nil && r.res.Body != nil {
		return r.res.Body.Close()
	}
	return nil
}
