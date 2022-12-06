package config

import (
	// "crypto/tls"
	// "fmt"
	// "net"

	"github.com/djedjethai/generation/internal/deleter"
	"github.com/djedjethai/generation/internal/getter"
	// "github.com/djedjethai/generation/internal/observability"
	"github.com/djedjethai/generation/internal/setter"
	// "github.com/djedjethai/generation/internal/storage"
)

type Services struct {
	Setter  setter.Setter
	Getter  getter.Getter
	Deleter deleter.Deleter
}

type PostgresDBParams struct {
	DbName   string
	Host     string
	User     string
	Password string
}

// type Config struct {
// 	PortGRPC         int
// 	EncryptKEY       string
// 	Port             string
// 	FileLoggerActive bool
// 	DBLoggerActive   bool
// 	Shards           int
// 	ItemsPerShard    int
// 	Protocol         string
// 	IsTracing        bool
// 	IsMetrics        bool
// 	ServiceName      string
// 	JaegerEndpoint   string
// 	//
// 	ServerTLSConfig *tls.Config
// 	PeerTLSConfig   *tls.Config
// 	BindAddr        string
// 	NodeName        string
// 	StartJoinAddrs  []string
// 	//
// 	ShardedMap     storage.ShardedMap
// 	Observability  observability.Observability
// 	PostgresParams PostgresDBParams
// 	Services       Services
//
// 	// Services
// }
//
// func (c Config) RPCAddr() (string, error) {
// 	host, _, err := net.SplitHostPort(c.BindAddr)
// 	if err != nil {
// 		return "", err
// 	}
// 	return fmt.Sprintf("%s:%d", host, c.PortGRPC), nil
// }
