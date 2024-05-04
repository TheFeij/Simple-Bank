package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DatabaseDriver            string        `mapstructure:"DATABASE_DRIVER"`
	DatabaseSource            string        `mapstructure:"DATABASE_SOURCE"`
	HTTPServerHost            string        `mapstructure:"HTTP_SERVER_HOST"`
	HTTPServerPort            string        `mapstructure:"HTTP_SERVER_PORT"`
	GrpcServerHost            string        `mapstructure:"GRPC_SERVER_HOST"`
	GrpcServerPort            string        `mapstructure:"GRPC_SERVER_PORT"`
	TokenSymmetricKey         string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenAccessTokenDuration  time.Duration `mapstructure:"TOKEN_ACCESS_TOKEN_DURATION"`
	TokenRefreshTokenDuration time.Duration `mapstructure:"TOKEN_REFRESH_TOKEN_DURATION"`
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
