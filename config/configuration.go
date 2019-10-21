package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	AppEnv   string
}

func InitConfig() Configuration {
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

	configuration.Server.Port = os.Getenv("PORT")
	configuration.AppEnv = os.Getenv("APP_ENV")

	return configuration
}
