package config

const (
	MemcachedPort = "127.0.0.1:11211"
	ServerPort    = ":9000"
	DBhost        = "localhost"
	DBport        = 5432
	DBuser        = "postgres"
	DBpassword    = "123"
	DBname        = "ozon"

	StudentsCacheExpiration  = 3600
	TasksCacheExpiration     = 604800
	SolutionsCacheExpiration = 600
)
