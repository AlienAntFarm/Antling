package main

import (
	"flag"
	"github.com/alienantfarm/antling/client"
	"github.com/alienantfarm/antling/scheduler"
	"github.com/alienantfarm/antling/utils/config"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"os"
	"time"
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
		c := config.Get()
		s := scheduler.InitScheduler()
		if c.Debug {
			flag.Set("v", "10") // totally arbitrary but who cares!
			flag.Parse()
		}
		glog.V(1).Infoln("debug mode enabled")

		// setup self
		self := client.NewAntling(c.Id, c.Client)

		// start main loop
		for {
			jobs, err := self.GetJobs()
			s.ProcessJobs(jobs)
			if err != nil {
				glog.Fatalf("%s", err)
			}
			time.Sleep(10 * time.Second)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		glog.Fatalf("%s", err)
	}
}
