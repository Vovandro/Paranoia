package mongodb

import (
	"context"
	"errors"
	"gitlab.com/devpro_studio/go_utils/decode"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"strings"
	"time"
)

type MongoDB struct {
	name   string
	config Config
	client *mongo.Client
	db     *mongo.Database

	counter     metric.Int64Counter
	timeCounter metric.Int64Histogram
}

type Config struct {
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Hosts    string `yaml:"hosts"`
	Mode     string `yaml:"mode"`
	URI      string `yaml:"uri"`
}

func New(name string) *MongoDB {
	return &MongoDB{
		name: name,
	}
}

func (t *MongoDB) Init(cfg map[string]interface{}) error {
	err := decode.Decode(cfg, &t.config, "yaml", decode.DecoderStrongFoundDst)
	if err != nil {
		return err
	}

	if t.config.Hosts == "" && t.config.URI == "" {
		return errors.New("hosts or uri is required")
	}

	var opt options.ClientOptions

	if t.config.URI != "" {
		opt.ApplyURI(t.config.URI)
	} else {
		opt.Hosts = strings.Split(t.config.Hosts, ",")
		modeConverted, err := readpref.ModeFromString(t.config.Mode)
		if err != nil {
			return err
		}

		opt.ReadPreference, _ = readpref.New(modeConverted)

		if t.config.User != "" {
			opt.Auth = &options.Credential{
				Username:   t.config.User,
				Password:   t.config.Password,
				AuthSource: t.config.Database,
			}
		}
	}

	t.client, err = mongo.Connect(context.TODO(), &opt)

	if err != nil {
		return err
	}

	t.db = t.client.Database(t.config.Database)

	t.counter, _ = otel.Meter("").Int64Counter("mongodb." + t.name + ".count")
	t.timeCounter, _ = otel.Meter("").Int64Histogram("mongodb." + t.name + ".time")

	return t.client.Ping(context.TODO(), nil)
}

func (t *MongoDB) Stop() error {
	if err := t.client.Disconnect(context.TODO()); err != nil {
		return err
	}

	return nil
}

func (t *MongoDB) Name() string {
	return t.name
}

func (t *MongoDB) Type() string {
	return "database"
}

func (t *MongoDB) Exists(ctx context.Context, collection string, query interface{}) bool {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	opt := options.Count()
	var limit int64 = 1
	opt.Limit = &limit

	find, err := t.db.Collection(collection).CountDocuments(ctx, query, opt)

	if err != nil {
		return false
	}

	return find != 0
}

func (t *MongoDB) Count(ctx context.Context, collection string, query interface{}, opt *options.CountOptions) int64 {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	find, err := t.db.Collection(collection).CountDocuments(ctx, query, opt)

	if err != nil {
		return 0
	}

	return find
}

func (t *MongoDB) FindOne(ctx context.Context, collection string, query interface{}, opt *options.FindOneOptions) (NoSQLRow, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	find := t.db.Collection(collection).FindOne(ctx, query, opt)

	if err := find.Err(); err != nil {
		return nil, err
	}

	return &MongoRow{find}, nil
}

func (t *MongoDB) FindOneAndUpdate(ctx context.Context, collection string, query interface{}, update interface{}, opt *options.FindOneAndUpdateOptions) (NoSQLRow, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	find := t.db.Collection(collection).FindOneAndUpdate(ctx, query, update, opt)

	if err := find.Err(); err != nil {
		return nil, err
	}

	return &MongoRow{find}, nil
}

func (t *MongoDB) Find(ctx context.Context, collection string, query interface{}, opt *options.FindOptions) (NoSQLRows, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	find, err := t.db.Collection(collection).Find(ctx, query, opt)

	if err != nil {
		return nil, err
	}

	return &MongoRows{find}, nil
}

func (t *MongoDB) Exec(ctx context.Context, collection string, query interface{}, opt *options.AggregateOptions) (NoSQLRows, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	aggregate, err := t.db.Collection(collection).Aggregate(ctx, query, opt)

	if err != nil {
		return nil, err
	}

	return &MongoRows{aggregate}, nil
}

func (t *MongoDB) Insert(ctx context.Context, collection string, query interface{}, opt *options.InsertOneOptions) (interface{}, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	res, err := t.db.Collection(collection).InsertOne(ctx, query, opt)

	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func (t *MongoDB) Update(ctx context.Context, collection string, query interface{}, update interface{}, opt *options.UpdateOptions) error {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	_, err := t.db.Collection(collection).UpdateMany(ctx, query, update, opt)

	if err != nil {
		return err
	}

	return nil

}

func (t *MongoDB) Delete(ctx context.Context, collection string, query interface{}, opt *options.DeleteOptions) int64 {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	res, err := t.db.Collection(collection).DeleteMany(ctx, query, opt)

	if err != nil {
		return 0
	}

	return res.DeletedCount
}

func (t *MongoDB) Batch(ctx context.Context, collection string, query []mongo.WriteModel, opt *options.BulkWriteOptions) (int64, error) {
	defer func(s time.Time) {
		t.timeCounter.Record(context.Background(), time.Since(s).Milliseconds())
	}(time.Now())
	t.counter.Add(context.Background(), 1)

	write, err := t.db.Collection(collection).BulkWrite(ctx, query, opt)

	if err != nil {
		return 0, err
	}

	return write.ModifiedCount + write.InsertedCount + write.UpsertedCount, nil
}

func (t *MongoDB) GetDb() *mongo.Database {
	return t.db
}
