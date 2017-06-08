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
			if job.State > structs.JOB_PENDING {
				continue
			}
			if job.Retries > structs.MAX_RETRIES {
				glog.Errorf("too much retries for job %d, entering failed state", job.Id)
				job.State = structs.JOB_ERROR
				continue
			}
			glog.Infof("starting job %d", job.Id)
			s.jobs[job.Id] = job

			// start the job
			go func(job *structs.Job) {
				if err := s.startJob(job); err != nil {
					glog.Infof("job starting failed, retrying later")
					job.Retries += 1
					delete(s.jobs, job.Id)
				} else {
					glog.Infof("job %d succeed", job.Id)
					job.State += 1
					glog.Infof(utils.MarshalJSON(job))
				}
			}(job) // not sure how this behave in asynchronous world, better copy address
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

	if err = os.MkdirAll(rootfs, 0755); err != nil {
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
		glog.Errorf("cache hit failed for %s: %s", job.Image.Archive, err)
		return err
	}
	if err := lxc.Start(job); err != nil {
		glog.Errorf("error when running container %s: %s", job.Image.Archive, err)
		return err
	}

	return nil
}

func (s *scheduler) ProcessJob(job *structs.Job) {
	s.channel <- job
}
