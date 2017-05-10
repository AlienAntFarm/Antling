package scheduler

import (
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/golang/glog"
)

type scheduler struct {
	jobs    []*structs.Job
	channel chan *structs.Job
}

var Scheduler *scheduler

func InitScheduler() *scheduler {
	if Scheduler != nil {
		glog.Fatalf("scheduler already inited, something bad is happening")
	} else {
		Scheduler = &scheduler{[]*structs.Job{}, make(chan *structs.Job, 1)}
		go Scheduler.start()
	}
	return Scheduler
}

func (s *scheduler) start() {
	for job := range s.channel {
		glog.Infof("processing job %d, with status %s", job.Id, structs.JOB_STATES[job.State])
	}
}

func (s *scheduler) ProcessJobs(jobs []*structs.Job) {
	for _, job := range jobs {
		s.channel <- job
	}
}
