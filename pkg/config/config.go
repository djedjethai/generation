package config

import (
	"context"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanOption) (context.Context, trace.Span)
}

type Observability struct {
	Requests    *metric.Int64Counter
	Labels      []label.KeyValue
	Tracer      Tracer
	IsTracing   bool
	IsMetrics   bool
	ServiceName string
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
	IsTracing        bool   `yaml:"isTracing"`
	IsMetrics        bool   `yaml:"isMetrics"`
	ServiceName      string `yaml:"serviceName"`
	JaegerEndPoint   string `yaml:"jaegerEndPoint"`
}

type PostgresDBParams struct {
	DbName   string
	Host     string
	User     string
	Password string
}
