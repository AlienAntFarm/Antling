package scheduler

import (
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/golang/glog"
)

type scheduler struct {
	jobs    map[int]*structs.Job
	channel chan *structs.Job
}

var Scheduler *scheduler

func InitScheduler() *scheduler {
	if Scheduler != nil {
		glog.Fatalf("scheduler already inited, something bad is happening")
	} else {
		Scheduler = &scheduler{
			make(map[int]*structs.Job),
			make(chan *structs.Job, 1),
		}
		go Scheduler.start()
	}
	return Scheduler
}

func (s *scheduler) start() {
	for job := range s.channel {
		if _, ok := s.jobs[job.Id]; !ok { // async stuff going around, avoid restart something existing
			s.jobs[job.Id] = job
			glog.Infof("starting job %d", job.Id)
			// start the job
		}
	}
}

func (s *scheduler) ProcessJob(job *structs.Job) {
	s.channel <- job
}
