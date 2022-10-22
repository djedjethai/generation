package main

import (
	// "fmt"
	"log"
	"net"

	// "net/http"

	// "github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/deleter"
	"github.com/djedjethai/generation0/pkg/getter"
	// "github.com/djedjethai/generation0/pkg/handlers/grcp"

	// "github.com/djedjethai/generation0/pkg/handlers/rest"
	// lgr "github.com/djedjethai/generation0/pkg/logger"
	"github.com/djedjethai/generation0/pkg/handlers/grpc"
	pb "github.com/djedjethai/generation0/pkg/handlers/grpc/proto/keyvalue"
	"github.com/djedjethai/generation0/pkg/setter"
	storage "github.com/djedjethai/generation0/pkg/storage"
	gglGrpc "google.golang.org/grpc"
)

// var store = make(map[string]*storage.Node)
var encryptK = "PX9PHFrdn79ljrjLDZHlV1t+BdxHRFf5"

var port = ":8080"

// default value
var fileLoggerActive = false
var dbLoggerActive = false

var shards = 2
var itemsPerShard = 25

func main() {

	// storage(infra layer)
	// the first arg is the number of shard, the second the number of item/shard
	var shardedMap storage.ShardedMap
	if shards > 0 && itemsPerShard > 0 {
		shardedMap = storage.NewShardedMap(shards, itemsPerShard)
	} else {
		log.Fatal("The key value store can not work without storage")
	}

	setSrv := setter.NewSetter(shardedMap)
	getSrv := getter.NewGetter(shardedMap)
	delSrv := deleter.NewDeleter(shardedMap)

	s := gglGrpc.NewServer()
	pb.RegisterKeyValueServer(s, &grpc.Server{
		SetSrv: setSrv,
		GetSrv: getSrv,
		DelSrv: delSrv,
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	// TODO rest implementation
	// var postgresConfig = config.PostgresDBParams{}
	// if dbLoggerActive {
	// 	if dbLoggerActive {
	// 		postgresConfig.Host = "localhost"
	// 		postgresConfig.DbName = "transactions"
	// 		postgresConfig.User = "postgres"
	// 		postgresConfig.Password = "password"
	// 	}
	// }

	// loggerFacade, err := lgr.NewLoggerFacade(setSrv, delSrv, fileLoggerActive, dbLoggerActive, postgresConfig, encryptK)
	// defer loggerFacade.CloseFileLogger()

	// // in case the srv crash, when start back it will read the logger and recover its state
	// // logger, err := initializeTransactionLog(setSrv, delSrv, fileLoggerActive)
	// if err != nil {
	// 	log.Panic("Logger(s) initialization failed: ", err)
	// }

	// // handler(application layer)
	// router := rest.Handler(setSrv, getSrv, delSrv, loggerFacade)

	// fmt.Printf("***** Service listening on port %s *****", port)
	// log.Fatal(http.ListenAndServe(port, router))
}
