/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package kafka

import (
	"context"
	GLogger "lruCache/poc/lib/logger"
	"lruCache/poc/src/config"
	"lruCache/poc/src/constants"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

var kafkaConfig = kafka.ConfigMap{
	"bootstrap.servers":     config.KafkaBroker,
	"security.protocol":     "SASL_SSL",
	"sasl.mechanisms":       "PLAIN",
	"group.id":              config.KafkaGroupID,
	"heartbeat.interval.ms": 2000,
	"session.timeout.ms":    6000,
	"max.poll.interval.ms":  300000,
	"sasl.username":         config.KafkaClsaKey,
	"sasl.password":         config.KafkaClsaSecret,
}

type MetaDataProvider struct {
}

func (l *MetaDataProvider) AppName(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) UserID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) CorrelationID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) LoggerName(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) TraceID(ctx context.Context) string {
	return uuid.NewString()
}
func (l *MetaDataProvider) SpanID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) Thread(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientDeviceID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientPlatform(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientVersion(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) AppPackageID(ctx context.Context) string {
	return ""
}
func (l *MetaDataProvider) ClientSessionID(ctx context.Context) string {
	return ""
}

func getLogger(level GLogger.LogLevel) *GLogger.LoggerService {
	var metaDataProvider GLogger.ILoggingMetaDataProvider = &MetaDataProvider{}
	return GLogger.NewLoggerService(level, metaDataProvider)
}

func TestProduceMessage(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CorrelationID, uuid.NewString())
	logger := getLogger(GLogger.DEBUG)
	_, err := InitializeKafkaProducer(ctx, &kafkaConfig, logger)
	if err != nil {
		logger.Error(ctx, "error_while_creating_kafka_producer")
		return
	}
	producer := GetKafkaProducerFromCache("pawananjanikumar", logger)
	err = producer.ProduceKafkaMessage(ctx, config.DummyVar, map[string]interface{}{
		"name": "pawananjanikumar",
	})
	if err != nil {
		logger.Error(ctx, "error_occurred_while_producing_message_to_kafka")
		return
	}
	logger.Info(ctx, "successfully_produced_kafka_message")
}
