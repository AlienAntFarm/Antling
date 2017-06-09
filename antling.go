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

		for _, job := range jobs { // go over jobs and infer what to do
			if job.State == structs.JOB_NEW {
				job.State += 1
			}
			if _, ok := self.Jobs[job.Id]; !ok {
				self.Jobs[job.Id] = job
			}
			s.ProcessJob(job)
		}

		// now update the server so it knows which jobs have been started, finished or errored
		if err := self.Update(); err != nil {
			glog.Errorf("%s", err)
		} else { // subject to race conditions !
			for _, job := range self.Jobs {
				if job.State > structs.JOB_PENDING {
					delete(self.Jobs, job.Id)
				}
			}
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
