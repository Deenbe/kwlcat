package global

import (
	"context"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	otelTrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

var (
	Tracer otelTrace.Tracer
	TracerProvider otelTrace.TracerProvider
)

func initialiseXrayTrace(name string, logger *zap.Logger) func() {
	ctx := context.Background()
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "0.0.0.0:4317"
	}

	logger.Sugar().Infof("otel exporter endpoint is %s", endpoint)

	driver := otlpgrpc.NewDriver(otlpgrpc.WithInsecure(), otlpgrpc.WithEndpoint(endpoint))

	exporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		log.Fatal("failed to initialise xtray exporter: %v", err)
	}

	bsp := trace.NewBatchSpanProcessor(exporter)
	idg := xray.NewIDGenerator()
	res, err := resource.New(ctx, resource.WithAttributes(), resource.WithAttributes(semconv.ServiceNameKey.String(name)))
	if err != nil {
		log.Fatal("failed to initialise resource")
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
		trace.WithIDGenerator(idg),
	)

	pusher := controller.New(
		processor.New(simple.NewWithExactDistribution(), exporter),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(1*time.Second),
	)

	err = pusher.Start(ctx)
	if err != nil {
		log.Fatal("failed to start metric controller: %v", err)
	}

	otel.SetTracerProvider(tp)
	global.SetMeterProvider(pusher.MeterProvider())
	otel.SetTextMapPropagator(xray.Propagator{})

	Tracer = otel.Tracer(name)
	TracerProvider = tp
	return func() {
		err := pusher.Stop(ctx)
		log.Fatal("failed to stop metric controller: %v", err)
	}
}

func InitialiseTrace(name string, logger *zap.Logger) func() {
	return initialiseXrayTrace(name, logger)
}
