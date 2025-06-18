/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package kafka

import (
	"context"
	"fmt"
	GLogger "lruCache/poc/lib/logger"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UniqueKafkaMessageHandler func(ctx context.Context, message *kafka.Message, logger *GLogger.LoggerService) error

type sKafkaConsumer struct {
	consumerName         string
	kafkaTopics          []string
	logger               *GLogger.LoggerService
	uniqueMessageHandler UniqueKafkaMessageHandler
}

type IKafkaConsumer interface {
	StartConsumer(ctx context.Context, kafkaConfig *kafka.ConfigMap) error
}

func InitializeKafkaConsumer(consumerName string, kafkaTopics []string, singleMessageHandler UniqueKafkaMessageHandler, logger *GLogger.LoggerService) IKafkaConsumer {
	return &sKafkaConsumer{
		consumerName:         consumerName,
		kafkaTopics:          kafkaTopics,
		logger:               logger,
		uniqueMessageHandler: singleMessageHandler,
	}
}

func (kc *sKafkaConsumer) StartConsumer(ctx context.Context, kafkaConfig *kafka.ConfigMap) error {
	ctx = context.WithValue(ctx, AppName, kc.consumerName)
	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		kc.logger.Error(ctx, fmt.Sprintf("error_while_creating_kafka_consumer"), err)
		return err
	}
	kafkaConfig.SetKey("sasl.username", "**")
	kafkaConfig.SetKey("sasl.password", "**")
	kc.logger.Info(ctx, "initiating_kafka_consumer", kafkaConfig)
	err = consumer.SubscribeTopics(kc.kafkaTopics, nil)
	if err != nil {
		kc.logger.Error(ctx, "error_while_connecting_to_kafka_consumer", err)
		return err
	}
	kc.logger.Info(ctx, fmt.Sprintf("consumer_started_successfully"), consumer)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	run := true
	tracer := otel.Tracer(fmt.Sprintf("tracer:%s", kc.consumerName))
	for run {
		select {
		case sig := <-sigchan:
			// Handle OS signal to shut down consumer gracefully
			kc.logger.Info(ctx, "os_interrupt_received", map[string]interface{}{"sig": sig})
			kc.logger.Info(ctx, "Caught  os.Interrupt %v: terminating\n run=false", nil)
			kc.logger.Info(ctx, fmt.Sprintf("terminating_consumer_listener"), nil)
			run = false
			break
		default:
			// Poll Kafka for new messages
			ev := consumer.Poll(10000)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				kc.processMessage(tracer, e)
				break
			case kafka.Error:
				// Handle Kafka errors
				kc.logger.Error(ctx, "kafka_error", e)
				run = false
				break
			case kafka.PartitionEOF:
				kc.logger.Info(ctx, "reached_end_of_partition", e)
				break
			}
		}
	}
	kc.logger.Info(ctx, "closing_consumer")
	err = consumer.Close()
	if err != nil {
		kc.logger.Error(ctx, "error_occurred_while_closing_consumer", err)
		return err
	}
	return nil
}

func (kc *sKafkaConsumer) processMessage(tracer trace.Tracer, kafkaMessage *kafka.Message) {
	kafkaMessageContext := context.Background()
	kafkaMessageContext = context.WithValue(kafkaMessageContext, AppName, kc.consumerName)
	msgSpanCtx, msgSpan := tracer.Start(kafkaMessageContext, "processingKafkaMessage")
	msgSpan.SetAttributes(
		attribute.String("kafka.topic", *kafkaMessage.TopicPartition.Topic),
		attribute.Int64("kafka.partition", int64(kafkaMessage.TopicPartition.Partition)),
		attribute.Int64("kafka.offset", int64(kafkaMessage.TopicPartition.Offset)),
		attribute.String("kafka.key", string(kafkaMessage.Key)),
	)
	defer msgSpan.End()
	kc.logger.Info(msgSpanCtx, "kafka_message_received", map[string]interface{}{"key": string(kafkaMessage.Key)})
	for _, header := range kafkaMessage.Headers {
		if header.Key == XCorrelationID {
			kc.logger.Info(msgSpanCtx, "kafka_message_header", map[string]interface{}{"headerKey": header.Key, "headerValue": string(header.Value)})
			correlationID := string(header.Value)
			msgSpanCtx = context.WithValue(msgSpanCtx, CorrelationID, correlationID)
			msgSpan.SetAttributes(
				attribute.String(XCorrelationID, correlationID),
			)
		}
	}
	kc.logger.Info(msgSpanCtx, "received_message", string(kafkaMessage.Value))
	err := kc.uniqueMessageHandler(msgSpanCtx, kafkaMessage, kc.logger)
	if err != nil {
		kc.logger.Error(msgSpanCtx, "error_occurred_while_processing_kafka_message", err.Error())
	}
	return
}
