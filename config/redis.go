package config

type RedisConfiguration struct {
	Addr     string
	Password string
	// connection pool
	MinIdleConns int
	MaxConnAge   int
	PoolSize     int
}
