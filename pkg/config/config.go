package config

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanOption) (context.Context, trace.Span)
}

type Config struct {
	EncryptKEY       string `yaml:"encryptK"`
	Port             string `yaml:"port"`
	PortGRPC         string `yaml:"portGrpc"`
	FileLoggerActive bool   `yaml:"fileLoggerActive"`
	DBLoggerActive   bool   `yaml:"dbLoggerActive"`
	Shards           int    `yaml:"shards"`
	ItemsPerShard    int    `yaml:"itemsPerShard"`
	Protocol         string `yaml:"protocol"`
}

type PostgresDBParams struct {
	DbName   string
	Host     string
	User     string
	Password string
}
