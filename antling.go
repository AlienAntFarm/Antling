package main

import (
	"flag"
	"github.com/alienantfarm/antling/utils"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "antling",
	Short: "Start the antling process",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// reinit args for glog
		os.Args = os.Args[:1]
		flag.Set("logtostderr", "true")
		flag.Parse()
	},
	Run: func(cmd *cobra.Command, args []string) {
		conf := utils.Config()
		if conf.Debug {
			flag.Set("v", "10") // totally arbitrary but who cares!
			flag.Parse()
		}
		glog.V(1).Infoln("debug mode enabled")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		glog.Fatalf("%s", err)
	}
}
