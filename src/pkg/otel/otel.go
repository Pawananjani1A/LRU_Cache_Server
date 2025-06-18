/*
	Author: Pawananjani Kumar (pawananjanimth1@gmail.com)
	CreatedAt: 28 Mar 2024*/

package otel

import (
	"context"
	GLogger "lruCache/poc/lib/logger"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(ctx context.Context, secure string, collectorURL string, serviceName string, logger *GLogger.LoggerService) func(context.Context) error {
	logger.Info(ctx, "initiating_otel_service")
	secureOption := otlptracegrpc.WithInsecure()
	if len(secure) > 0 {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	}
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
		),
	)
	if err != nil {
		logger.Error(ctx, "error_occured_while_initiating_otel_exporter", err)
		log.Fatal(err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("modules.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.Error(ctx, "Could not set resources for starting open-telemetry tracing: ", err)
		log.Fatal(err)
	}
	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	logger.Info(ctx, "successfully_started_otel_service")
	return exporter.Shutdown
}
