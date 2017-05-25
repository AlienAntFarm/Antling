package scheduler

import (
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/alienantfarm/antling/client"
	"github.com/alienantfarm/antling/lxc"
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
		// retrieve lxc images
		files, err := ioutil.ReadDir(utils.Config.Paths.LXC)
		if err != nil {
			glog.Fatalf("reading images cache failed %s", err)
		}
		for _, image := range files {
			glog.Infof("caching image %s", image.Name())
			Scheduler.cache = append(Scheduler.cache, image.Name())
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

func (s *scheduler) checkLXC(image string) (err error) {
	var reader io.ReadCloser

	glog.Infof("looking for %s in cache", image)

	for _, entry := range s.cache {
		if entry == image {
			return
		}
	}
	rootfs := path.Join(utils.Config.Paths.LXC, image, "rootfs")

	if err = os.MkdirAll(rootfs, 0770); err != nil {
		return
	}
	defer func() { utils.RemoveOnFail(path.Dir(rootfs), err) }()

	if err = os.Chdir(rootfs); err != nil {
		return
	}

	glog.Infof("caching image: %s", image)
	if reader, err = s.client.Images.Get(image); err != nil {
		return
	} else if err = lxc.DeflateLXC(reader); err != nil {
		return
	} else {
		s.cache = append(s.cache, image)
		return
	}
}

func (s *scheduler) startJob(job *structs.Job) error {
	// check lxc dir and download if image is not here
	if err := s.checkLXC(job.Image.Archive); err != nil {
		glog.Errorf("cache hit failed %s", err)
		return err
	}

	return nil
}

func (s *scheduler) ProcessJob(job *structs.Job) {
	s.channel <- job
}
