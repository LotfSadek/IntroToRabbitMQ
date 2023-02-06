package main

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ConnectionString string
	AMQPURL          string
}

func LoadConfig(filePath string) (*Config, error) {
	connectionString := os.Getenv("CONNECTION_STRING")
	if connectionString == "" {
		return nil, errors.New("CONNECTION_STRING environment variable is not set")
	}

	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		return nil, errors.New("AMQP_URL environment variable is not set")
	}

	v := viper.New()
	v.SetConfigFile(filePath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
