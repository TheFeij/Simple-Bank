package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database Database `mapstructure:"DATABASE"`
	Server   Server   `mapstructure:"SERVER"`
	Token    Token    `mapstructure:"TOKEN"`
}

type Database struct {
	Driver string `mapstructure:"DRIVER"`
	Source string `mapstructure:"SOURCE"`
}

type Server struct {
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
}

type Token struct {
	SymmetricKey        string        `mapstructure:"SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path, name string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
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
