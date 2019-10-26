package config

type RedisConfiguration struct {
	Addr     string
	Password string
	// connection pool
	MinIdleConns int `mapstructure:"min_idle_conns"`
	MaxConnAge   int `mapstructure:"max_conn_age"`
	PoolSize     int `mapstructure:"pool_size"`
}
