package agent

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/djedjethai/generation/internal/config"
	"github.com/djedjethai/generation/internal/deleter"
	"github.com/djedjethai/generation/internal/discovery"
	"github.com/djedjethai/generation/internal/getter"
	"github.com/djedjethai/generation/internal/handlers/grpc"
	"github.com/djedjethai/generation/internal/handlers/rest"
	"github.com/djedjethai/generation/internal/logger"
	"github.com/djedjethai/generation/internal/observability"
	"github.com/djedjethai/generation/internal/setter"
	"github.com/djedjethai/generation/internal/storage"
	"github.com/hashicorp/raft"
	"github.com/soheilhy/cmux"
	gglGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Agent struct {
	config Config

	mux          cmux.CMux
	server       *gglGrpc.Server
	Storage      *storage.DistributedStorage
	membership   *discovery.Membership
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
	Bootstrap       bool   // TODO to add
	DataDir         string // TODO to add
	//
	// ShardedMap     storage.ShardedMap
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

func New(cfg Config) (*Agent, error) {
	a := &Agent{
		config: cfg,
	}

	err := a.setupMux()
	if err != nil {
		return a, err
	}

	// set the storage
	err = a.setupStorage(cfg.Shards, cfg.ItemsPerShard)
	if err != nil {
		return a, err
	}

	// set services
	a.setupServices()

	// set loggerFacade
	err = a.setupLoggerFacade()
	if err != nil {
		return a, err
	}

	// set servers
	err = a.setupServers()
	if err != nil {
		return a, err
	}

	// set the membership
	err = a.setupMembership()
	if err != nil {
		return a, err
	}

	go a.serve()

	return a, nil
}

func (a *Agent) setupMux() error {
	rpcAddr := fmt.Sprintf(
		":%d",
		a.config.PortGRPC,
	)
	ln, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		return err
	}
	a.mux = cmux.New(ln)
	return nil
}

func (a *Agent) setupStorage(shards, itemsPerShard int) error {
	if shards > 0 && itemsPerShard > 0 {
		// TODO replace shardedMap with distributed.go, New
		// shardedMap := storage.NewShardedMap(shards, itemsPerShard, a.config.Observability)
		// TODO replace ....
		// a.config.ShardedMap = shardedMap

		raftLn := a.mux.Match(func(reader io.Reader) bool {
			b := make([]byte, 1)
			if _, err := reader.Read(b); err != nil {
				return false
			}
			return bytes.Compare(b, []byte{byte(storage.RaftRPC)}) == 0
		})

		logConfig := storage.Config{}
		logConfig.Raft.StreamLayer = storage.NewStreamLayer(
			raftLn,
			a.config.ServerTLSConfig,
			a.config.PeerTLSConfig,
		)
		logConfig.Raft.LocalID = raft.ServerID(a.config.NodeName)
		logConfig.Raft.Bootstrap = a.config.Bootstrap

		var err error
		a.Storage, err = storage.NewDistributedStorage(
			a.config.DataDir,
			logConfig,
			shards,
			itemsPerShard,
			a.config.Observability,
		)
		if err != nil {
			return err
		}

		if a.config.Bootstrap {
			err = a.Storage.WaitForLeader(3 * time.Second)
		}

		return err

	} else {
		return errors.New("The key value store can not work without storage")
	}
}

func (a *Agent) setupServices() {
	setSrv := setter.NewSetter(a.Storage, a.config.Observability)
	getSrv := getter.NewGetter(a.Storage, a.config.Observability)
	delSrv := deleter.NewDeleter(a.Storage, a.config.Observability)

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

func (a *Agent) setupServers() error {

	if a.config.Protocol == "grpc" {
		return a.runGRPC()
	} else if a.config.Protocol == "http" {
		// return a.runHTTP()
		return errors.New("HTTP protocol is not implemented")
	} else {
		return errors.New("Error start server, protocol is not defined")
	}
}

func (a *Agent) runGRPC() error {
	// TODO remove that, it should go trought the config
	l, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", a.config.PortGRPC))
	// TODO if run in docker
	// l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", a.config.PortGRPC))
	// 	if err != nil {
	// 		return func() {}, err
	// 	}

	// TODO to remove into configs...
	// set tls
	serverTLSConfig, _ := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: l.Addr().String(),
	})
	a.config.ServerTLSConfig = serverTLSConfig
	// end to remove

	var opts []gglGrpc.ServerOption
	if a.config.ServerTLSConfig != nil {
		// serverCreds := credentials.NewTLS(serverTLSConfig)
		serverCreds := credentials.NewTLS(a.config.ServerTLSConfig)
		opts = append(opts, gglGrpc.Creds(serverCreds))
	}

	var err error
	a.server, err = grpc.NewGRPCServer(a.config.Services, a.config.LoggerFacade, opts...)
	if err != nil {
		return err
	}

	grpcLn := a.mux.Match(cmux.Any())
	go func() {
		if err := a.server.Serve(grpcLn); err != nil {
			_ = a.Shutdown()
		}
	}()
	return err
}

func (a *Agent) setupMembership() error {
	rpcAddr, err := a.config.RPCAddr()
	if err != nil {
		return err
	}
	a.membership, err = discovery.New(a.Storage, discovery.Config{
		NodeName: a.config.NodeName,
		BindAddr: a.config.BindAddr,
		Tags: map[string]string{
			"rpc_addr": rpcAddr,
		},
		StartJoinAddrs: a.config.StartJoinAddrs,
	})
	return err
}

func (a *Agent) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()

	if a.shutdown {
		return nil
	}

	a.shutdown = true
	close(a.shutdowns)

	shutdown := []func() error{
		a.membership.Leave,
		func() error {
			a.server.GracefulStop()
			return nil
		},
		a.Storage.Close,
	}
	for _, fn := range shutdown {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Agent) serve() error {
	if err := a.mux.Serve(); err != nil {
		_ = a.Shutdown()
		return err
	}
	return nil
}

// TODO ........
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
