/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package cmd

import (
	"context"
	"fmt"
	GLogger "lruCache/poc/lib/logger"
	"lruCache/poc/src/config"
	"lruCache/poc/src/constants"
	communications2 "lruCache/poc/src/internal/communications"
	"lruCache/poc/src/modules/communications"
	"lruCache/poc/src/modules/database"
	awssecretsmanager "lruCache/poc/src/pkg/aws/secretsmanager"
	pkgKafka "lruCache/poc/src/pkg/kafka"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
)

const awsRegion = "ap-south-1"

func GetKafkaConfig(ctx context.Context, logger *GLogger.LoggerService) (kafka.ConfigMap, error) {
	var kafkaConfig = kafka.ConfigMap{
		"bootstrap.servers":     config.KafkaBroker,
		"security.protocol":     "SASL_SSL",
		"sasl.mechanisms":       "PLAIN",
		"group.id":              config.KafkaGroupID,
		"heartbeat.interval.ms": 10000,
	}
	if config.AppEnv == constants.EnvTesting {
		kafkaConfig["sasl.username"] = config.KafkaClsaKey
		kafkaConfig["sasl.password"] = config.KafkaClsaSecret
	} else {
		smClient := awssecretsmanager.GetSecretManagerClient(ctx, awsRegion, logger)
		kafkaSecrets, err := smClient.GetSecret(ctx, config.KafkaSecretArn)
		if err != nil {
			logger.Error(ctx, "error_occured_while_fetching_kafka_secrets", err)
			return nil, err
		}
		kafkaConfig["sasl.username"] = kafkaSecrets["clsa_key"]
		kafkaConfig["sasl.password"] = kafkaSecrets["clsa_secret"]
	}
	return kafkaConfig, nil
}

func EstablishDBConnection(ctx context.Context, logger *GLogger.LoggerService) (database.DBInterface, error) {
	var db database.DBInterface
	if config.AppEnv == constants.EnvTesting || config.AppEnv == constants.EnvDevelopment || config.AppEnv == constants.EnvStaging {
		logger.Info(ctx, fmt.Sprintf("Test env:%s", config.MongoDbUrlGroupSpends))
		db = database.NewDBInstance(logger, config.MongoDbUrlGroupSpends, database.OpenURI, "", database.ConnectionOptions{}, "")
	} else {
		db = database.NewDBInstance(logger, config.MongoDbUrlGroupSpends, database.AwsSecrets, config.AWSSecretMongoDB, database.ConnectionOptions{}, config.MongoDbUrlGroupSpends)
	}
	err := db.Connect(ctx)
	if err != nil {
		err = errors.Wrap(err, "error_encountered_with_db_connection")
		logger.Error(ctx, err)
		log.Fatal(err)
	}
	return db, err
}

func InitializeKafkaProducer(ctx context.Context, logger *GLogger.LoggerService) error {
	var kafkaConfig, err = GetKafkaConfig(ctx, logger)
	if err != nil {
		logger.Error(ctx, "error_while_creating_kafka_config", err)
		return err
	}
	_, err = pkgKafka.InitializeKafkaProducer(ctx, &kafkaConfig, logger)
	if err != nil {
		logger.Error(ctx, "error_occurred_while_initializing_kafka_producer", err)
		return err
	}
	return nil
}

func InitializeKafkaConsumer(ctx context.Context, consumerName string, kafkaTopics []string, uniqueMessageHandler pkgKafka.UniqueKafkaMessageHandler, logger *GLogger.LoggerService) error {
	var kafkaConfig, err = GetKafkaConfig(ctx, logger)
	if err != nil {
		logger.Error(ctx, "error_while_creating_kafka_config", err)
		return err
	}
	consumer := pkgKafka.InitializeKafkaConsumer(fmt.Sprintf("%s:%s:%s", config.AppName, constants.AppNameSuffixConsumer, consumerName), kafkaTopics, uniqueMessageHandler, logger)
	err = consumer.StartConsumer(ctx, &kafkaConfig)
	if err != nil {
		logger.Error(ctx, "error_while_running_kafka_consumer", err)
		return err
	}
	return nil
}

func InitializeCommunications() {
	communications.NewCommunications(
		config.AppName,
		&constants.CommunicationTemplates,
		&constants.CommSSEEventNameMap,
		config.SSE_Kafka_TopicName,
		communications2.BuildSSEEvent,
		nil,
		"",
		nil,
		nil,
		"",
		nil,
	)
}
