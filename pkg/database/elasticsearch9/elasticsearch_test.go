package elasticsearch9

import (
	"context"
	"os"
	"reflect"
	"testing"
)

type testESDoc struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

func TestElastic_Index_Get_Delete(t *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t.Skip()
		return
	}
	es := initESTest("test_index")
	defer closeESTest(es)

	name := "foo"
	doc := testESDoc{ID: 10, Name: &name}

	id, err := es.Index(context.Background(), es.name, "10", doc, true)
	if err != nil {
		t.Fatalf("index err: %v", err)
	}
	if id == "" {
		t.Fatalf("empty id")
	}

	row, err := es.Get(context.Background(), es.name, "10")
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
	var got testESDoc
	if err := row.Scan(&got); err != nil {
		t.Fatalf("scan err: %v", err)
	}
	if !reflect.DeepEqual(got, doc) {
		t.Fatalf("got %+v want %+v", got, doc)
	}

	if err := es.Delete(context.Background(), es.name, "10", true); err != nil {
		t.Fatalf("delete err: %v", err)
	}
}

func TestElastic_DeleteByQuery(t *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t.Skip()
		return
	}
	es := initESTest("test_index")
	defer closeESTest(es)

	name1 := "t1"
	name2 := "t2"
	es.Index(context.Background(), es.name, "1", testESDoc{ID: 1, Name: &name1}, true)
	es.Index(context.Background(), es.name, "2", testESDoc{ID: 2, Name: &name2}, true)

	rows, err := es.Search(context.Background(), []string{es.name}, nil, 0, 10)
	if err != nil {
		t.Fatalf("search err: %v", err)
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		var d testESDoc
		if err := rows.Scan(&d); err != nil {
			t.Fatalf("scan err: %v", err)
		}
		cnt++
	}
	if cnt < 2 {
		t.Fatalf("want >=2 got %d", cnt)
	}

	if err := es.DeleteByQuery(context.Background(), []string{es.name}, map[string]any{"match_all": map[string]any{}}, true); err != nil {
		t.Fatalf("delete by query err: %v", err)
	}

	rows, err = es.Search(context.Background(), []string{es.name}, nil, 0, 10)
	if err != nil {
		t.Fatalf("search err: %v", err)
	}
	defer rows.Close()
	cnt = 0
	for rows.Next() {
		var d testESDoc
		if err := rows.Scan(&d); err != nil {
			t.Fatalf("scan err: %v", err)
		}
		cnt++
	}
	if cnt != 0 {
		t.Fatalf("want 0 got %d", cnt)
	}
}

func TestElastic_Search(t *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t.Skip()
		return
	}
	es := initESTest("test_search")
	defer closeESTest(es)

	name1 := "t1"
	name2 := "t2"
	es.Index(context.Background(), es.name, "1", testESDoc{ID: 1, Name: &name1}, true)
	es.Index(context.Background(), es.name, "2", testESDoc{ID: 2, Name: &name2}, true)

	rows, err := es.Search(context.Background(), []string{es.name}, nil, 0, 10)
	if err != nil {
		t.Fatalf("search err: %v", err)
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		var d testESDoc
		if err := rows.Scan(&d); err != nil {
			t.Fatalf("scan err: %v", err)
		}
		cnt++
	}
	if cnt < 2 {
		t.Fatalf("want >=2 got %d", cnt)
	}
}

func TestElastic_SearchSource_IncludesExcludes(t *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t.Skip()
		return
	}
	es := initESTest("test_searchsource9")
	defer closeESTest(es)

	name1 := "sx1"
	name2 := "sx2"
	es.Index(context.Background(), es.name, "101", testESDoc{ID: 101, Name: &name1}, true)
	es.Index(context.Background(), es.name, "102", testESDoc{ID: 102, Name: &name2}, true)

	rows, err := es.SearchSource(context.Background(), []string{es.name}, nil, 0, 10, []string{"id"}, []string{"name"})
	if err != nil {
		t.Fatalf("searchsource err: %v", err)
	}
	defer rows.Close()

	// Expect that only id is present in _source, and name is omitted
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
}

func TestElastic_BulkIndex(t *testing.T) {
	if os.Getenv("PARANOIA_INTEGRATED_TESTS") != "Y" {
		t.Skip()
		return
	}
	es := initESTest("test_bulk9")
	defer closeESTest(es)

	name1 := "b1"
	name2 := "b2"
	res, err := es.BulkIndex(context.Background(), es.name, []BulkItem{
		{ID: "b101", Document: testESDoc{ID: 101, Name: &name1}},
		{ID: "b102", Document: testESDoc{ID: 102, Name: &name2}},
	}, true)
	if err != nil {
		t.Fatalf("bulk err: %v", err)
	}
	if len(res.Errors) != 0 {
		t.Fatalf("bulk errors: %+v", res.Errors)
	}
	if len(res.IDs) < 2 {
		t.Fatalf("bulk ids len < 2: %+v", res.IDs)
	}

	// verify documents present
	rows, err := es.Search(context.Background(), []string{es.name}, nil, 0, 10)
	if err != nil {
		t.Fatalf("search err: %v", err)
	}
	defer rows.Close()
	found := 0
	for rows.Next() {
		var d testESDoc
		if err := rows.Scan(&d); err != nil {
			t.Fatalf("scan err: %v", err)
		}
		if (d.ID == 101 && d.Name != nil && *d.Name == name1) || (d.ID == 102 && d.Name != nil && *d.Name == name2) {
			found++
		}
	}
	if found < 2 {
		t.Fatalf("expected to find 2 docs, got %d", found)
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

func closeESTest(es *ElasticSearch) { es.Stop() }
