/*
   Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
   CreatedAt: 28 Mar 2024*/

package main

import (
	"context"
	"fmt"
	"log"
	"lruCache/poc/src/cmd"
	"lruCache/poc/src/config"
	"lruCache/poc/src/constants"
	appgin "lruCache/poc/src/internal/app/gin"
	"lruCache/poc/src/internal/helpers"
	"lruCache/poc/src/pkg/otel"

	"github.com/pkg/errors"
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

	cmd.InitializeCommunications()
	server := appgin.SetupServer(ctx, l)
	err := server.Run(config.Port)
	if err != nil {
		err = errors.Wrap(err, "error_occurred_while_starting_gin_server")
		l.Error(ctx, err)
		log.Fatal("could_not_start_gin_server " + err.Error())
	}

	//// Graceful shutdown handling
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, os.Interrupt)
	//
	//// Block until a signal is received
	//<-quit
	//fmt.Println("Shutting down server...")
	//
	//// Create a context with a timeout
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//
	//// Attempt to gracefully shutdown the server
	//if err := server.Shutdo(ctx); err != nil {
	//	fmt.Printf("Error during server shutdown: %v\n", err)
	//} else {
	//	fmt.Println("Server gracefully stopped")
	//}
}
