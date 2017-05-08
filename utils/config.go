package utils

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/viper"
)

type Configuration struct {
	Id      int    `json:"id"`
	Debug   bool   `json:"debug"`
	Anthive string `json:"anthive"`
}

var config *Configuration

func Config() *Configuration {
	if config == nil {
		config = &Configuration{}

		viper.SetConfigName("config")
		viper.AddConfigPath(fmt.Sprintf("$%s/", viper.Get("CONFIG")))
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			glog.Info("no config found, registring the node")
		} else if err != nil {
			glog.Fatalf("when reading config file: %s", err)
		}
		err = viper.Unmarshal(config)
		if err != nil {
			glog.Fatalf("when unmarshalling the json: %s", err)
		}
	}
	return config
}

func init() {
	viper.Set("PROJECT", "github.com/alienantfarm/antling")
	viper.Set("CONFIG", "ANTLING_CONFIG")
}
