package noSql

import (
	"context"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.mongodb.org/mongo-driver/bson"
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

func (t *MongoDB) Exists(ctx context.Context, collection string, query bson.D) bool {
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

func (t *MongoDB) Count(ctx context.Context, collection string, query bson.D, opt *options.CountOptions) int64 {
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

func (t *MongoDB) FindOne(ctx context.Context, collection string, query bson.D, opt *options.FindOneOptions) (interfaces.NoSQLRow, error) {
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

func (t *MongoDB) Find(ctx context.Context, collection string, query bson.D, opt *options.FindOptions) (interfaces.NoSQLRows, error) {
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

func (t *MongoDB) Exec(ctx context.Context, collection string, query bson.D, opt *options.AggregateOptions) (interfaces.NoSQLRows, error) {
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

func (t *MongoDB) Insert(ctx context.Context, collection string, query bson.D, opt *options.InsertOneOptions) (interface{}, error) {
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

func (t *MongoDB) Update(ctx context.Context, collection string, query bson.D, update bson.D, opt *options.UpdateOptions) error {
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

func (t *MongoDB) Delete(ctx context.Context, collection string, query bson.D, opt *options.DeleteOptions) int64 {
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

	return write.ModifiedCount + write.InsertedCount, nil
}

func (t *MongoDB) GetDb() *mongo.Database {
	return t.db
}
