package configs

import (
	"github.com/spf13/viper"
)

func LoadConfig() {
	var err error

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		panic("Error reading .env file")
	}

	viper.AutomaticEnv()
}
