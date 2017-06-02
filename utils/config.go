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

const CONFIG = "antling.toml"
const CONFIG_LXC = "lxc.conf"
const IMAGES_PREFIX = "/static/images"

type Configuration struct {
	Id      int
	Debug   bool
	Dev     bool
	Anthive string
	Paths   struct {
		Templates string
		LXC       string
		Conf      string
	}
	Templates struct {
		Path    string
		ConfLXC *template.Template
		Conf    *template.Template
	}
}

func (c *Configuration) Save() error {
	if file, err := os.Create(c.Paths.Conf); err != nil {
		return err
	} else {
		defer file.Close()
		return c.Templates.Conf.Execute(file, c)
	}
}

func (c *Configuration) LoadTemplates() error {
	tplts := map[string]**template.Template{
		CONFIG_LXC: &c.Templates.ConfLXC,
		CONFIG:     &c.Templates.Conf,
	}
	for name, tplt := range tplts {
		name = path.Join(c.Paths.Templates, name)
		if t, err := template.ParseFiles(name); err != nil {
			return err
		} else {
			*tplt = t
		}
	}
	return nil
}

func (c *Configuration) CreatePaths() error {
	paths := []string{
		c.Paths.LXC,
	}

	for _, p := range paths {
		if err := os.MkdirAll(p, 0755); err != nil {
			return err
		}
	}
	return nil
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
		wd, err := os.Getwd()
		if err != nil {
			glog.Fatalf("%s", err)
		}
		viper.Set("Paths.Templates", path.Join(wd, "templates"))
		viper.Set("Paths.Conf", path.Join(CONFIG)) // create config in place
		viper.Set("Paths.LXC", path.Join(wd, "var", "antling"))
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

	if err := Config.LoadTemplates(); err != nil {
		glog.Fatalf("%s", err)
	}
	if err := Config.CreatePaths(); err != nil {
		glog.Fatalf("%s", err)
	}

	glog.V(1).Infoln("debug mode enabled")
	glog.V(2).Infof("%q", Config)
}

func LogLevel() int {
	//discard error
	level, _ := strconv.Atoi(flag.Lookup("v").Value.String())
	return level
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
	viper.Set("Paths.LXC", path.Join(sep, "var", "lib", "lxc"))
	viper.Set("Paths.Templates", path.Join(sep, "usr", "share", "antling", "templates"))
	viper.Set("Paths.Conf", path.Join("/etc", CONFIG))

	viper.BindEnv("Anthive", "ANTHIVE_URL")

	viper.SetConfigName(CONFIG[:len(CONFIG)-5])

	viper.AddConfigPath("/etc")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$ANTLING_CONFIG/")
}
