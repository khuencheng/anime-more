package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"log"
	_ "embed"
)

var config *viper.Viper

//go:embed conf.toml
var confFile []byte

func Init() {
	config = viper.New()
	config.SetConfigType("toml")
	err := config.ReadConfig(bytes.NewBuffer(confFile))
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
	fmt.Println(config.AllKeys())
}

func GetConfig() *viper.Viper {
	return config
}
