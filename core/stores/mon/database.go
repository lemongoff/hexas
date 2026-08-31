package mon

import (
	"context"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/v2/mongo/options"
)

type (
	Database struct {
		name string

		client      *mongo.Client
		database    *mongo.Database
		brk         breaker.Breaker
		opts        []Option
		mutex       sync.RWMutex
		collections map[string]Collection
	}
)

func MustNewDatabase(uri, db string, opts ...Option) *Database {
	database, err := NewDatabase(uri, db, opts...)
	logx.Must(err)
	return database
}

func NewDatabase(uri, db string, opts ...Option) (*Database, error) {
	clientOptions := mongoOptions.Client().ApplyURI(uri)
	for _, opt := range append([]Option{defaultTimeoutOption()}, opts...) {
		opt(clientOptions)
	}
	cli, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}
	if err := cli.Ping(context.Background(), nil); err != nil {
		_ = cli.Disconnect(context.Background())
		return nil, err
	}

	name := strings.Join([]string{uri, db}, "/")
	brk := breaker.GetBreaker(uri)
	return newDatabase(name, cli, cli.Database(db), brk, opts...), nil
}

func newDatabase(name string, client *mongo.Client, database *mongo.Database, brk breaker.Breaker,
	opts ...Option) *Database {
	return &Database{
		name:        name,
		client:      client,
		database:    database,
		brk:         brk,
		opts:        opts,
		collections: map[string]Collection{},
	}
}

func (db *Database) Close() error {
	return db.client.Disconnect(context.Background())
}

func (db *Database) Collection(collection string) Collection {
	db.mutex.RLock()
	coll := db.collections[collection]
	db.mutex.RUnlock()
	if coll != nil {
		return coll
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()
	if coll = db.collections[collection]; coll != nil {
		return coll
	}
	coll = newCollection(db.database.Collection(collection), db.brk)
	db.collections[collection] = coll
	return coll
}

func (db *Database) CreateIndex(ctx context.Context, collection string, unique bool, fields ...string) error {
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: 1})
	}
	_, err := db.Collection(collection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    keys,
		Options: mongoOptions.Index().SetUnique(unique),
	})
	return err
}
