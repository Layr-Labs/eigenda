package tracing

import (
	"context"
	"errors"
	"net/http"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// InitTelemetry bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func InitTelemetry(ctx context.Context, cfg TracingConfig) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	if !cfg.Enabled {
		return
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	var traceExporter tracesdk.SpanExporter
	if cfg.Endpoint != "" {
		traceExporter, err = otelTraceExporter(cfg.Endpoint, ctx)
		if err != nil {
			handleErr(err)
			return
		}
	} else {
		traceExporter, err = stdoutTraceExporter()
		if err != nil {
			handleErr(err)
			return
		}
	}

	otelResource, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcessExecutableName(),
		resource.WithProcessOwner(),
		resource.WithProcessRuntimeVersion(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
		),
	)

	if err != nil {
		handleErr(err)
		return
	}

	shutdownFuncs = append(shutdownFuncs, traceExporter.Shutdown)
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(traceExporter,
			tracesdk.WithBatchTimeout(5*time.Second)),
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(cfg.SampleRatio)),
		tracesdk.WithResource(otelResource),
	)
	otel.SetTracerProvider(tracerProvider)
	return

}

func otelTraceExporter(endpoint string, ctx context.Context) (*otlptrace.Exporter, error) {
	traceExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
	if err != nil {
		return nil, err
	}

	return traceExporter, nil
}

func stdoutTraceExporter() (*stdouttrace.Exporter, error) {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	return traceExporter, nil
}

// TraceOperation Wraps a block with a span
func TraceOperation(ctx context.Context, spanName string) (context.Context, trace.Span) {
	// If ctx is a gin context, then we need to walk to the request context from it
	req, ok := ctx.Value(0).(*http.Request)
	if ok {
		ctx = req.Context()
	}
	t := otel.Tracer("github.com/Layr-Labs/eigenda/tracing")
	ctx, span := t.Start(ctx, spanName)

	// Get stack trace
	stackTrace := make([]byte, 1<<16)
	stackTrace = stackTrace[:runtime.Stack(stackTrace, false)]

	span.SetAttributes(attribute.String("stack.trace", string(stackTrace)))

	return ctx, span
}
