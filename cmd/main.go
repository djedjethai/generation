package main

import (
	"fmt"
	"github.com/djedjethai/generation/pkg/config"
	"github.com/djedjethai/generation/pkg/deleter"
	"github.com/djedjethai/generation/pkg/getter"
	"github.com/djedjethai/generation/pkg/handlers/grpc"
	"github.com/djedjethai/generation/pkg/handlers/rest"
	lgr "github.com/djedjethai/generation/pkg/logger"
	"github.com/djedjethai/generation/pkg/setter"
	storage "github.com/djedjethai/generation/pkg/storage"
	gglGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {

	cfg, obs, postgresConfig, err := setupSrv()
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

	loggerFacade, err := lgr.NewLoggerFacade(services, cfg.DBLoggerActive, postgresConfig)

	// in case the srv crash, when start back it will read the logger and recover its state
	// logger, err := initializeTransactionLog(setSrv, delSrv, fileLoggerActive)
	if err != nil {
		log.Panic("Logger(s) initialization failed: ", err)
	}

	switch cfg.Protocol {
	case "http":
		runHTTP(&services, loggerFacade, cfg.Port)
	case "grpc":
		srv, l := runGRPC(&services, loggerFacade, cfg.PortGRPC)
		if err := srv.Serve(l); err != nil {
			log.Fatal("Error run grpc server: ", err)
		}
		defer srv.Stop()
		defer l.Close()

	default:
		log.Fatalln("Invalid protocol...")
	}

}

func runGRPC(services *config.Services, loggerFacade *lgr.LoggerFacade, port string) (*gglGrpc.Server, net.Listener) {

	l, err := net.Listen("tcp", fmt.Sprintf("%s%s", "127.0.0.1", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// set tls
	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: l.Addr().String(),
	})

	serverCreds := credentials.NewTLS(serverTLSConfig)

	server, err := grpc.NewGRPCServer(services, loggerFacade, gglGrpc.Creds(serverCreds))
	if err != nil {
		log.Fatal("Error create GRPC server: ", err)
	}

	return server, l
}

func runHTTP(services *config.Services, loggerFacade *lgr.LoggerFacade, port string) {
	// handler(application layer)
	hdl := rest.NewHandler(services, loggerFacade)
	router := hdl.Multiplex()

	fmt.Printf("***** Service listening on port %s *****", port)
	log.Fatal(http.ListenAndServe(port, router))
}
