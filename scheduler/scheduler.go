package scheduler

import (
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/alienantfarm/antling/client"
	"github.com/alienantfarm/antling/utils"
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type scheduler struct {
	jobs    map[int]*structs.Job
	channel chan *structs.Job
	cache   []string
	client  *client.Client
}

var Scheduler *scheduler

func InitScheduler() *scheduler {
	if Scheduler != nil {
		glog.Fatalf("scheduler already inited, something bad is happening")
	} else {
		glog.Infoln("init scheduler")

		Scheduler = &scheduler{
			make(map[int]*structs.Job),
			make(chan *structs.Job, 1),
			[]string{},
			client.NewClient(),
		}
		// retrieve cache images
		files, err := ioutil.ReadDir(utils.Config.Cache)
		if err != nil {
			glog.Fatalf("reading images cache failed %s", err)
		}
		for _, image := range files {
			glog.Infof("caching image %s", image.Name())
			Scheduler.cache = append(
				Scheduler.cache, path.Join(utils.Config.Cache, image.Name()),
			)
		}
		go Scheduler.start()
	}
	return Scheduler
}

func (s *scheduler) start() {
	for job := range s.channel {
		if _, ok := s.jobs[job.Id]; !ok { // async stuff going around, avoid restart something existing
			glog.Infof("starting job %d", job.Id)
			// start the job
			if err := s.startJob(job); err != nil {
				glog.Infof("job starting failed, retrying later")
				continue
			}
			s.jobs[job.Id] = job
		}
	}
}

func (s *scheduler) checkCache(image string) (entry string, err error) {
	glog.Infof("looking for %s in cache", image)
	for _, entry = range s.cache {
		if path.Base(entry) == image {
			return
		}
	}
	entry = path.Join(utils.Config.Cache, image)
	writer, err := os.Create(entry)
	if err != nil {
		return
	}
	reader, err := s.client.Images.Get(image)
	if err != nil {
		return
	}
	_, err = io.Copy(writer, reader)
	if err != nil {
		return
	}
	s.cache = append(s.cache, entry)
	return
}

func (s *scheduler) startJob(job *structs.Job) error {
	// check cache and download if image is not here
	_, err := s.checkCache(job.Image.Archive)
	if err != nil {
		glog.Errorf("cache hit failed %s", err)
		return err
	}

	return nil
}

func (s *scheduler) ProcessJob(job *structs.Job) {
	s.channel <- job
}
