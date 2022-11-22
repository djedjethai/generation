package main

import (
	"fmt"
	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/deleter"
	"github.com/djedjethai/generation0/pkg/getter"
	"github.com/djedjethai/generation0/pkg/handlers/grpc"
	pb "github.com/djedjethai/generation0/pkg/handlers/grpc/proto/keyvalue"
	"github.com/djedjethai/generation0/pkg/handlers/rest"
	lgr "github.com/djedjethai/generation0/pkg/logger"
	"github.com/djedjethai/generation0/pkg/setter"
	storage "github.com/djedjethai/generation0/pkg/storage"
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

	// set logger
	var postgresConfig = config.PostgresDBParams{}
	if cfg.DBLoggerActive {
		postgresConfig.Host = "localhost"
		postgresConfig.DbName = "transactions"
		postgresConfig.User = "postgres"
		postgresConfig.Password = "password"
	}

	loggerFacade, err := lgr.NewLoggerFacade(setSrv, delSrv, cfg.FileLoggerActive, cfg.DBLoggerActive, postgresConfig, cfg.EncryptKEY)
	defer loggerFacade.CloseFileLogger()

	// in case the srv crash, when start back it will read the logger and recover its state
	// logger, err := initializeTransactionLog(setSrv, delSrv, fileLoggerActive)
	if err != nil {
		log.Panic("Logger(s) initialization failed: ", err)
	}

	switch cfg.Protocol {
	case "http":
		runHTTP(setSrv, getSrv, delSrv, loggerFacade, cfg.Port)
	case "grpc":
		runGRPC(setSrv, getSrv, delSrv, loggerFacade, cfg.PortGRPC)
	default:
		log.Fatalln("Invalid protocol...")
	}

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
