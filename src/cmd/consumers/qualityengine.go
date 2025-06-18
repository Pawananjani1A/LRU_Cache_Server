/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package main

import (
	"context"
	"fmt"
	"lruCache/poc/src/cmd"
	"lruCache/poc/src/config"
	"lruCache/poc/src/constants"
	"lruCache/poc/src/internal/helpers"
	"lruCache/poc/src/modules/database"
	"lruCache/poc/src/pkg/otel"
	"log"

	"github.com/pkg/errors"
)

func main() {
	l := helpers.InitLogger()
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.AppName, fmt.Sprintf("%s:%s", config.AppName, constants.AppNameSuffixConsumer))
	cleanup := otel.InitTracer(ctx, config.Secure, config.CollectorURL, config.ServiceName, l)
	defer func() {
		err := cleanup(ctx)
		err = errors.Wrap(err, "open-telemetry tracer could not be closed successfully")
		l.Error(ctx, err)
	}()
	db, err := cmd.EstablishDBConnection(ctx, l)
	defer func(db database.DBInterface) {
		_ = db.Disconnect(ctx)
	}(db)
	err = cmd.InitializeKafkaProducer(ctx, l)
	if err != nil {
		err = errors.Wrap(err, "error_occurred_while_initializing_kafka_producer")
		l.Error(ctx, err)
		log.Fatal("could_not_start_gin_server " + err.Error())
	}
	cmd.InitializeCommunications()

}
