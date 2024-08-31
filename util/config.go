package util

import "github.com/spf13/viper"

type Config struct {
	DBSource   string `mapstructure:"DB_SOURCE"`
	DBDriver   string `mapstructure:"DB_DRIVER"`
	HTTPServer string `mapstructure:"HTTP_SERVER_ADDRESS"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
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
