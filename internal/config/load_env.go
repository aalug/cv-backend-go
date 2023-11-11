package config

import (
	"errors"
	"github.com/spf13/viper"
	"os"
)

func LoadConfigFromFile(path string) (cfg Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	return
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() (cfg Config, err error) {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	dbSource := os.Getenv("DB_SOURCE")
	dbDriver := os.Getenv("DB_DRIVER")

	if serverAddress == "" || dbSource == "" || dbDriver == "" {
		return Config{}, errors.New("missing required environment variable")
	}

	cfg.ServerAddress = serverAddress
	cfg.DBSource = dbSource
	cfg.DBDriver = dbDriver
	return cfg, nil
}
