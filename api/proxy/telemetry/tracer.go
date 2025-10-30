package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// InitTracer initializes the OpenTelemetry tracer provider with OTLP exporter.
// Returns a shutdown function that should be called when the application terminates.
func InitTracer(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	if !cfg.Enabled {
		// Return a no-op shutdown function if tracing is not enabled
		return func(context.Context) error { return nil }, nil
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Set up OTLP exporter (using HTTP)
	var opts []otlptracehttp.Option
	opts = append(opts, otlptracehttp.WithEndpoint(cfg.ExporterEndpoint))

	if cfg.ExporterInsecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Create tracer provider with sampling configuration
	// TraceIDRatioBased handles edge cases: returns AlwaysSample for rate >= 1.0, NeverSample for rate <= 0.0
	sampler := sdktrace.TraceIDRatioBased(cfg.TraceSampleRate)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator to propagate trace context (W3C Trace Context)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Return shutdown function
	return tp.Shutdown, nil
}
