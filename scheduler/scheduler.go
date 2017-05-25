package scheduler

import (
	"archive/tar"
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
		// retrieve lxc images
		files, err := ioutil.ReadDir(utils.Config.LXC)
		if err != nil {
			glog.Fatalf("reading images cache failed %s", err)
		}
		for _, image := range files {
			glog.Infof("caching image %s", image.Name())
			Scheduler.cache = append(
				Scheduler.cache, path.Join(utils.Config.LXC, image.Name()),
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

func (s *scheduler) deflateLXC(reader *tar.Reader) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		info := header.FileInfo()
		target := path.Join(pwd, header.Name)
		if info.IsDir() {
			if err = os.MkdirAll(target, info.Mode()); err != nil {
				return err
			}
			continue
		}
		// do not deflate manifest.json for now
		if info.Name() == "manifest.json" {
			continue
		}

		file, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, reader)
		if err != nil {
			return err
		}
	}
	// now add to in memory cache
	s.cache = append(s.cache, path.Dir(pwd)) // this remove rootfs from path
	return nil
}

func (s *scheduler) checkLXC(image string) error {
	glog.Infof("looking for %s in cache", image)

	for _, entry := range s.cache {
		if path.Base(entry) == image {
			return nil
		}
	}
	rootfs := path.Join(utils.Config.LXC, image, "rootfs")

	if err := os.MkdirAll(rootfs, 0770); err != nil {
		return err
	}
	if err := os.Chdir(rootfs); err != nil {
		return err
	}

	if reader, err := s.client.Images.Get(image); err != nil {
		return err
	} else {
		glog.Infof("caching image: %s", image)
		defer reader.Close()
		return s.deflateLXC(tar.NewReader(reader))
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
