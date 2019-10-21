package config

type DatabaseConfiguration struct {
	Driver   string
	Dbname   string
	Username string
	Password string
	Host     string
	Port     string
	LogMode  bool
	// connection pool
	MaxOpenConns int
	MaxIdleConns int
}
