package noSql

import (
	"context"
	"fmt"
	"github.com/aerospike/aerospike-client-go/v7"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"strconv"
	"strings"
	"time"
)

type Aerospike struct {
	Name   string
	Config AerospikeConfig
	app    interfaces.IService
	client *aerospike.Client

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

type AerospikeConfig struct {
	User     string
	Password string
	Hosts    string
}

func NewAerospike(name string, cfg AerospikeConfig) *Aerospike {
	return &Aerospike{
		Name:   name,
		Config: cfg,
	}
}

func (t *Aerospike) Init(app interfaces.IService) error {
	t.app = app
	var err error

	cp := aerospike.NewClientPolicy()

	cp.User = t.Config.User
	cp.Password = t.Config.Password
	cp.Timeout = 3 * time.Second
	hostsArr := make([]*aerospike.Host, 0)

	for _, s := range strings.Split(t.Config.Hosts, ",") {
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

	t.counter, _ = otel.Meter("").Int64Counter("aerospike." + t.Name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("aerospike." + t.Name + ".time")

	return nil
}

func (t *Aerospike) Stop() error {
	t.client.Close()

	return nil
}

func (t *Aerospike) String() string {
	return t.Name
}

func (t *Aerospike) Exists(ctx context.Context, key interface{}, query interface{}, args ...interface{}) bool {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *aerospike.BasePolicy

	if len(args) > 1 {
		if val, ok := args[1].(*aerospike.BasePolicy); ok {
			opt = val
		}
	}

	if opt == nil {
		opt = &aerospike.BasePolicy{}
		opt.SendKey = true
	}

	find, err := t.client.Exists(opt, key.(*aerospike.Key))

	if err != nil {
		return false
	}

	return find
}

func (t *Aerospike) Count(ctx context.Context, key interface{}, query interface{}, args ...interface{}) int64 {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	if t.Exists(ctx, key, query, args...) {
		return 1
	}

	return 0
}

func (t *Aerospike) FindOne(ctx context.Context, key interface{}, query interface{}, args ...interface{}) (interfaces.NoSQLRow, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	bins := make([]string, len(args))
	var opt *aerospike.BasePolicy

	if len(args) > 0 {
		if val, ok := args[0].(*aerospike.BasePolicy); ok {
			opt = val
		}
	}

	if opt == nil {
		opt = &aerospike.BasePolicy{}
		opt.SendKey = true
	}

	for i := 0; i < len(args); i++ {
		if val, ok := args[i].(string); ok {
			bins[i] = val
		} else if val, ok := args[i].([]string); ok {
			bins = val
			break
		}
	}

	find, err := t.client.Get(opt, key.(*aerospike.Key), bins...)

	if err != nil {
		return nil, err
	}

	return &ASRow{find}, nil
}

func (t *Aerospike) Find(ctx context.Context, _ interface{}, query interface{}, args ...interface{}) (interfaces.NoSQLRows, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *aerospike.QueryPolicy

	if len(args) > 0 {
		if val, ok := args[0].(*aerospike.QueryPolicy); ok {
			opt = val
		}
	}

	if opt == nil {
		opt = &aerospike.QueryPolicy{}
		opt.SendKey = true
	}

	q, err := t.client.Query(opt, query.(*aerospike.Statement))

	if err != nil {
		return nil, err
	}

	return &ASRows{rows: q}, nil
}

func (t *Aerospike) Exec(ctx context.Context, key interface{}, query interface{}, args ...interface{}) (interfaces.NoSQLRows, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *aerospike.WritePolicy
	var arg1, arg2 string

	if len(args) > 0 {
		if val, ok := args[0].(*aerospike.WritePolicy); ok {
			opt = val
		}
	}

	if opt == nil {
		if len(args) < 2 {
			return nil, fmt.Errorf("invalid query args")
		}

		opt = &aerospike.WritePolicy{}
		opt.SendKey = true

		arg1 = args[0].(string)
		arg2 = args[1].(string)
	} else {
		if len(args) < 3 {
			return nil, fmt.Errorf("invalid query args")
		}

		arg1 = args[1].(string)
		arg2 = args[2].(string)
	}

	_, err := t.client.Execute(opt, key.(*aerospike.Key), arg1, arg2)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *Aerospike) Insert(ctx context.Context, key interface{}, query interface{}, args ...interface{}) (interface{}, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *aerospike.WritePolicy

	if len(args) > 0 {
		if val, ok := args[0].(*aerospike.WritePolicy); ok {
			opt = val
		}
	}

	if opt == nil {
		opt = &aerospike.WritePolicy{}
		opt.SendKey = true
	}

	var bins []*aerospike.Bin

	if val, ok := query.(*aerospike.Bin); ok {
		bins = append(bins, val)
	} else if val, ok := query.([]*aerospike.Bin); ok {
		bins = val
	} else {
		return nil, fmt.Errorf("invalid query type")
	}

	err := t.client.PutBins(opt, key.(*aerospike.Key), bins...)

	if err != nil {
		return nil, err
	}

	return query, nil
}

func (t *Aerospike) Update(ctx context.Context, key interface{}, query interface{}, args ...interface{}) error {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *aerospike.WritePolicy

	if len(args) > 0 {
		if val, ok := args[0].(*aerospike.WritePolicy); ok {
			opt = val
		}
	}

	if opt == nil {
		opt = &aerospike.WritePolicy{}
		opt.SendKey = true
	}

	var bins []*aerospike.Bin

	if val, ok := query.(*aerospike.Bin); ok {
		bins = append(bins, val)
	} else if val, ok := query.([]*aerospike.Bin); ok {
		bins = val
	} else {
		return fmt.Errorf("invalid query type")
	}

	err := t.client.PutBins(opt, key.(*aerospike.Key), bins...)

	if err != nil {
		return err
	}

	return nil

}

func (t *Aerospike) Delete(ctx context.Context, key interface{}, query interface{}, args ...interface{}) int64 {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	if _, ok := key.(*aerospike.Key); ok {
		var opt *aerospike.WritePolicy

		if len(args) > 0 {
			if val, ok := args[0].(*aerospike.WritePolicy); ok {
				opt = val
			}
		}

		if opt == nil {
			opt = &aerospike.WritePolicy{}
			opt.SendKey = true
		}

		_, err := t.client.Delete(opt, key.(*aerospike.Key))

		if err != nil {
			return 0
		}
	} else if keys, ok := key.([]*aerospike.Key); ok {
		var opt *aerospike.BatchPolicy

		if len(args) > 0 {
			if val, ok := args[0].(*aerospike.BatchPolicy); ok {
				opt = val
			}
		}

		if opt == nil {
			opt = &aerospike.BatchPolicy{}
			opt.SendKey = true
		}

		_, err := t.client.BatchDelete(opt, nil, keys)

		if err != nil {
			return 0
		}
	} else {
		return 0
	}

	return 1
}

func (t *Aerospike) Batch(ctx context.Context, key interface{}, query interface{}, typeOp string, args ...interface{}) (int64, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *aerospike.BatchPolicy

	switch typeOp {
	case "operate":
		err := t.client.BatchOperate(opt, query.([]aerospike.BatchRecordIfc))

		if err != nil {
			return 0, err
		}

		return 1, nil

	default:
		break
	}

	return 0, fmt.Errorf("batch query usupported type")
}

func (t *Aerospike) GetDb() interface{} {
	return t.client
}
