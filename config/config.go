package config

import (
	"github.com/spf13/viper"
	"log"
)

var config *viper.Viper

func Init(env string) {
	v := viper.New()
	v.SetConfigName(env)
	v.SetConfigType("yaml")
	v.AddConfigPath("config/")

	err := v.ReadInConfig()

	if err != nil {
		log.Fatal("Unable to parse the configuration file:", err)
	}

	config = v
}

func GetConfig() *viper.Viper {
	return config
}
