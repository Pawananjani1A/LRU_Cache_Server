/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package database

import (
	"context"
	"fmt"
	GLogger "lruCache/poc/lib/logger"
	awssecretsmanager "lruCache/poc/src/pkg/aws/secretsmanager"
	mongolib "lruCache/poc/src/pkg/mongo"
	"reflect"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type databaseMongo struct {
	log               *GLogger.LoggerService
	ctx               context.Context
	connectionURI     string
	authMethod        AuthMethod
	awsSecretARN      string
	awsRegion         string
	connectionOptions ConnectionOptions
	dbClient          *mongo.Client
	defaultDB         string
}

func (m *databaseMongo) makeConnectURIWithAwsSecrets(ctx context.Context, awsRegion string) (string, error) {
	smClient := awssecretsmanager.GetSecretManagerClient(ctx, awsRegion, m.log)
	mongoSecrets, err := smClient.GetSecret(ctx, m.awsSecretARN)
	if err != nil {
		err := errors.Wrap(err, "error_occured_while_trying_to_fetch_mongo_secrets")
		m.log.Error(ctx, err)
		return "", err
	}
	return fmt.Sprintf(m.connectionURI, mongoSecrets["username"], mongoSecrets["password"], m.defaultDB), nil
}

func (m *databaseMongo) Connect(ctx context.Context) error {
	var connectionURI string
	var err error
	if m.authMethod == OpenURI || m.authMethod == AwsIAM {
		connectionURI = m.connectionURI
	} else if m.authMethod == AwsSecrets {
		connectionURI, err = m.makeConnectURIWithAwsSecrets(ctx, m.awsRegion)
	}
	db := mongolib.NewDB(connectionURI, m.log, m.connectionOptions.maxPoolSize, m.connectionOptions.minPoolSize, m.connectionOptions.maxConnIdleTime, m.connectionOptions.enableDDTrace, m.connectionOptions.bsonOptions)
	client, err := db.Connect(ctx)
	if err != nil {
		err = errors.Wrap(err, "error_while_establishing_connection_with_mongo")
		m.log.Error(ctx, err)
		return err
	}
	m.dbClient = client
	setDBClient(client)
	m.log.Info(ctx, "successfully_established_db_connection")
	return nil
}

// Disconnect Close terminates the connection to the MongoDB server.
func (m *databaseMongo) Disconnect(ctx context.Context) error {
	if m.dbClient != nil {
		// Close the database client connection.
		err := m.dbClient.Disconnect(ctx)
		if err != nil {
			m.log.Error(ctx, fmt.Errorf("error_while_closing_mongo_connection: %w", err))
			return err
		}
		m.log.Info(ctx, "mongo_connection_closed_successfully")
	}
	return nil
}

func (m *databaseMongo) Find(ctx context.Context, dbName string, collectionName string, filter bson.M, sortOptions *options.FindOptions, results any) error {
	resultsVal := reflect.ValueOf(results)
	if resultsVal.Kind() != reflect.Ptr || resultsVal.Elem().Kind() != reflect.Slice {
		m.log.Error(ctx, "resultsVal must be a pointer to a slice")
		return errors.New(DBErrInvalidResultsParamReceived)
	}
	collection := m.dbClient.Database(dbName).Collection(collectionName)
	m.log.Debug(ctx, "db.find().query", map[string]interface{}{
		"dbName":         dbName,
		"collectionName": collectionName,
		"filter":         filter,
		"sortOptions":    sortOptions,
	})
	cursor, err := collection.Find(ctx, filter, sortOptions)
	if err != nil {
		m.log.Error(ctx, "error_occurred_in_find_query", err)
		return errors.New(DBErrFindQueryError)
	}
	defer cursor.Close(ctx)
	// Iterate through the cursor and decode results
	sliceVal := resultsVal.Elem()
	elemType := sliceVal.Type().Elem()
	for cursor.Next(ctx) {
		elemPtr := reflect.New(elemType)
		if err = cursor.Decode(elemPtr.Interface()); err != nil {
			return fmt.Errorf("error decoding cursor result: %w", err)
		}
		sliceVal.Set(reflect.Append(sliceVal, elemPtr.Elem()))
	}
	m.log.Debug(ctx, "db.find().result", map[string]interface{}{
		"dbName":         dbName,
		"collectionName": collectionName,
		"filter":         filter,
		"sortOptions":    sortOptions,
		"results":        results,
	})
	return nil
}

func (m *databaseMongo) Insert(ctx context.Context, dbName string, collectionName string, data any) error {
	collection := m.dbClient.Database(dbName).Collection(collectionName)
	m.log.Debug(ctx, "db.insert().query", map[string]interface{}{
		"dbName":         dbName,
		"collectionName": collectionName,
		"dataToInsert":   data,
	})
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		m.log.Error(ctx, "db_insert_failed", err)
		return err
	}
	m.log.Debug(ctx, "db_insert_success", res)
	return nil
}

func (m *databaseMongo) UpdateOne(ctx context.Context, dbName string, collectionName string, filter bson.M, dataToUpdate bson.M, updateOptions *options.UpdateOptions, throwNoMatchedDocumentErr bool) (*mongo.UpdateResult, error) {
	collection := m.dbClient.Database(dbName).Collection(collectionName)
	m.log.Debug(ctx, "db.updateOne().query", map[string]interface{}{
		"dbName":         dbName,
		"collectionName": collectionName,
		"filter":         filter,
		"updateOptions":  updateOptions,
		"dataToUpdate":   dataToUpdate,
	})
	result, err := collection.UpdateOne(ctx, filter, dataToUpdate, updateOptions)
	if err != nil {
		m.log.Error(ctx, "db_update_failed", err)
		return nil, errors.New(DBErrUpdateFailed)
	}
	m.log.Debug(ctx, "db.updateOne().query", map[string]interface{}{
		"dbName":         dbName,
		"collectionName": collectionName,
		"filter":         filter,
		"updateOptions":  updateOptions,
		"dataToUpdate":   dataToUpdate,
		"results":        result,
	})
	if result.MatchedCount == 0 && throwNoMatchedDocumentErr == true {
		m.log.Debug(ctx, "db_update_no_match", result)
		return result, errors.New(DBErrNoDocumentsMatched)
	}
	m.log.Debug(ctx, "db_update_success", result)
	return result, nil
}

func (m *databaseMongo) UpdateMany(ctx context.Context, dbName string, collectionName string, filter bson.M, dataToUpdate bson.M, updateOptions *options.UpdateOptions, throwNoMatchedDocumentErr bool) (*mongo.UpdateResult, error) {
	collection := m.dbClient.Database(dbName).Collection(collectionName)
	m.log.Debug(ctx, "db.updateMany().query", map[string]interface{}{
		"dbName":         dbName,
		"collectionName": collectionName,
		"filter":         filter,
		"updateOptions":  updateOptions,
		"dataToUpdate":   dataToUpdate,
	})
	result, err := collection.UpdateMany(ctx, filter, dataToUpdate, updateOptions)
	if err != nil {
		m.log.Error(ctx, "db_update_failed", err)
		return nil, errors.New(DBErrUpdateFailed)
	}
	m.log.Debug(ctx, "db.updateMany().result", map[string]interface{}{
		"dbName":         dbName,
		"collectionName": collectionName,
		"filter":         filter,
		"updateOptions":  updateOptions,
		"dataToUpdate":   dataToUpdate,
		"results":        result,
	})
	if result.MatchedCount == 0 && throwNoMatchedDocumentErr == true {
		m.log.Debug(ctx, "db_update_no_match", result)
		return result, errors.New(DBErrNoDocumentsMatched)
	}
	m.log.Debug(ctx, "db_update_success", result)
	return result, nil
}
