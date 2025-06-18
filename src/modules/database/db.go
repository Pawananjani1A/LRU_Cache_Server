/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package database

import (
	"context"
	GLogger "lruCache/poc/lib/logger"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbInstance DBInterface
var dbClient *mongo.Client
var dbMutex = sync.RWMutex{}

type DBInterface interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Find(ctx context.Context, dbName string, collectionName string, filter bson.M, sortOptions *options.FindOptions, results any) error
	Insert(ctx context.Context, dbName string, collectionName string, data any) error
	UpdateOne(ctx context.Context, dbName string, collectionName string, filter bson.M, update bson.M, updateOptions *options.UpdateOptions, throwNoMatchedDocumentErr bool) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, dbName string, collectionName string, filter bson.M, dataToUpdate bson.M, updateOptions *options.UpdateOptions, throwNoMatchedDocumentErr bool) (*mongo.UpdateResult, error)
}

func GetDBInstance() DBInterface {
	return dbInstance
}

func setDBInstance(instance DBInterface) {
	dbInstance = instance
}

func GetDBClient() *mongo.Client {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	return dbClient
}

func setDBClient(client *mongo.Client) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	dbClient = client
}

type ConnectionOptions struct {
	maxPoolSize     uint64
	minPoolSize     uint64
	maxConnIdleTime *time.Duration
	enableDDTrace   bool
	bsonOptions     *options.BSONOptions
}

func NewDBInstance(
	log *GLogger.LoggerService,
	connectionURI string,
	authMethod AuthMethod,
	awsSecretARN string,
	connectionOptions ConnectionOptions,
	defaultDbName string) DBInterface {
	var instance DBInterface = &databaseMongo{
		log:               log,
		connectionURI:     connectionURI,
		authMethod:        authMethod,
		awsSecretARN:      awsSecretARN,
		connectionOptions: connectionOptions,
		dbClient:          nil,
		defaultDB:         defaultDbName,
	}
	setDBInstance(instance)
	return instance
}
