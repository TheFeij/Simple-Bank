package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DatabaseDriver           string        `mapstructure:"DATABASE_DRIVER"`
	DatabaseSource           string        `mapstructure:"DATABASE_SOURCE"`
	ServerHost               string        `mapstructure:"SERVER_HOST"`
	ServerPort               string        `mapstructure:"SERVER_PORT"`
	TokenSymmetricKey        string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenAccessTokenDuration time.Duration `mapstructure:"TOKEN_ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path, name string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("json")

	viper.AutomaticEnv()

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
