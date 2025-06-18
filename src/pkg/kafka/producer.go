/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package kafka

import (
	"context"
	"encoding/json"
	GLogger "lruCache/poc/lib/logger"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

/*
	Limitation: 1 kafkaConfig :: Many kafkaTopics
*/

var kafkaMutex = sync.RWMutex{}

var (
	kafkaProducerCache *kafka.Producer
	kafkaConfigCache   *kafka.ConfigMap
)

func getKafkaProducer() *kafka.Producer {
	kafkaMutex.RLock()
	defer kafkaMutex.RUnlock()
	return kafkaProducerCache
}

func setKafkaProducer(kafkaProducer *kafka.Producer) {
	kafkaMutex.Lock()
	defer kafkaMutex.Unlock()
	kafkaProducerCache = kafkaProducer
}

func getKafkaConfig() *kafka.ConfigMap {
	kafkaMutex.RLock()
	defer kafkaMutex.RUnlock()
	return kafkaConfigCache
}

func setKafkaConfig(kafkaConfig *kafka.ConfigMap) {
	kafkaMutex.Lock()
	defer kafkaMutex.Unlock()
	kafkaConfigCache = kafkaConfig
}

type sKafkaProducer struct {
	kafkaConfig   *kafka.ConfigMap
	kafkaProducer *kafka.Producer
	producedBy    string
	logger        *GLogger.LoggerService
}

type IKafkaProducer interface {
	ProduceKafkaMessage(ctx context.Context, topicName string, message interface{}) error
}

func InitializeKafkaProducer(ctx context.Context, kafkaConfig *kafka.ConfigMap, logger *GLogger.LoggerService) (IKafkaProducer, error) {
	kp := sKafkaProducer{
		kafkaConfig: kafkaConfig,
		logger:      logger,
	}
	setKafkaConfig(kafkaConfig)
	kp.logger.Info(ctx, "creating_new_kafka_producer_client")
	configMap := kafkaConfigCache
	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		setKafkaProducer(nil)
		kp.logger.Error(ctx, "error_while_creating_kafka_producer", err)
		return nil, err
	}
	setKafkaProducer(producer)
	kp.kafkaProducer = producer
	kp.logger.Info(ctx, "successfully_created_kafka_producer_client")
	if err != nil {
		kp.logger.Error(ctx, "error_occurred_while_creating_kafka_producer", err)
		return nil, err
	}
	return &kp, nil
}

func GetKafkaProducerFromCache(producedBy string, logger *GLogger.LoggerService) IKafkaProducer {
	kafkaConfig := getKafkaConfig()
	kafkaProducer := getKafkaProducer()
	return &sKafkaProducer{
		kafkaConfig:   kafkaConfig,
		kafkaProducer: kafkaProducer,
		logger:        logger,
		producedBy:    producedBy,
	}
}

func (kp *sKafkaProducer) ProduceKafkaMessage(ctx context.Context, topicName string, message interface{}) error {
	payload, _ := json.Marshal(message)
	deliveryChan := make(chan kafka.Event)
	correlationID := GetDefaultValueFromContext(ctx, CorrelationID)
	requestID := GetDefaultValueFromContext(ctx, RequestID)
	kp.logger.Info(ctx, "producing_message_to_kafka", map[string]interface{}{
		"payload":     string(payload),
		CorrelationID: correlationID,
		RequestID:     requestID,
	})
	err := kp.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
		Value:          payload,
		Headers: []kafka.Header{
			{Key: string(XCorrelationID), Value: []byte(correlationID)},
			{Key: "x-event-source", Value: []byte(kp.producedBy)},
			{Key: string(XRequestID), Value: []byte(requestID)},
			{Key: "x-message-id", Value: []byte(uuid.NewString())},
		},
	}, deliveryChan)
	if err != nil {
		kp.logger.Error(ctx, "failed_to_produce_messages: ", err)
		kp.logger.Error(ctx, "error_occurred_while_producing_kafka_message")
		//	trigger error notification for kafka event not produced
		return err
	} else {
		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			kp.logger.Error(ctx, "kafka_message_delivery_failed: %v", m.TopicPartition.Error)
			return err
		} else {
			kp.logger.Info(ctx, "kafka_message_successfully_produced: (offset: %v)", m.TopicPartition.Offset)
		}
	}
	return nil
}
