package utils

import (
	"flag"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"strconv"
	"text/template"
)

const CONFIG_NAME = "antling.toml"
const IMAGES_PREFIX = "/static/images"

type Configuration struct {
	Id        int
	Debug     bool
	Dev       bool
	Anthive   string
	Templates string
	LXC       string
}

func (c *Configuration) Save() (err error) {
	var (
		file     *os.File
		filepath string
		tmpl     *template.Template
	)
	if viper.GetBool("Dev") {
		filepath = CONFIG_NAME // create file in place
	} else {
		filepath = path.Join("/etc", CONFIG_NAME)
	}
	if file, err = os.Create(filepath); err != nil {
		return
	}
	defer file.Close()
	filepath = path.Join(c.Templates, CONFIG_NAME)
	if tmpl, err = template.ParseFiles(filepath); err != nil {
		return
	}
	return tmpl.Execute(file, c)
}

func PreRun(cmd *cobra.Command, args []string) {
	// reinit args for glog
	os.Args = os.Args[:1]

	// make glog happy and log to the correct place
	flag.Set("logtostderr", "true")
	flag.Parse()

	viper.ReadInConfig()
	// check dev mode, and reset some configs
	if viper.GetBool("Dev") {
		glog.Infof("dev mode enabled")
		viper.Set("Templates", path.Join(".", "templates"))
	}

	viper.Unmarshal(Config) // this will load default config

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			glog.Infof("creating a new config file")
			// a bit crappy but should handle this correctly
			RegisterCommand.Run(cmd, args)
			viper.ReadInConfig()
		} else {
			glog.Fatalf("when reading config file: %s", err)
		}
	}

	if err := viper.Unmarshal(Config); err != nil {
		glog.Errorf("%q", Config)
		glog.Fatalf("when unmarshalling the config: %s", err)
	}
	if Config.Debug { // debug is same as -vvvvv
		verbosity = 5
	}
	flag.Set("v", strconv.Itoa(verbosity))
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
	devMsg := "dev mode, instead of getting path from the system use those at $PWD"
	verboseMsg := "verbose output, can be stacked to increase verbosity"

	Command.PersistentFlags().Bool("debug", false, debugMsg)
	Command.PersistentFlags().Bool("dev", false, devMsg)
	Command.PersistentFlags().CountVarP(&verbosity, "verbose", "v", verboseMsg)

	Command.AddCommand(RegisterCommand)

	viper.BindPFlag("Debug", Command.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("Dev", Command.PersistentFlags().Lookup("dev"))

	viper.Set("PROJECT", "github.com/alienantfarm/antling")

	// set some paths
	viper.Set("LXC", path.Join(sep, "var", "lib", "lxc"))
	viper.Set("Templates", path.Join(sep, "usr", "share", "antling", "templates"))

	viper.BindEnv("Anthive", "ANTHIVE_URL")

	viper.SetConfigName(CONFIG_NAME[:len(CONFIG_NAME)-5])

	viper.AddConfigPath("/etc")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$ANTLING_CONFIG/")
}
