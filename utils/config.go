package utils

import (
	"encoding/json"
	"flag"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

type Configuration struct {
	Id      int    `json:"id"`
	Debug   bool   `json:"debug"`
	Anthive string `json:"anthive"`
	Images  string `json:"-"`
}

func (c *Configuration) Save() error {
	f, err := os.Create("config.json") // create file in place
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(c)
}

func PreRun(cmd *cobra.Command, args []string) {
	// reinit args for glog
	os.Args = os.Args[:1]

	viper.Unmarshal(Config)     // this will load default config
	err := viper.ReadInConfig() // load config
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			glog.Infof("creating a new config file")
			// a bit crappy but should handle this correctly
			RegisterCommand.Run(cmd, args)
			viper.ReadInConfig()
		} else {
			glog.Fatalf("when reading config file: %s", err)
		}
	}

	err = viper.Unmarshal(Config)
	if err != nil {
		glog.Errorf("%q", Config)
		glog.Fatalf("when unmarshalling the json: %s", err)
	}
	if Config.Debug { // debug is same as -vvvvv
		verbosity = 5
	}
	flag.Set("v", strconv.Itoa(verbosity))
	flag.Set("logtostderr", "true")
	flag.Parse()
	glog.V(1).Infoln("debug mode enabled")
	glog.V(2).Infof("%q", Config)
}

var (
	verbosity int
	sep       = string(os.PathSeparator)
	Config    = &Configuration{}
	Command   = &cobra.Command{
		Use:              "antling",
		Short:            "Start the antling process",
		PersistentPreRun: PreRun,
	}
	RegisterCommand = &cobra.Command{
		Use:   "register",
		Short: "Will register this node to the AlienAntFarm API",
	}
)

func init() {
	debugMsg := "trigger debug logs, same as -vvvvv, take precedence over verbose flag"
	verboseMsg := "verbose output, can be stacked to increase verbosity"

	Command.PersistentFlags().Bool("debug", false, debugMsg)
	Command.PersistentFlags().CountVarP(&verbosity, "verbose", "v", verboseMsg)

	Command.AddCommand(RegisterCommand)

	viper.BindPFlag("Debug", Command.PersistentFlags().Lookup("debug"))
	viper.Set("PROJECT", "github.com/alienantfarm/antling")
	viper.Set("CONFIG", "ANTLING_CONFIG")

	// set Images path
	viper.Set("Images", Urlize("assets", "images"))
	viper.BindEnv("Anthive", "ANTHIVE_URL")

	viper.SetConfigName("config")
	viper.AddConfigPath("$" + viper.GetString("CONFIG") + sep)
	viper.AddConfigPath(".")
}
