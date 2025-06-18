/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package app_grpc

import (
	"context"
	proto "lruCache/poc/gen"
	GLogger "lruCache/poc/lib/logger"
	chathandler "lruCache/poc/src/internal/handlers/grpc/chats"
	"net"

	"google.golang.org/grpc"
)

func SetupServer(ctx context.Context, log *GLogger.LoggerService) *grpc.Server {
	log.Info(ctx, "setting_up_routes...")

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create a new connection pool
	var conn []*chathandler.Connection

	pool := &chathandler.Pool{
		Connection: conn,
	}

	// Register the pool with the gRPC server
	proto.RegisterBroadcastServer(grpcServer, pool)

	// Create a TCP listener at port 8080
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Error(ctx, "Error creating the server %v", err)
	}

	log.Info(ctx, "Server started at port :8080")

	// Start serving requests at port 8080
	if err := grpcServer.Serve(listener); err != nil {
		log.Error(ctx, "Error creating the server %v", err)
	}

	return grpcServer
}
