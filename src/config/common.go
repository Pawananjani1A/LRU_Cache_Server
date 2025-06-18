/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package config

import (
	"os"
)

var (
	Port          = os.Getenv("app_port")
	GrpcPort      = os.Getenv("grpc_port")
	AppName       = os.Getenv("app_name")
	AppVersion    = os.Getenv("app_version")
	Application   = os.Getenv("application")
	DDServiceName = os.Getenv("dd_service_name")
	ServiceName   = os.Getenv("service_name")
	AwsRegion     = os.Getenv("aws_region")
	CollectorURL  = os.Getenv("otel_exporter_url")
	Secure        = os.Getenv("secure_otel")
)

var (
	AWSSecretMongoDB = os.Getenv("awssecrets_mongodb")
)
var (
	KafkaBroker           = os.Getenv("kafka_broker")
	KafkaGroupID          = os.Getenv("kafka_groupid")
	KafkaSecretArn        = os.Getenv("kafka_secretarn")
	KafkaClsaKey          = os.Getenv("kafka_clsakey")
	KafkaClsaSecret       = os.Getenv("kafka_clsasecret")
	MongoDbUrlGroupSpends = os.Getenv("mongodbUrl_groupSpends")
	SSE_Kafka_TopicName   = os.Getenv("gs_sse_topic")
	DummyVar              = os.Getenv("dummyVar")
)
