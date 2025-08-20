package elasticsearch9

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	es9 "github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type ElasticSearch struct {
	name   string
	config Config
	client *es9.Client

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

type Config struct {
	Addresses string `yaml:"addresses"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	CloudID   string `yaml:"cloud_id"`
	APIKey    string `yaml:"api_key"`
}

func New(name string) *ElasticSearch { return &ElasticSearch{name: name} }

func (t *ElasticSearch) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Addresses == "" && t.config.CloudID == "" {
		return errors.New("addresses or cloud_id is required")
	}

	esCfg := es9.Config{}
	if t.config.CloudID != "" {
		esCfg.CloudID = t.config.CloudID
	} else {
		esCfg.Addresses = strings.Split(t.config.Addresses, ",")
	}
	if t.config.APIKey != "" {
		esCfg.APIKey = t.config.APIKey
	} else if t.config.Username != "" {
		esCfg.Username = t.config.Username
		esCfg.Password = t.config.Password
	}

	t.client, err = es9.NewClient(esCfg)
	if err != nil {
		return err
	}

	t.counter, _ = otel.Meter("").Int64Counter("elasticsearch." + t.name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("elasticsearch." + t.name + ".time")

	resp, err := t.client.Info()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		body, _ := io.ReadAll(resp.Body)
		return errors.New(string(body))
	}
	return nil
}

func (t *ElasticSearch) Stop() error  { return nil }
func (t *ElasticSearch) Name() string { return t.name }
func (t *ElasticSearch) Type() string { return "database" }

func (t *ElasticSearch) Index(ctx context.Context, index string, id string, document interface{}, refresh bool) (string, error) {
	defer func(s time.Time) { t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds()) }(time.Now())
	t.counter.Add(context.Background(), 1)
	body, err := json.Marshal(document)
	if err != nil {
		return "", err
	}
	req := esapi.IndexRequest{Index: index, DocumentID: id, Body: bytes.NewReader(body), Refresh: func() string {
		if refresh {
			return "true"
		} else {
			return "false"
		}
	}()}
	res, err := req.Do(ctx, t.client)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.IsError() {
		data, _ := io.ReadAll(res.Body)
		return "", errors.New(string(data))
	}
	var out struct {
		ID string `json:"_id"`
	}
	_ = json.NewDecoder(res.Body).Decode(&out)
	if out.ID != "" {
		return out.ID, nil
	}
	return id, nil
}

func (t *ElasticSearch) Get(ctx context.Context, index string, id string) (NoSQLRow, error) {
	defer func(s time.Time) { t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds()) }(time.Now())
	t.counter.Add(context.Background(), 1)
	res, err := t.client.Get(index, id, t.client.Get.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		data, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		return nil, errors.New(string(data))
	}
	return &ESRow{res: res}, nil
}

func (t *ElasticSearch) Search(ctx context.Context, index []string, query map[string]any, from, size int) (NoSQLRows, error) {
	defer func(s time.Time) { t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds()) }(time.Now())
	t.counter.Add(context.Background(), 1)
	body := map[string]any{"from": from, "size": size}
	if query != nil {
		body["query"] = query
	} else {
		body["query"] = map[string]any{"match_all": map[string]any{}}
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	res, err := t.client.Search(t.client.Search.WithContext(ctx), t.client.Search.WithIndex(index...), t.client.Search.WithBody(bytes.NewReader(b)))
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		data, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		return nil, errors.New(string(data))
	}
	return &ESRows{res: res}, nil
}

func (t *ElasticSearch) Delete(ctx context.Context, index string, id string, refresh bool) error {
	defer func(s time.Time) { t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds()) }(time.Now())
	t.counter.Add(context.Background(), 1)
	req := esapi.DeleteRequest{Index: index, DocumentID: id, Refresh: func() string {
		if refresh {
			return "true"
		} else {
			return "false"
		}
	}()}
	res, err := req.Do(ctx, t.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		data, _ := io.ReadAll(res.Body)
		return errors.New(string(data))
	}
	return nil
}

func (t *ElasticSearch) Update(ctx context.Context, index string, id string, doc interface{}, refresh bool) error {
	defer func(s time.Time) { t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds()) }(time.Now())
	t.counter.Add(context.Background(), 1)
	body := map[string]any{"doc": doc}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req := esapi.UpdateRequest{Index: index, DocumentID: id, Body: bytes.NewReader(b), Refresh: func() string {
		if refresh {
			return "true"
		} else {
			return "false"
		}
	}()}
	res, err := req.Do(ctx, t.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		data, _ := io.ReadAll(res.Body)
		return errors.New(string(data))
	}
	return nil
}

func (t *ElasticSearch) GetClient() interface{} { return t.client }
