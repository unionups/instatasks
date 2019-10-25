package config

type Cache struct {
	NewUserExpiration int
	UserExpiration    int
}

type Superadmin struct {
	Username string
	Password string
}

type ServerConfiguration struct {
	Port          string
	AesPassphrase string
	Superadmin
	Cache
}
