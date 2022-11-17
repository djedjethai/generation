package config

import (
	"context"
	"github.com/djedjethai/generation0/pkg/logger"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	// "go.uber.org/zap"
)

type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanOption) (context.Context, trace.Span)
}

type Observability struct {
	Logger      *logger.SrvLogger
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
