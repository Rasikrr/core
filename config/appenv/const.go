package appenv

const (
	AppEnv = "env"

	PostgresDSN      = "postgres_dsn"
	PostgresDBName   = "postgres_db"
	PostgresUser     = "postgres_user"
	PostgresPassword = "postgres_password"
	PostgresHost     = "postgres_host"
	PostgresPort     = "postgres_port"

	RedisPassword    = "redis_password"
	RedisHost        = "redis_host"
	RedisPort        = "redis_port"
	RedisUser        = "redis_user"
	RedisDB          = "redis_db"
	RedisPoolSize    = "redis_pool_size"
	RedisMinIdle     = "redis_min_idle"
	RedisMaxIdle     = "redis_max_idle"
	RedisReadTimeout = "redis_read_timeout"
)
