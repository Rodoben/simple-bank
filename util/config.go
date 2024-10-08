package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource     string        `mapstructure:"DB_SOURCE"`
	DBDriver     string        `mapstructure:"DB_DRIVER"`
	HTTPServer   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	AuthTokenKey string        `mapstructure:"AUTH_TOKEN_KEY"`
	TokenExpiry  time.Duration `mapstructure:"TOKEN_EXPIRY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
