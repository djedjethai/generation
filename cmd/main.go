package main

import (
	"context"
	"fmt"
	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/deleter"
	"github.com/djedjethai/generation0/pkg/getter"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	// "github.com/djedjethai/generation0/pkg/handlers/grcp"

	"github.com/djedjethai/generation0/pkg/handlers/grpc"
	pb "github.com/djedjethai/generation0/pkg/handlers/grpc/proto/keyvalue"
	"github.com/djedjethai/generation0/pkg/handlers/rest"
	lgr "github.com/djedjethai/generation0/pkg/logger"
	"github.com/djedjethai/generation0/pkg/setter"
	storage "github.com/djedjethai/generation0/pkg/storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	gglGrpc "google.golang.org/grpc"
)

var labels = []label.KeyValue{
	label.Key("application").String(serviceName),
	label.Key("container_id").String(os.Getenv("HOSTNAME")),
}
var requests metric.Int64Counter
var serviceName = "golru"

var encryptK = "PX9PHFrdn79ljrjLDZHlV1t+BdxHRFf5"
var port = ":8080"
var portGrpc = ":50051"
var fileLoggerActive = false
var dbLoggerActive = false
var shards = 3
var itemsPerShard = 20
var protocol = "http"

func main() {

	// conf, err := getConf()
	// if err != nil {
	// 	log.Fatal("Err reading the config file: ", err)
	// }

	log.Println("debug1")
	// ================ prometheus config =================
	prometheusExporter, err := prometheus.NewExportPipeline(prometheus.Config{})
	if err != nil {
		fmt.Println(err)
	}
	log.Println("debug2")

	// Get the meter provider from the exporter.
	mp := prometheusExporter.MeterProvider()

	log.Println("debug3")
	// Set it as the global meter provider.
	otel.SetMeterProvider(mp)

	log.Println("debug4")
	// // Register the exporter as the handler for the "/metrics" pattern.
	// http.Handle("/metrics", prometheusExporter)
	// // Start the HTTP server listening on port 3000.
	// log.Fatal(http.ListenAndServe(":3000", nil))
	go runPrometheusEndPoint(prometheusExporter)

	log.Println("debug5")
	// meter := otel.GetMeterProvider().Meter("golru")

	err = buildRequestsCounter()
	if err != nil {
		log.Println("Error from build request counter: ", err)
	}

	log.Println("debug6")

	ctx := context.Background()
	go updateMetrics(ctx)

	log.Println("debug7")
	// =========================================
	// storage(infra layer)
	// the first arg is the number of shard, the second the number of item/shard
	var shardedMap storage.ShardedMap
	if shards > 0 && itemsPerShard > 0 {
		shardedMap = storage.NewShardedMap(shards, itemsPerShard)
	} else {
		log.Fatal("The key value store can not work without storage")
	}

	setSrv := setter.NewSetter(shardedMap, labels, &requests)
	getSrv := getter.NewGetter(shardedMap)
	delSrv := deleter.NewDeleter(shardedMap)

	// set logger
	var postgresConfig = config.PostgresDBParams{}
	if dbLoggerActive {
		if dbLoggerActive {
			postgresConfig.Host = "localhost"
			postgresConfig.DbName = "transactions"
			postgresConfig.User = "postgres"
			postgresConfig.Password = "password"
		}
	}

	loggerFacade, err := lgr.NewLoggerFacade(setSrv, delSrv, fileLoggerActive, dbLoggerActive, postgresConfig, encryptK)
	defer loggerFacade.CloseFileLogger()

	// in case the srv crash, when start back it will read the logger and recover its state
	// logger, err := initializeTransactionLog(setSrv, delSrv, fileLoggerActive)
	if err != nil {
		log.Panic("Logger(s) initialization failed: ", err)
	}

	switch protocol {
	case "http":
		runHTTP(setSrv, getSrv, delSrv, loggerFacade, port)
	case "grpc":
		runGRPC(setSrv, getSrv, delSrv, loggerFacade, portGrpc)
	default:
		log.Fatalln("Invalid protocol...")
	}

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

func updateMetrics(ctx context.Context) {
	// Retrieve the meter from the meter provider.
	meter := otel.GetMeterProvider().Meter(serviceName)
	// Create the instruments that we'll use to report memory
	// and goroutine values. Error values ignored for brevity.
	mem, _ := meter.NewInt64UpDownCounter("memory_usage_bytes",
		metric.WithDescription("Amount of memory used."),
	)
	goroutines, _ := meter.NewInt64UpDownCounter("num_goroutines",
		metric.WithDescription("Number of running goroutines."),
	)
	var m runtime.MemStats
	for {
		runtime.ReadMemStats(&m)
		// Report the values to the instruments, and receive
		// metric.Measurement values in return.
		mMem := mem.Measurement(int64(m.Sys))
		mGoroutines := goroutines.Measurement(int64(runtime.NumGoroutine()))
		// Provide the measurements (and teh context and
		// labels) to the meter.
		meter.RecordBatch(ctx, labels, mMem, mGoroutines)
		time.Sleep(5 * time.Second)
	}
}

// func buildRuntimeObservers() {
// 	meter := otel.GetMeterProvider().Meter(serviceName)
// 	m := runtime.MemStats{}
// 	meter.NewInt64UpDownSumObserver("memory_usage_bytes",
// 		func(_ context.Context, result metric.Int64ObserverResult) {
// 			runtime.ReadMemStats(&m)
// 			result.Observe(int64(m.Sys), labels...)
// 		},
// 		metric.WithDescription("Amount of memory used."),
// 	)
// 	meter.NewInt64UpDownSumObserver("num_goroutines",
// 		func(_ context.Context, result metric.Int64ObserverResult) {
// 			result.Observe(int64(runtime.NumGoroutine()), labels...)
// 		},
// 		metric.WithDescription("Number of running goroutines."),
// 	)
// }

func getConf() (*config.Config, error) {

	path, _ := os.Getwd()

	configPath := filepath.Join(path, "../config.yaml")

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	c := &config.Config{}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func runGRPC(setSrv setter.Setter, getSrv getter.Getter, delSrv deleter.Deleter, loggerFacade *lgr.LoggerFacade, port string) {
	s := gglGrpc.NewServer()
	pb.RegisterKeyValueServer(s, &grpc.Server{
		SetSrv:       setSrv,
		GetSrv:       getSrv,
		DelSrv:       delSrv,
		LoggerFacade: loggerFacade,
	})

	// lis, err := net.Listen("tcp", ":50051")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func runHTTP(setSrv setter.Setter, getSrv getter.Getter, delSrv deleter.Deleter, loggerFacade *lgr.LoggerFacade, port string) {
	// handler(application layer)
	router := rest.Handler(setSrv, getSrv, delSrv, loggerFacade)

	fmt.Printf("***** Service listening on port %s *****", port)
	log.Fatal(http.ListenAndServe(port, router))
}
