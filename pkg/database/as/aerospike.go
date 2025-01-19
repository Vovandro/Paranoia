package as

import (
	"context"
	"errors"
	"fmt"
	"github.com/aerospike/aerospike-client-go/v7"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"strconv"
	"strings"
	"time"
)

type Aerospike struct {
	name   string
	config Config
	client *aerospike.Client

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Hosts    string `yaml:"hosts"`
}

func New(name string) *Aerospike {
	return &Aerospike{
		name: name,
	}
}

func (t *Aerospike) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Hosts == "" {
		return errors.New("hosts is required")
	}

	cp := aerospike.NewClientPolicy()

	cp.User = t.config.User
	cp.Password = t.config.Password
	cp.Timeout = 3 * time.Second
	hostsArr := make([]*aerospike.Host, 0)

	for _, s := range strings.Split(t.config.Hosts, ",") {
		item := strings.Split(s, ":")
		p, _ := strconv.ParseInt(item[1], 10, 64)
		hostsArr = append(hostsArr, aerospike.NewHost(item[0], int(p)))
	}

	t.client, err = aerospike.NewClientWithPolicyAndHost(cp, hostsArr...)

	if err != nil {
		return err
	}

	var policy aerospike.BasePolicy
	var wPolicy aerospike.WritePolicy
	var bPolicy aerospike.BatchPolicy

	policy.SendKey = true
	wPolicy.SendKey = true
	bPolicy.SendKey = true

	t.client.DefaultPolicy = &policy
	t.client.DefaultWritePolicy = &wPolicy
	t.client.DefaultBatchPolicy = &bPolicy

	t.counter, _ = otel.Meter("").Int64Counter("aerospike." + t.name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("aerospike." + t.name + ".time")

	return nil
}

func (t *Aerospike) Stop() error {
	t.client.Close()

	return nil
}

func (t *Aerospike) Name() string {
	return t.name
}

func (t *Aerospike) Type() string {
	return "database"
}

func (t *Aerospike) Exists(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) bool {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	if policy == nil {
		policy = &aerospike.BasePolicy{}
		policy.SendKey = true
	}

	find, err := t.client.Exists(policy, key)

	if err != nil {
		return false
	}

	return find
}

func (t *Aerospike) Count(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy) int64 {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	if t.Exists(ctx, key, policy) {
		return 1
	}

	return 0
}

func (t *Aerospike) FindOne(ctx context.Context, key *aerospike.Key, policy *aerospike.BasePolicy, bins []string) (NoSQLRow, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	if policy == nil {
		policy = &aerospike.BasePolicy{}
		policy.SendKey = true
	}

	find, err := t.client.Get(policy, key, bins...)

	if err != nil {
		return nil, err
	}

	return &ASRow{find}, nil
}

func (t *Aerospike) Find(ctx context.Context, query *aerospike.Statement, policy *aerospike.QueryPolicy) (NoSQLRows, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	if policy == nil {
		policy = &aerospike.QueryPolicy{}
		policy.SendKey = true
	}

	q, err := t.client.Query(policy, query)

	if err != nil {
		return nil, err
	}

	return &ASRows{rows: q}, nil
}

func (t *Aerospike) Exec(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy, packageName string, functionName string) (NoSQLRows, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	if policy == nil {
		policy = &aerospike.WritePolicy{}
		policy.SendKey = true
	}

	_, err := t.client.Execute(policy, key, packageName, functionName)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Insert query is *aerospike.Bin or []*aerospike.Bin or *aerospike.BinMap
func (t *Aerospike) Insert(ctx context.Context, key *aerospike.Key, query interface{}, policy *aerospike.WritePolicy) (interface{}, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	if policy == nil {
		policy = &aerospike.WritePolicy{}
		policy.SendKey = true
	}

	var err error

	if val, ok := query.(*aerospike.Bin); ok {
		err = t.client.PutBins(policy, key, val)
	} else if val, ok := query.([]*aerospike.Bin); ok {
		err = t.client.PutBins(policy, key, val...)
	} else if val, ok := query.(*aerospike.BinMap); ok {
		err = t.client.Put(policy, key, *val)
	} else {
		return nil, fmt.Errorf("invalid query type")
	}

	if err != nil {
		return nil, err
	}

	return query, nil
}

func (t *Aerospike) Delete(ctx context.Context, key *aerospike.Key, policy *aerospike.WritePolicy) int64 {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	if policy == nil {
		policy = &aerospike.WritePolicy{}
		policy.SendKey = true
	}

	_, err := t.client.Delete(policy, key)

	if err != nil {
		return 0
	}

	return 1
}

func (t *Aerospike) DeleteMany(ctx context.Context, keys []*aerospike.Key, policy *aerospike.BatchPolicy, policyDelete *aerospike.BatchDeletePolicy) int64 {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	if policy == nil {
		policy = &aerospike.BatchPolicy{}
		policy.SendKey = true
	}

	_, err := t.client.BatchDelete(policy, policyDelete, keys)

	if err != nil {
		return 0
	}

	return 1
}

func (t *Aerospike) Operate(ctx context.Context, query []aerospike.BatchRecordIfc) (int64, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *aerospike.BatchPolicy

	err := t.client.BatchOperate(opt, query)

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (t *Aerospike) GetDb() *aerospike.Client {
	return t.client
}
