package main

import (
	"fmt"
	pb "github.com/djedjethai/generation/api/v1/keyvalue"
	"github.com/djedjethai/generation/pkg/config"
	"github.com/djedjethai/generation/pkg/deleter"
	"github.com/djedjethai/generation/pkg/getter"
	"github.com/djedjethai/generation/pkg/handlers/grpc"
	"github.com/djedjethai/generation/pkg/handlers/rest"
	lgr "github.com/djedjethai/generation/pkg/logger"
	"github.com/djedjethai/generation/pkg/setter"
	storage "github.com/djedjethai/generation/pkg/storage"
	gglGrpc "google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {

	cfg, obs, err := setupSrv()
	if err != nil {
		os.Exit(1)
	}

	// storage(infra layer)
	// the first arg is the number of shard, the second the number of item/shard
	var shardedMap storage.ShardedMap
	if cfg.Shards > 0 && cfg.ItemsPerShard > 0 {
		shardedMap = storage.NewShardedMap(cfg.Shards, cfg.ItemsPerShard, obs)
	} else {
		log.Fatal("The key value store can not work without storage")
	}

	setSrv := setter.NewSetter(shardedMap, obs)
	getSrv := getter.NewGetter(shardedMap, obs)
	delSrv := deleter.NewDeleter(shardedMap, obs)

	services := config.Services{setSrv, getSrv, delSrv}

	// set logger
	var postgresConfig = config.PostgresDBParams{}
	if cfg.DBLoggerActive {
		// postgresConfig.Host = "localhost"
		postgresConfig.Host = "postgres" // in the docker-compose network
		postgresConfig.DbName = "transactions"
		postgresConfig.User = "postgres"
		postgresConfig.Password = "password"
	}

	loggerFacade, err := lgr.NewLoggerFacade(setSrv, delSrv, cfg.DBLoggerActive, postgresConfig)

	// in case the srv crash, when start back it will read the logger and recover its state
	// logger, err := initializeTransactionLog(setSrv, delSrv, fileLoggerActive)
	if err != nil {
		log.Panic("Logger(s) initialization failed: ", err)
	}

	switch cfg.Protocol {
	case "http":
		runHTTP(&services, loggerFacade, cfg.Port)
	case "grpc":
		runGRPC(&services, loggerFacade, cfg.PortGRPC)
	default:
		log.Fatalln("Invalid protocol...")
	}

}

func runGRPC(services *config.Services, loggerFacade *lgr.LoggerFacade, port string) {
	s := gglGrpc.NewServer()
	pb.RegisterKeyValueServer(s, &grpc.Server{
		Services:     services,
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

func runHTTP(services *config.Services, loggerFacade *lgr.LoggerFacade, port string) {
	// handler(application layer)
	hdl := rest.NewHandler(services, loggerFacade)
	router := hdl.Multiplex()

	fmt.Printf("***** Service listening on port %s *****", port)
	log.Fatal(http.ListenAndServe(port, router))
}
