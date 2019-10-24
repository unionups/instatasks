package config

type Cache struct {
	NewUserExpiration int
	UserExpiration    int
}

type ServerConfiguration struct {
	Port string
	Cache
}
