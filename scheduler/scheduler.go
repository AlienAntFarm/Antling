package scheduler

import (
	"github.com/alienantfarm/antling/client"
	"github.com/golang/glog"
)

type scheduler struct{}

var s *scheduler

func Get() *scheduler {
	if s == nil {
		s = &scheduler{}
	}
	return s
}

func (s *scheduler) ProcessJobs(jobs []*client.Job) {
	glog.Infof("processing jobs %q", jobs)
}
