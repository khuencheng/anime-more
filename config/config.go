package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var config *viper.Viper

//go:embed conf.toml
var confFile []byte

func Init() {

}

func GetConfig() *viper.Viper {
	config = viper.New()
	config.SetConfigType("toml")
	err := config.ReadConfig(bytes.NewBuffer(confFile))
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
	fmt.Println(config.AllKeys())
	return config
}
