package config

type Cache struct {
	NewUserExpiration int `mapstructure:"new_user_expiration"`
	UserExpiration    int `mapstructure:"user_expiration"`
	NewTaskExpiration int `mapstructure:"new_task_expiration"`
	TaskExpiration    int `mapstructure:"task_expiration"`
}

type Superadmin struct {
	Username string
	Password string
}

type ServerConfiguration struct {
	Port            string
	AesPassphrase   string
	RsaKeySize      int  `mapstructure:"rsa_key_size"`
	ConnectionLimit int  `mapstructure:"connection_limt"`
	BodyCrypt       bool `mapstructure:"body_crypt"`
	Superadmin
	Cache
}
