package noSql

import (
	"context"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"strings"
	"time"
)

type MongoDB struct {
	Name   string
	Config MongoDBConfig
	app    interfaces.IEngine
	client *mongo.Client
	db     *mongo.Database

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

type MongoDBConfig struct {
	Database string        `yaml:"database"`
	User     string        `yaml:"user,omitempty"`
	Password string        `yaml:"password,omitempty"`
	Hosts    string        `yaml:"hosts,omitempty"`
	Mode     readpref.Mode `yaml:"mode,omitempty"`
	URI      string        `yaml:"uri,omitempty"`
}

func NewMongoDB(name string, cfg MongoDBConfig) *MongoDB {
	return &MongoDB{
		Name:   name,
		Config: cfg,
	}
}

func (t *MongoDB) Init(app interfaces.IEngine) error {
	t.app = app
	var err error
	var opt options.ClientOptions

	if t.Config.URI != "" {
		opt.ApplyURI(t.Config.URI)
	} else {
		opt.Hosts = strings.Split(t.Config.Hosts, ",")
		opt.ReadPreference, _ = readpref.New(t.Config.Mode)

		if t.Config.User != "" {
			opt.Auth = &options.Credential{
				Username:   t.Config.User,
				Password:   t.Config.Password,
				AuthSource: t.Config.Database,
			}
		}
	}

	t.client, err = mongo.Connect(context.TODO(), &opt)

	if err != nil {
		return err
	}

	t.db = t.client.Database(t.Config.Database)

	t.counter, _ = otel.Meter("").Int64Counter("mongodb." + t.Name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("mongodb." + t.Name + ".time")

	return t.client.Ping(context.TODO(), nil)
}

func (t *MongoDB) Stop() error {
	if err := t.client.Disconnect(context.TODO()); err != nil {
		return err
	}

	return nil
}

func (t *MongoDB) String() string {
	return t.Name
}

func (t *MongoDB) Exists(ctx context.Context, key interface{}, query interface{}, args ...interface{}) bool {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *options.CountOptions

	if len(args) > 0 {
		if val, ok := args[0].(*options.CountOptions); ok {
			opt = val
		}
	}

	if opt == nil {
		opt = options.Count()
	}

	var limit int64 = 1
	opt.Limit = &limit

	find, err := t.db.Collection(key.(string)).CountDocuments(ctx, query, opt)

	if err != nil {
		return false
	}

	return find != 0
}

func (t *MongoDB) Count(ctx context.Context, key interface{}, query interface{}, args ...interface{}) int64 {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *options.CountOptions

	if len(args) > 0 {
		if val, ok := args[0].(*options.CountOptions); ok {
			opt = val
		}
	}

	find, err := t.db.Collection(key.(string)).CountDocuments(ctx, query, opt)

	if err != nil {
		return 0
	}

	return find
}

func (t *MongoDB) FindOne(ctx context.Context, key interface{}, query interface{}, args ...interface{}) (interfaces.NoSQLRow, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *options.FindOneOptions

	if len(args) > 0 {
		if val, ok := args[0].(*options.FindOneOptions); ok {
			opt = val
		}
	}

	find := t.db.Collection(key.(string)).FindOne(ctx, query, opt)

	if err := find.Err(); err != nil {
		return nil, err
	}

	return &MongoRow{find}, nil
}

func (t *MongoDB) Find(ctx context.Context, key interface{}, query interface{}, args ...interface{}) (interfaces.NoSQLRows, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *options.FindOptions

	if len(args) > 0 {
		if val, ok := args[0].(*options.FindOptions); ok {
			opt = val
		}
	}

	find, err := t.db.Collection(key.(string)).Find(ctx, query, opt)

	if err != nil {
		return nil, err
	}

	return &MongoRows{find}, nil
}

func (t *MongoDB) Exec(ctx context.Context, key interface{}, query interface{}, args ...interface{}) (interfaces.NoSQLRows, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *options.AggregateOptions

	if len(args) > 0 {
		if val, ok := args[0].(*options.AggregateOptions); ok {
			opt = val
		}
	}

	aggregate, err := t.db.Collection(key.(string)).Aggregate(ctx, query, opt)

	if err != nil {
		return nil, err
	}

	return &MongoRows{aggregate}, nil
}

func (t *MongoDB) Insert(ctx context.Context, key interface{}, query interface{}, args ...interface{}) (interface{}, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *options.InsertOneOptions

	if len(args) > 0 {
		if val, ok := args[0].(*options.InsertOneOptions); ok {
			opt = val
		}
	}

	res, err := t.db.Collection(key.(string)).InsertOne(ctx, query, opt)

	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func (t *MongoDB) Update(ctx context.Context, key interface{}, query interface{}, args ...interface{}) error {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *options.UpdateOptions
	var update interface{}

	if len(args) > 0 {
		update = args[0]

		if len(args) > 1 {
			if val, ok := args[1].(*options.UpdateOptions); ok {
				opt = val
			}
		}
	} else {
		return fmt.Errorf("exec update change is empty")
	}

	_, err := t.db.Collection(key.(string)).UpdateMany(ctx, query, update, opt)

	if err != nil {
		return err
	}

	return nil

}

func (t *MongoDB) Delete(ctx context.Context, key interface{}, query interface{}, args ...interface{}) int64 {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *options.DeleteOptions

	if len(args) > 0 {
		if val, ok := args[0].(*options.DeleteOptions); ok {
			opt = val
		}
	}

	res, err := t.db.Collection(key.(string)).DeleteMany(ctx, query, opt)

	if err != nil {
		return 0
	}

	return res.DeletedCount
}

/*
Batch

key - collection name

query []mongo.WriteModel

typeOp in [bulk]
*/
func (t *MongoDB) Batch(ctx context.Context, key interface{}, query interface{}, typeOp string, args ...interface{}) (int64, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	var opt *options.BulkWriteOptions

	switch typeOp {
	case "bulk":
		write, err := t.db.Collection(key.(string)).BulkWrite(ctx, query.([]mongo.WriteModel), opt)

		if err != nil {
			return 0, err
		}

		return write.ModifiedCount + write.InsertedCount, nil

	default:
		break
	}

	return 0, fmt.Errorf("batch query usupported type")
}

func (t *MongoDB) GetDb() interface{} {
	return t.db
}
