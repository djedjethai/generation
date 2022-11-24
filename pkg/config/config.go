package config

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
