package main

import (
	"context"
	"fmt"
	"github.com/djedjethai/generation/internal/agent"
	"github.com/djedjethai/generation/internal/config"
	"github.com/djedjethai/generation/internal/observability"
	// "github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	// "strconv"
	"time"
	// "go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"runtime"
)

var labels = []label.KeyValue{
	label.Key("application").String(appName),
	label.Key("service").String(serviceName),
	// label.Key("container_id").String(os.Getenv("HOSTNAME")),
	label.Key("container_id").String("1234"),
}

var requests metric.Int64Counter
var appName string
var serviceName string

func setupSrv(cfg *agent.Config) error {

	// config from var env
	port := os.Getenv("PORT")
	// port_grpc := os.Getenv("PORT_GRPC")
	protocol := os.Getenv("PROTOCOL")
	app_name := os.Getenv("APP_NAME")
	service_name := os.Getenv("SERVICE_NAME")
	// TODO add this varenv
	_ = os.Getenv("DATABASE_HOST")

	protocol = "grpc" // uncomment to switch to grpc

	setVarEnv(&protocol, &port, &app_name, &service_name)

	appName = app_name
	serviceName = service_name

	cfg.Port = port
	cfg.Protocol = protocol

	obs := observability.Observability{
		Requests:    &requests,
		Labels:      labels,
		IsTracing:   isTracing,
		IsMetrics:   isMetrics,
		ServiceName: serviceName,
	}
	cfg.Observability = &obs

	srvLog := observability.NewSrvLogger(logMode)
	obs.Logger = srvLog

	// set logger
	initLogger(logMode)

	// tracing is on
	if obs.IsTracing {
		// jaeger config, http://localhost:16686/search
		tr, err := configJaeger(cfg.JaegerEndpoint)
		if err != nil {
			log.Fatal("Error when configuring Jaeger: ", err)
		}

		obs.Tracer = tr
	}

	// metrics is on
	if obs.IsMetrics {
		// prometheus config, 127.0.0.1:9090
		configPrometheus()
	}

	// set logger
	if cfg.DBLoggerActive {
		var postgresConfig = config.PostgresDBParams{}
		// postgresConfig.Host = "localhost"
		postgresConfig.Host = "postgres" // in the docker-compose network
		postgresConfig.DbName = "transactions"
		postgresConfig.User = "postgres"
		postgresConfig.Password = "password"
		cfg.PostgresParams = postgresConfig
	}

	return nil
}

func setVarEnv(protocol, port, app_name, service_name *string) {
	if len(*protocol) == 0 {
		*protocol = "http"
	}
	if len(*port) == 0 && *protocol == "http" {
		*port = "8080"
	}
	if len(*app_name) == 0 {
		*app_name = "golru"
	}
	if len(*service_name) == 0 {
		*service_name = "generation"
	}
}

func initLogger(logMode string) {
	var cfg zap.Config
	if logMode == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {

		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.TimeKey = fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05")) // Turn off timestamp output
	cfg.Sampling = &zap.SamplingConfig{
		Initial:    100, // Allow first 3 events/second
		Thereafter: 100, // Allows 1 per 3 thereafter
		Hook: func(e zapcore.Entry, d zapcore.SamplingDecision) {
			if d == zapcore.LogDropped {
				fmt.Println("event dropped...")
			}
		},
	}
	logger, _ := cfg.Build()   // Constructs the new logger
	zap.ReplaceGlobals(logger) // Replace Zap's global logger
}

func configPrometheus() {
	prometheusExporter, err := prometheus.NewExportPipeline(prometheus.Config{})
	if err != nil {
		fmt.Println(err)
	}

	// Get the meter provider from the exporter.
	mp := prometheusExporter.MeterProvider()

	// Set it as the global meter provider.
	otel.SetMeterProvider(mp)

	// // Register the exporter as the handler for the "/metrics" pattern.
	// http.Handle("/metrics", prometheusExporter)
	// // Start the HTTP server listening on port 3000.
	// log.Fatal(http.ListenAndServe(":3000", nil))
	go runPrometheusEndPoint(prometheusExporter)

	// meter := otel.GetMeterProvider().Meter("golru")

	err = buildRequestsCounter()
	if err != nil {
		log.Println("Error from build request counter: ", err)
	}

	buildRuntimeObservers()
}

func configJaeger(jaegerEndpoint string) (observability.Tracer, error) {
	// stdExporter, err := stdout.NewExporter(
	// 	stdout.WithPrettyPrint(),
	// )
	// if err != nil {
	// 	log.Println("Error creating a Jaeger new exporter: ", err)
	// 	var ct config.Tracer
	// 	return ct, err
	// }

	jaegerExporter, err := jaeger.NewRawExporter(
		jaeger.WithCollectorEndpoint(jaegerEndpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: appName,
		}),
	)
	if err != nil {
		log.Println("Error creating a Jaeger new rawExporter: ", err)
		var ct observability.Tracer
		return ct, err
	}

	tp := sdktrace.NewTracerProvider(
		// sdktrace.WithSyncer(stdExporter),
		sdktrace.WithSyncer(jaegerExporter),
		// sdktrace.WithResource(resource.NewWithAttributes(
		// 	semconv.SchemaURL,
		// 	semconv.ServiceNameKey.String(serviceName))),
	)

	otel.SetTracerProvider(tp)

	// Setting the global tracer provider makes it discoverable via the otel.GetTracerPro
	// vider function. This allows libraries and other dependencies that use the OpenTele‐
	// metry API to more easily discover the SDK and emit telemetry data:
	// gtp := otel.GetTracerProvider(tp)
	tr := otel.GetTracerProvider().Tracer(serviceName)

	return tr, nil
}

func runPrometheusEndPoint(prometheusExporter *prometheus.Exporter) {
	// Register the exporter as the handler for the "/metrics" pattern.
	http.Handle("/metrics", prometheusExporter)
	// Start the HTTP server listening on port 3000.
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func buildRequestsCounter() error {
	var err error
	// Retrieve the meter from the meter provider.
	meter := otel.GetMeterProvider().Meter(serviceName)
	// Get an Int64Counter for a metric called "fibonacci_requests_total".
	requests, err = meter.NewInt64Counter("golru_requests_total",
		metric.WithDescription("Total number of golru requests."),
	)
	return err
}

// the NewInt64UpDownSumObserver accepts the name of the metric as a
// string, something called a Int64ObserverFunc, and zero or more instrument
// options (such as the metric description)
func buildRuntimeObservers() {
	meter := otel.GetMeterProvider().Meter(serviceName)
	m := runtime.MemStats{}
	meter.NewInt64UpDownSumObserver("memory_usage_bytes",
		func(_ context.Context, result metric.Int64ObserverResult) {
			runtime.ReadMemStats(&m)
			result.Observe(int64(m.Sys), labels...)
		},
		metric.WithDescription("Amount of memory used."),
	)
	meter.NewInt64UpDownSumObserver("num_goroutines",
		func(_ context.Context, result metric.Int64ObserverResult) {
			result.Observe(int64(runtime.NumGoroutine()), labels...)
		},
		metric.WithDescription("Number of running goroutines."),
	)
}
