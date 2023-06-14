package config

const (
	Service       = "api"
	Environment   = "development"
	MemcachedPort = "127.0.0.1:11211"
	HttpPort      = ":9000"

	GrpcPort       = ":50051"
	PrometheusPort = ":9091"
	TracesPort     = "http://localhost:14268"
	DBhost         = "localhost"
	DBport         = 5432
	DBuser         = "postgres"
	DBpassword     = "123"
	DBname         = "ozon"

	StudentsCacheExpiration = 3600
)
