/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024
*/

package mongo

import (
	"context"
	GLogger "lruCache/poc/lib/logger"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mongotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go.mongodb.org/mongo-driver/mongo"
)

type Mongo interface {
	Connect(ctx context.Context) (*mongo.Client, error)
	GetMongoClient() *mongo.Client
}

type db struct {
	ConnectionURI   string
	MaxPoolSize     uint64
	MinPoolSize     uint64
	MaxConnIdleTime *time.Duration
	Log             *GLogger.LoggerService
	EnableDDTrace   bool
	BSONOptions     *options.BSONOptions
	dbClient        *mongo.Client
}

func NewDB(connectionURI string, log *GLogger.LoggerService, maxPoolSize uint64, minPoolSize uint64, maxConnIdleTime *time.Duration, enableDDTrace bool, BSONOptions *options.BSONOptions) Mongo {
	return &db{ConnectionURI: connectionURI, MaxPoolSize: maxPoolSize, MinPoolSize: minPoolSize, MaxConnIdleTime: maxConnIdleTime, Log: log, EnableDDTrace: enableDDTrace, BSONOptions: BSONOptions}
}

func (m *db) Connect(ctx context.Context) (*mongo.Client, error) {
	if m.EnableDDTrace == true {
		mongotrace.NewMonitor(mongotrace.WithAnalytics(true))
	}
	clientOptions := options.Client().ApplyURI(m.ConnectionURI)
	if m.MaxPoolSize > 1 {
		clientOptions.SetMaxPoolSize(m.MaxPoolSize)
	}
	if m.MinPoolSize > 1 {
		clientOptions.SetMinPoolSize(m.MinPoolSize)
	} else {
		clientOptions.SetMinPoolSize(1)
	}
	if m.MaxConnIdleTime != nil {
		clientOptions.SetMaxConnIdleTime(*m.MaxConnIdleTime)
		//	The default is 0, meaning a connection can remain unused indefinitely
	}
	if m.BSONOptions != nil {
		clientOptions.SetBSONOptions(m.BSONOptions)
	} else {
		clientOptions.SetBSONOptions(&options.BSONOptions{
			NilMapAsEmpty: true,
			// NilMapAsEmpty causes the driver to marshal nil Go maps as empty BSON
			// documents instead of BSON null.
			// Empty BSON documents take up slightly more space than BSON null, but
			// preserve the ability to use document update operations like "$set" that
			// do not work on BSON null.
			NilSliceAsEmpty: true,
			// NilSliceAsEmpty causes the driver to marshal nil Go slices as empty BSON
			// arrays instead of BSON null.
			// Empty BSON arrays take up slightly more space than BSON null, but
			// preserve the ability to use array update operations like "$push" or
			// "$addToSet" that do not work on BSON null.
			NilByteSliceAsEmpty: true,
			OmitZeroStruct:      true,
			DefaultDocumentD:    true,
			// DefaultDocumentD causes the driver to always unmarshal documents into the
		})
	}
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		m.Log.Error(ctx, "db_connection_failed", err)
		return nil, err
	}
	m.dbClient = client
	return m.dbClient, nil
}

func (m *db) GetMongoClient() *mongo.Client {
	return m.dbClient
}
