package lxc

import (
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/alienantfarm/antling/utils"
	"github.com/golang/glog"
	"gopkg.in/lxc/go-lxc.v2"
	"io/ioutil"
	"os"
)

type pipes struct {
	stdin, stderr, stdout, logs *os.File
}

func newPipes() (p *pipes, err error) {
	p = &pipes{}

	if p.stdin, err = os.Open(os.DevNull); err != nil {
		return
	}
	for _, pipe := range []**os.File{&p.stderr, &p.stdout, &p.logs} {
		if *pipe, err = ioutil.TempFile("", ""); err != nil {
			return
		}
	}
	return
}

func (p *pipes) Close() {
	p.stdin.Close()
	for _, pipe := range []*os.File{p.stderr, p.stdout, p.logs} {
		pipe.Close()
		if err := os.Remove(pipe.Name()); err != nil {
			glog.Errorf("deleting custom pipe failed %s", err)
		}
	}
}

func (p *pipes) Log() error {
	for i, pipe := range []*os.File{p.stderr, p.stdout, p.logs} {
		pipe.Seek(0, os.SEEK_SET)
		if b, err := ioutil.ReadAll(pipe); err != nil {
			return err
		} else {
			glog.Infof("%s:\n %s", []string{"stderr", "stdout", "logs"}[i], b)
		}
	}
	return nil
}

func Start(job *structs.Job) error {
	// open temporary files for default pipes
	p, err := newPipes()
	if err != nil {
		glog.Infof("opening temp files for container output failed")
		return err
	}
	defer p.Close()

	//create and start the container
	c, err := lxc.NewContainer(job.Image.Archive, utils.Config.Paths.LXC)
	if err != nil {
		return err
	}
	c.SetLogLevel(LOG_LEVELS[utils.LogLevel()]) // pass log level to container
	c.SetLogFile(p.logs.Name())
	defer p.Log()
	options := lxc.DefaultAttachOptions
	options.StdinFd = p.stdin.Fd()
	options.StdoutFd = p.stdout.Fd()
	options.StderrFd = p.stderr.Fd()

	if err := c.Start(); err != nil {
		return err
	}

	glog.Infof("RunCommand")
	if _, err := c.RunCommand([]string{"ls"}, options); err != nil {
		return err
	}

	glog.Infof("Done")
	if err := c.Stop(); err != nil {
		return err
	}
	return nil
}
