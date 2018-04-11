package config

import (
	"github.com/spf13/viper"
	"log"
)

var config *viper.Viper

//Initialises the config for the specified environment
func Initialize(env string) {
	v := viper.New()
	v.SetConfigName(env)
	v.SetConfigType("yaml")
	v.AddConfigPath("config/")

	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("Unable to parse the configuration file: ", err)
	}

	config = v
}

//Returns the initialized config. Do not call before calling Initialize!
func GetConfig() *viper.Viper {
	return config
}
