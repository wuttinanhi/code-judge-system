package configs

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() {
	var err error

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	err = viper.ReadInConfig()
	if err == nil {
		log.Println("Using .env file")
	}

	viper.AutomaticEnv()
}
