package elasticsearch8

import (
	"context"
	"os"
	"testing"
)

type testESDoc struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

func TestElastic_SearchSource_IncludesExcludes(t *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t.Skip()
		return
	}
	es := initESTest("test_searchsource8")
	defer closeESTest(es)

	name1 := "sx1"
	name2 := "sx2"
	if _, err := es.Index(context.Background(), es.name, "201", testESDoc{ID: 201, Name: &name1}, true); err != nil {
		t.Fatalf("index err: %v", err)
	}
	if _, err := es.Index(context.Background(), es.name, "202", testESDoc{ID: 202, Name: &name2}, true); err != nil {
		t.Fatalf("index err: %v", err)
	}

	rows, err := es.SearchSource(context.Background(), []string{es.name}, nil, 0, 10, []string{"id"}, []string{"name"})
	if err != nil {
		t.Fatalf("searchsource err: %v", err)
	}
	defer rows.Close()

	cnt := 0
	for rows.Next() {
		var d testESDoc
		if err := rows.Scan(&d); err != nil {
			t.Fatalf("scan err: %v", err)
		}
		if d.ID == 0 {
			t.Fatalf("expected id to be present")
		}
		if d.Name != nil {
			t.Fatalf("expected name to be nil due to excludes, got %v", *d.Name)
		}
		cnt++
	}
	if cnt < 2 {
		t.Fatalf("want >=2 got %d", cnt)
	}

	// sanity: Get should still return full document
	row, err := es.Get(context.Background(), es.name, "201")
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
	var got testESDoc
	if err := row.Scan(&got); err != nil {
		t.Fatalf("scan err: %v", err)
	}
	if got.Name == nil || *got.Name != name1 || got.ID != 201 {
		t.Fatalf("unexpected full doc: %+v", got)
	}
}

func initESTest(name string) *ElasticSearch {
	host := os.Getenv("PARANOIA_INTEGRATED_SERVER")
	es := New(name)
	if err := es.Init(map[string]interface{}{"addresses": "http://" + host + ":9200", "username": "elastic", "password": "changeme"}); err != nil {
		panic(err)
	}
	return es
}

func closeESTest(es *ElasticSearch) { _ = es.Stop() }
