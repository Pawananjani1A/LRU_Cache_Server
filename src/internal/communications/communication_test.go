/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package communications

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	GLogger "lruCache/poc/lib/logger"
	"lruCache/poc/src/config"
	"lruCache/poc/src/constants"
	awssecretsmanager "lruCache/poc/src/pkg/aws/secretsmanager"
)

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
		smClient := awssecretsmanager.GetSecretManagerClient(ctx, "ap-south-1", logger)
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
