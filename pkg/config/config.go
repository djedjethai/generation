package config

import (
	"github.com/djedjethai/generation/pkg/deleter"
	"github.com/djedjethai/generation/pkg/getter"
	"github.com/djedjethai/generation/pkg/setter"
)

type Services struct {
	Setter  setter.Setter
	Getter  getter.Getter
	Deleter deleter.Deleter
}

type Config struct {
	EncryptKEY       string
	Port             string
	PortGRPC         string
	FileLoggerActive bool
	DBLoggerActive   bool
	Shards           int
	ItemsPerShard    int
	Protocol         string
	IsTracing        bool
	IsMetrics        bool
	ServiceName      string
	JaegerEndpoint   string
}

type PostgresDBParams struct {
	DbName   string
	Host     string
	User     string
	Password string
}
