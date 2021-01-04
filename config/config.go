package config

import (
	"github.com/spf13/viper"
	"log"
)

var config *viper.Viper

func Init() {
	config = viper.New()
	config.SetConfigType("toml")
	config.SetConfigName("dev")
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
}

func GetConfig() *viper.Viper {
	return config
}
