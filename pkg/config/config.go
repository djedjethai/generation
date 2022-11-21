package config

import (
	"context"
	"github.com/djedjethai/generation0/pkg/serviceLogger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanOption) (context.Context, trace.Span)
}

type Observability struct {
	Logger      *serviceLogger.SrvLogger
	Requests    *metric.Int64Counter
	Labels      []label.KeyValue
	Tracer      Tracer
	IsTracing   bool
	IsMetrics   bool
	ServiceName string
}

func (o *Observability) StartTrace(ctx context.Context, traceName string) (context.Context, func()) {

	if o.IsTracing {
		ctx1, sp := o.Tracer.Start(context.Background(), traceName)
		return ctx1, func() {
			defer sp.End()
		}
	}

	return ctx, func() {}
}

func (o *Observability) CarryOnTrace(ctx context.Context, traceName string) func() {

	if o.IsTracing {
		tr := otel.GetTracerProvider().Tracer(o.ServiceName)
		_, sp := tr.Start(ctx, traceName)
		return func() {
			defer sp.End()
		}
	}

	return func() {}
}

func (o *Observability) AddMetrics(ctx context.Context) {
	if o.IsMetrics {
		o.Requests.Add(ctx, 1, o.Labels...)
	}
}

func (o *Observability) AddMetricsAndSpecificLabel(ctx context.Context, key, val string) {
	if o.IsMetrics {
		lb := label.Key(key).String(val)
		o.Requests.Add(ctx, 1, lb)
	}
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
