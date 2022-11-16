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
