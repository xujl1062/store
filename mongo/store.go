package mongo

import (
	"context"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/xujl1062/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type db struct {
	client *mongo.Client
	logger *log.Logger
}

func NewStore(uri string, logger *log.Logger) (store.Store, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return db{
		client: client,
		logger: logger,
	}, nil
}

func (db db) Get(ctx context.Context, key string, criteria map[string]interface{}, out interface{}) (err error) {
	c, err := parsingKey(db.client, key)
	if err != nil {
		return InvalidKeyError
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidPtrError{reflect.TypeOf(out)}
	}
	return c.FindOne(ctx, criteria).Decode(out)
}

func (db db) Save(ctx context.Context, key string, entity interface{}) error {
	c, err := parsingKey(db.client, key)
	if err != nil {
		return err
	}
	childCtx, _ := context.WithTimeout(ctx, 5*time.Second)
	_, err = c.InsertOne(childCtx, entity)
	if err != nil {
		return errors.Wrap(err, "Insert department error")
	}
	db.logger.Printf("Success insert %v \n", entity)
	return nil
}

func (db db) List(ctx context.Context, key string, criteria map[string]interface{}, target interface{}) ([]interface{}, error) {
	c, err := parsingKey(db.client, key)
	if err != nil {
		return nil, err
	}
	rv := reflect.ValueOf(target)
	if reflect.Ptr != rv.Kind() || rv.IsNil() {
		return nil, errors.New("Param target must be ptr")
	}
	childCtx, _ := context.WithTimeout(ctx, 5*time.Second)
	cur, err := c.Find(childCtx, criteria)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	list := make([]interface{}, 0)
	for cur.Next(context.Background()) {
		v := reflect.New(reflect.TypeOf(target).Elem()).Interface()
		err := cur.Decode(v)
		if err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}

func (db db) Update(ctx context.Context, key string, criteria map[string]interface{}, entity interface{}) error {
	c, err := parsingKey(db.client, key)
	if err != nil {
		return err
	}
	result, err := c.UpdateOne(ctx, criteria, entity)
	if err != nil {
		db.logger.Println("Error ", "update doc ", err.Error())
		return err
	}
	db.logger.Println("Success ", "update count: ", result.ModifiedCount)
	return nil
}

func (db db) FindAndUpdate(ctx context.Context, key string, criteria map[string]interface{}, update interface{}) error {
	c, err := parsingKey(db.client, key)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(update)
	if reflect.Ptr != rv.Kind() || rv.IsNil() {
		return errors.New("Param target must be ptr")
	}

	result := c.FindOneAndUpdate(ctx, criteria, update)
	if result.Err() != nil {
		return result.Err()
	}

	if err := result.Decode(update); err != nil {
		return err
	}

	return nil
}

func parsingKey(client *mongo.Client, key string) (*mongo.Collection, error) {
	paths := strings.Split(key, "/")
	if len(paths) != 2 {
		return nil, errors.New("Expected key format is \"database/collection\" ")
	}
	db, col := paths[0], paths[1]
	c := client.Database(db).Collection(col)
	return c, nil
}
