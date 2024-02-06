package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database Database `mapstructure:"Database"`
	Server   Server   `mapstructure:"Server"`
}

type Database struct {
	Driver string `mapstructure:"Driver"`
	Source string `mapstructure:"Source"`
}

type Server struct {
	Host string `mapstructure:"Host"`
	Port string `mapstructure:"Port"`
}

func LoadConfig(path string) (Config, error) {
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
