package agent

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/djedjethai/generation/internal/config"
	"github.com/djedjethai/generation/internal/deleter"
	"github.com/djedjethai/generation/internal/getter"
	"github.com/djedjethai/generation/internal/handlers/grpc"
	"github.com/djedjethai/generation/internal/handlers/rest"
	"github.com/djedjethai/generation/internal/logger"
	"github.com/djedjethai/generation/internal/observability"
	"github.com/djedjethai/generation/internal/setter"
	"github.com/djedjethai/generation/internal/storage"
	gglGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"
	"sync"
)

type Agent struct {
	config       Config
	server       *grpc.Server
	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
}

type Config struct {
	PortGRPC int
	// EncryptKEY       string
	Port             string
	FileLoggerActive bool
	DBLoggerActive   bool
	Shards           int
	ItemsPerShard    int
	Protocol         string
	IsTracing        bool
	IsMetrics        bool
	ServiceName      string
	JaegerEndpoint   string
	//
	ServerTLSConfig *tls.Config
	PeerTLSConfig   *tls.Config
	BindAddr        string
	NodeName        string
	StartJoinAddrs  []string
	//
	ShardedMap     storage.ShardedMap
	Observability  *observability.Observability
	PostgresParams config.PostgresDBParams
	Services       config.Services
	LoggerFacade   *logger.LoggerFacade
}

func (c Config) RPCAddr() (string, error) {
	host, _, err := net.SplitHostPort(c.BindAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, c.PortGRPC), nil
}

func New(cfg Config) (*Agent, func(), error) {
	a := &Agent{
		config: cfg,
	}

	// set the storage
	err := a.setupStorage(cfg.Shards, cfg.ItemsPerShard)
	if err != nil {
		return a, func() {}, err
	}

	// set services
	a.setupServices()

	// set loggerFacade
	err = a.setupLoggerFacade()
	if err != nil {
		return a, func() {}, err
	}

	// set servers
	fn, err := a.setupServers()
	if err != nil {
		return a, func() {}, err
	}
	return a, fn, nil
}

func (a *Agent) setupStorage(shards, itemsPerShard int) error {
	if shards > 0 && itemsPerShard > 0 {
		// TODO replace shardedMap with distributed.go, New
		shardedMap := storage.NewShardedMap(shards, itemsPerShard, a.config.Observability)
		// TODO replace ....
		a.config.ShardedMap = shardedMap
		return nil
	} else {
		return errors.New("The key value store can not work without storage")
	}
}

func (a *Agent) setupServices() {
	setSrv := setter.NewSetter(a.config.ShardedMap, a.config.Observability)
	getSrv := getter.NewGetter(a.config.ShardedMap, a.config.Observability)
	delSrv := deleter.NewDeleter(a.config.ShardedMap, a.config.Observability)

	a.config.Services = config.Services{setSrv, getSrv, delSrv}
}

func (a *Agent) setupLoggerFacade() error {
	// TODO see the story of *services or not....
	lgrF, err := logger.NewLoggerFacade(a.config.Services, a.config.DBLoggerActive, a.config.PostgresParams)
	if err != nil {
		return err
	}

	a.config.LoggerFacade = lgrF
	return nil
}

func (a *Agent) setupServers() (func(), error) {

	if a.config.Protocol == "grpc" {
		return a.runGRPC()
	} else if a.config.Protocol == "http" {
		return a.runHTTP()
	} else {
		return func() {}, errors.New("Error start server, protocol is not defined")
	}
}

func (a *Agent) runGRPC() (func(), error) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", a.config.PortGRPC))
	// TODO if run in docker
	// l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", a.config.PortGRPC))
	if err != nil {
		return func() {}, err
	}

	// set tls
	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: l.Addr().String(),
	})

	serverCreds := credentials.NewTLS(serverTLSConfig)

	server, err := grpc.NewGRPCServer(a.config.Services, a.config.LoggerFacade, gglGrpc.Creds(serverCreds))
	if err != nil {
		return func() {}, err
	}

	err = server.Serve(l)
	if err != nil {
		return func() {}, err
	}
	return func() {
		defer server.Stop()
		defer l.Close()
	}, nil

}

func (a *Agent) runHTTP() (func(), error) {
	hdl := rest.NewHandler(a.config.Services, a.config.LoggerFacade)
	router := hdl.Multiplex()

	fmt.Printf("***** Service listening on port %s *****", a.config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", a.config.Port), router)
	if err != nil {
		return func() {}, err
	}
	return func() {}, nil
}
