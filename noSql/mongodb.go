package noSql

import (
	"context"
	"fmt"
	"gitlab.com/devpro_studio/Paranoia/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Name     string
	Database string
	Options  *options.ClientOptions
	app      interfaces.IService
	client   *mongo.Client
	db       *mongo.Database
}

func (t *MongoDB) Init(app interfaces.IService) error {
	t.app = app
	var err error

	t.client, err = mongo.Connect(context.TODO(), t.Options)

	if err != nil {
		return err
	}

	t.db = t.client.Database(t.Database)

	return nil
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
