package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	Redis    RedisConfiguration
	AppEnv   string
}

var Config *Configuration

func InitConfig() *Configuration {

	var configuration Configuration

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	configuration.AppEnv = os.Getenv("APP_ENV")

	configuration.Server.Port = os.Getenv("PORT")

	configuration.Server.AesPassphrase = os.Getenv("AES_PASSPHRASE")
	configuration.Server.Superadmin.Username = os.Getenv("SUPERADMIN_USERNAME")
	configuration.Server.Superadmin.Password = os.Getenv("SUPERADMIN_PASSWORD")

	configuration.Database.Username = os.Getenv("POSTGRES_USER")
	configuration.Database.Password = os.Getenv("POSTGRES_PASWORD")

	configuration.Redis.Password = os.Getenv("REDIS_PASSWORD")

	Config = &configuration

	return Config
}

func GetConfig() *Configuration {
	return Config
}
