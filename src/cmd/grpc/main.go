package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"lruCache/poc/src/cmd"
	"lruCache/poc/src/config"
	"lruCache/poc/src/constants"
	appgrpc "lruCache/poc/src/internal/app/grpc"
	"lruCache/poc/src/internal/helpers"
	"lruCache/poc/src/modules/database"
	"lruCache/poc/src/pkg/otel"
	"log"
	"net"
)

func main() {
	l := helpers.InitLogger()
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.AppName, fmt.Sprintf("%s:%s", config.AppName, constants.AppNameSuffixHTTPServer))
	cleanup := otel.InitTracer(ctx, config.Secure, config.CollectorURL, config.ServiceName, l)
	defer func() {
		err := cleanup(ctx)
		err = errors.Wrap(err, "open-telemetry tracer could not be closed successfully")
		l.Error(ctx, err)
	}()
	//err := cmd.InitializeKafkaProducer(ctx, l)
	//if err != nil {
	//	err = errors.Wrap(err, "error_occurred_while_initializing_kafka_producer")
	//	l.Error(ctx, err)
	//	log.Fatal("could_not_start_gin_server " + err.Error())
	//}
	db, err := cmd.EstablishDBConnection(ctx, l)
	defer func(db database.DBInterface) {
		_ = db.Disconnect(ctx)
	}(db)

	cmd.InitializeCommunications()

	// Create a new gRPC server
	grpcServer := appgrpc.SetupServer(ctx, l)

	// Create a TCP listener at port 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GrpcPort))

	// Start serving requests at port 8080
	if err := grpcServer.Serve(listener); err != nil {
		err = errors.Wrap(err, "error_occurred_while_starting_grpc_server")
		l.Error(ctx, err)
		log.Fatal("could_not_start_grpc_server " + err.Error())
	}

	if err != nil {
		log.Fatalf("Error creating the server %v", err)
	}

	l.Info(ctx, "Server started at port :8080")

}
