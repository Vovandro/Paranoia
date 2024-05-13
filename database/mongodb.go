package database

import (
	"Paranoia/interfaces"
	"context"
	"fmt"
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

func (t *MongoDB) Exists(ctx context.Context, query interface{}, args ...interface{}) bool {
	var opt *options.CountOptions
	var col string

	if len(args) > 0 {
		col = args[0].(string)
	} else {
		return false
	}

	if len(args) > 1 {
		if val, ok := args[1].(*options.CountOptions); ok {
			opt = val
		}
	}

	if opt == nil {
		opt = options.Count()
	}

	var limit int64 = 1
	opt.Limit = &limit

	find, err := t.db.Collection(col).CountDocuments(ctx, query, opt)

	if err != nil {
		return false
	}

	return find != 0
}

func (t *MongoDB) Count(ctx context.Context, query interface{}, args ...interface{}) int64 {

	var opt *options.CountOptions
	var col string

	if len(args) > 0 {
		col = args[0].(string)
	} else {
		return 0
	}

	if len(args) > 1 {
		if val, ok := args[1].(*options.CountOptions); ok {
			opt = val
		}
	}

	find, err := t.db.Collection(col).CountDocuments(ctx, query, opt)

	if err != nil {
		return 0
	}

	return find
}

func (t *MongoDB) FindOne(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	var opt *options.FindOneOptions
	var col string

	if len(args) > 0 {
		col = args[0].(string)
	} else {
		return fmt.Errorf("exec query collection is empty")
	}

	if len(args) > 1 {
		if val, ok := args[1].(*options.FindOneOptions); ok {
			opt = val
		}
	}

	find := t.db.Collection(col).FindOne(ctx, query, opt)

	if err := find.Err(); err != nil {
		return err
	}

	err := find.Decode(model)

	if err != nil {
		return err
	}

	return nil
}

func (t *MongoDB) Find(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	var opt *options.FindOptions
	var col string

	if len(args) > 0 {
		col = args[0].(string)
	} else {
		return fmt.Errorf("exec query collection is empty")
	}

	if len(args) > 1 {
		if val, ok := args[1].(*options.FindOptions); ok {
			opt = val
		}
	}

	find, err := t.db.Collection(col).Find(ctx, query, opt)

	if err != nil {
		return err
	}

	err = find.Decode(model)

	if err != nil {
		return err
	}

	return nil
}

func (t *MongoDB) Exec(ctx context.Context, query interface{}, model interface{}, args ...interface{}) error {
	var opt *options.AggregateOptions
	var col string

	if len(args) > 0 {
		col = args[0].(string)
	} else {
		return fmt.Errorf("exec query collection is empty")
	}

	if len(args) > 1 {
		if val, ok := args[1].(*options.AggregateOptions); ok {
			opt = val
		}
	}

	aggregate, err := t.db.Collection(col).Aggregate(ctx, query, opt)

	if err != nil {
		return err
	}

	if model != nil {
		err = aggregate.Decode(model)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *MongoDB) Update(ctx context.Context, query interface{}, args ...interface{}) error {
	var opt *options.UpdateOptions
	var col string
	var update interface{}

	if len(args) > 0 {
		col = args[0].(string)
	} else {
		return fmt.Errorf("exec query collection is empty")
	}

	if len(args) > 1 {
		update = args[1]

		if len(args) > 2 {
			if val, ok := args[2].(*options.UpdateOptions); ok {
				opt = val
			}
		}
	} else {
		return fmt.Errorf("exec update change is empty")
	}

	_, err := t.db.Collection(col).UpdateMany(ctx, query, update, opt)

	if err != nil {
		return err
	}

	return nil

}

func (t *MongoDB) Delete(ctx context.Context, query interface{}, args ...interface{}) int64 {
	var opt *options.DeleteOptions
	var col string

	if len(args) > 0 {
		col = args[0].(string)
	} else {
		return 0
	}

	if len(args) > 1 {
		if val, ok := args[1].(*options.DeleteOptions); ok {
			opt = val
		}
	}

	res, err := t.db.Collection(col).DeleteMany(ctx, query, opt)

	if err != nil {
		return 0
	}

	return res.DeletedCount
}

func (t *MongoDB) GetDb() interface{} {
	return t.db
}
