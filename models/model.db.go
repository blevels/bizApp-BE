package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"backend/config"
)

type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

type CollectionHelper interface {
	Find(context.Context, interface{}) (*mongo.Cursor, error)
	FindOne(context.Context, interface{}) SingleResultHelper
	InsertOne(context.Context, interface{}) (interface{}, error)
	DeleteOne(ctx context.Context, filter interface{}) (int64, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (interface{}, error)
}

type SingleResultHelper interface {
	Decode(v interface{}) error
}

type ClientHelper interface {
	Database(string) DatabaseHelper
	Connect() error
	StartSession() (mongo.Session, error)
	Ping() error
}

type mongoClient struct {
	cl *mongo.Client
}

type mongoDatabase struct {
	db *mongo.Database
}

type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoSession struct {
	mongo.Session
}


func CreateDBHelperService() DatabaseHelper {
	c := config.CreateConfigService()

	cl, err := NewClient(c)
	if err != nil {
		log.Fatal(err)
	}

	err = cl.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = cl.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return NewDatabase(c, cl)
}

func NewClient(cnf *config.ConfigService) (ClientHelper, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s", cnf.Config.Database.User, cnf.Config.Database.Password, cnf.Config.Database.Host, cnf.Config.Database.Port)

	c, err := mongo.NewClient(options.Client().SetAuth(
		options.Credential{
			Username:   cnf.Config.Database.User,
			Password:   cnf.Config.Database.Password,
			AuthSource: "admin",
		}).ApplyURI(url))

	return &mongoClient{cl: c}, err
}

func NewDatabase(cnf *config.ConfigService, client ClientHelper) DatabaseHelper {
	return client.Database(cnf.Config.Database.Database)
}

func (mc *mongoClient) Database(dbName string) DatabaseHelper {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.cl.StartSession()
	return &mongoSession{session}, err
}

func (mc *mongoClient) Connect() error {
	// mongo client does not use context on connect method. There is a ticket
	// with a request to deprecate this functionality and another one with
	// explanation why it could be useful in synchronous requests.
	// https://jira.mongodb.org/browse/GODRIVER-1031
	// https://jira.mongodb.org/browse/GODRIVER-979
	return mc.cl.Connect(nil)
}

func (mc *mongoClient) Ping() error {
	return mc.cl.Ping(context.TODO(), nil)
}

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	cur, err := mc.coll.Find(ctx, filter)
	return cur, err
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResultHelper {
	singleResult := mc.coll.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	return id.InsertedID, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.coll.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (mc *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (interface{}, error) {
	res, err := mc.coll.UpdateOne(ctx, filter, update)
	return res, err
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}