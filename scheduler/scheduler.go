package scheduler

import (
	"github.com/alienantfarm/antling/client"
	"github.com/golang/glog"
)

type scheduler struct {
	jobs    []*client.Job
	channel chan *client.Job
}

var Scheduler *scheduler

func InitScheduler() *scheduler {
	if Scheduler != nil {
		glog.Fatalf("scheduler already inited, something bad is happening")
	} else {
		Scheduler = &scheduler{[]*client.Job{}, make(chan *client.Job, 1)}
		go Scheduler.start()
	}
	return Scheduler
}

func (s *scheduler) start() {
	for job := range s.channel {
		glog.Infof("processing job %d, with status %s", job.Id, client.JOB_STATES[job.State])
	}
}

func (s *scheduler) ProcessJobs(jobs []*client.Job) {
	for _, job := range jobs {
		s.channel <- job
	}
}
