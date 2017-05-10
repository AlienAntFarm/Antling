package config

import (
	"encoding/json"
	"fmt"
	"github.com/alienantfarm/antling/client"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type Configuration struct {
	Id      int            `json:"id"`
	Debug   bool           `json:"debug"`
	Anthive string         `json:"anthive"`
	Client  *client.Client `json:"-"`
}

var config *Configuration

func Get() *Configuration {
	var encoder *json.Encoder
	var configFile = viper.GetString("CONFIG_FILE")

	if config == nil {
		config = &Configuration{Client: client.Get()}

		viper.SetConfigName(strings.TrimSuffix(configFile, filepath.Ext(configFile)))
		viper.AddConfigPath(fmt.Sprintf("$%s/", viper.Get("CONFIG")))
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// request the API for a new id
			glog.Info("no config found, registring the node")
			antling, err := config.Client.Antling.Create()
			if err != nil {
				glog.Fatalf("%s", err)
			}
			config.Id = antling.Id

			// prepare the file for encoding our config
			out, err := os.Create(configFile)
			if err != nil {
				glog.Fatalf("%s", err)
			}
			defer out.Close()
			encoder = json.NewEncoder(out)
		} else if err != nil {
			glog.Fatalf("when reading config file: %s", err)
		}

		// load the config into our struct
		err = viper.Unmarshal(config)
		if err != nil {
			glog.Fatalf("when unmarshalling the json: %s", err)
		}

		// if an id has been requested the encoder is set up
		if encoder != nil {
			glog.Infof("writing config to a file %s", configFile)
			err = encoder.Encode(config)
			if err != nil {
				glog.Fatalf("%s", err)
			}
		}
	}
	glog.Infof("config loaded")
	return config
}

func init() {
	viper.Set("PROJECT", "github.com/alienantfarm/antling")
	viper.Set("CONFIG", "ANTLING_CONFIG")
	viper.Set("CONFIG_FILE", "config.json")
	viper.BindEnv("Anthive", "ANTHIVE_URL")
}
