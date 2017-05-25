package main

import (
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/alienantfarm/antling/client"
	"github.com/alienantfarm/antling/scheduler"
	"github.com/alienantfarm/antling/utils"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"time"
)

func run(cmd *cobra.Command, args []string) {
	s := scheduler.InitScheduler()

	// setup self
	self := client.NewClient().Antling

	// start main loop
	for {
		// retrieve jobs from server
		jobs, err := self.GetJobs()
		if err != nil {
			glog.Errorf("%s", err)
		}

		// go over jobs and start new jobs
		for _, job := range jobs {
			if job.State == structs.JOB_NEW {
				job.State += 1
			}
			s.ProcessJob(job)
		}

		// now update the server so it nows which jobs have been started
		err = self.Update()
		if err != nil {
			glog.Errorf("%s", err)
		}

		time.Sleep(10 * time.Second)
	}
}

func runRegister(cmd *cobra.Command, args []string) {
	antling := client.NewClient().Antling
	if err := antling.Create(); err != nil {
		glog.Fatalf("%s", err)
	}
	utils.Config.Id = antling.Id
	if err := utils.Config.Save(); err != nil {
		glog.Fatalf("%s", err)
	}
}

func main() {
	utils.Command.Run = run
	utils.RegisterCommand.Run = runRegister
	if err := utils.Command.Execute(); err != nil {
		glog.Fatalf("%s", err)
	}
}
