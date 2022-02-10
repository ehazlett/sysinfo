package main

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc/credentials"
)

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return hostname
}

func newExporter(ctx context.Context, endpoint string, tlsInsecure bool) (*otlptrace.Exporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithHeaders(map[string]string{}),
	}
	if tlsInsecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	} else {
		opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	client := otlptracegrpc.NewClient(opts...)
	return otlptrace.New(ctx, client)
}

func tracerProvider(endpoint, service, version string) (*tracesdk.TracerProvider, error) {
	// TODO: support tls
	exporter, err := newExporter(context.Background(), endpoint, true)
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("host", getHostname()),
			attribute.String("version", version),
		)),
	)
	return tp, nil
}

func newTraceProvider(endpoint, serviceName, version string) (*tracesdk.TracerProvider, error) {
	tp, err := tracerProvider(endpoint, serviceName, version)
	if err != nil {
		return nil, err
	}

	// register global provider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}
