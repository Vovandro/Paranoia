package noSql

import (
	"context"
	"fmt"
	"github.com/aerospike/aerospike-client-go/v7"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Aerospike struct {
	Name   string
	Hosts  string
	User   string
	Pass   string
	app    interfaces.IService
	client *aerospike.Client
}

func (t *Aerospike) Init(app interfaces.IService) error {
	t.app = app
	var err error

	cp := aerospike.NewClientPolicy()

	cp.User = t.User
	cp.Password = t.Pass
	cp.Timeout = 3 * time.Second
	hostsArr := make([]*aerospike.Host, 0)

	for _, s := range strings.Split(t.Hosts, ",") {
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

	return nil
}

func (t *Aerospike) Stop() error {
	t.client.Close()

	return nil
}

func (t *Aerospike) String() string {
	return t.Name
}

func (t *Aerospike) Exists(ctx context.Context, query interface{}, args ...interface{}) bool {
	var opt *aerospike.BasePolicy
	var key *aerospike.Key

	if len(args) > 0 {
		key = args[0].(*aerospike.Key)
	} else {
		return false
	}

	if len(args) > 1 {
		if val, ok := args[1].(*aerospike.BasePolicy); ok {
			opt = val
		}
	}

	if opt == nil {
		opt = &aerospike.BasePolicy{}
		opt.SendKey = true
	}

	find, err := t.client.Exists(opt, key)

	if err != nil {
		return false
	}

	return find
}

func (t *Aerospike) Count(ctx context.Context, query interface{}, args ...interface{}) int64 {
	return 0
}

func (t *Aerospike) FindOne(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	if key, ok := query.(*aerospike.Key); ok {
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

		find, err := t.client.Get(opt, key, bins...)

		if err != nil {
			return err
		}

		if _, ok := model.(map[string]interface{}); ok {
			for k, v := range find.Bins {
				model.(map[string]interface{})[k] = v
			}
		} else {
			err := t.Scan(find.Bins, model)

			if err != nil {
				return err
			}
		}
	} else if q, ok := query.(*aerospike.Statement); ok {
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

		query, err := t.client.Query(opt, q)

		if err != nil {
			return err
		}
		defer query.Close()

		res := <-query.Results()

		e := t.Scan(res.Record.Bins, model)

		if e != nil {
			return e
		}
	} else {
		return fmt.Errorf("invalid query type: %T", query)
	}

	return nil
}

func (t *Aerospike) Find(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
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
		return err
	}
	defer q.Close()

	for r := range q.Results() {
		item := reflect.New(reflect.TypeOf(model))
		e := t.Scan(r.Record.Bins, item)

		if e != nil {
			return e
		}

		model = append(model.([]interface{}), item.Interface())
	}

	return nil
}

func (t *Aerospike) Exec(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	var opt *aerospike.WritePolicy
	var arg1, arg2 string

	if len(args) > 0 {
		if val, ok := args[0].(*aerospike.WritePolicy); ok {
			opt = val
		}
	}

	if opt == nil {
		if len(args) < 2 {
			return fmt.Errorf("invalid query args")
		}

		opt = &aerospike.WritePolicy{}
		opt.SendKey = true

		arg1 = args[0].(string)
		arg2 = args[1].(string)
	} else {
		if len(args) < 3 {
			return fmt.Errorf("invalid query args")
		}

		arg1 = args[1].(string)
		arg2 = args[2].(string)
	}

	q, err := t.client.Execute(opt, query.(*aerospike.Key), arg1, arg2)

	if err != nil {
		return err
	}

	model = q

	return nil
}

func (t *Aerospike) Insert(ctx context.Context, query interface{}, args ...interface{}) (interface{}, error) {
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

	for i := 0; i < len(args); i++ {
		if val, ok := args[i].(*aerospike.Bin); ok {
			bins[i] = val
		} else if val, ok := args[i].([]*aerospike.Bin); ok {
			bins = val
			break
		}
	}

	err := t.client.PutBins(opt, query.(*aerospike.Key), bins...)

	if err != nil {
		return nil, err
	}

	return query, nil
}

func (t *Aerospike) Update(ctx context.Context, query interface{}, args ...interface{}) error {
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

	for i := 0; i < len(args); i++ {
		if val, ok := args[i].(*aerospike.Bin); ok {
			bins[i] = val
		} else if val, ok := args[i].([]*aerospike.Bin); ok {
			bins = val
			break
		}
	}

	err := t.client.PutBins(opt, query.(*aerospike.Key), bins...)

	if err != nil {
		return err
	}

	return nil

}

func (t *Aerospike) Delete(ctx context.Context, query interface{}, args ...interface{}) int64 {
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

	_, err := t.client.Delete(opt, query.(*aerospike.Key))

	if err != nil {
		return 0
	}

	return 1
}

func (t *Aerospike) Batch(ctx context.Context, typeOp string, query interface{}, args ...interface{}) (int64, error) {
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

func (t *Aerospike) Scan(from aerospike.BinMap, to interface{}) error {
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
