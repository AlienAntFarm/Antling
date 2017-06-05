package lxc

import (
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/alienantfarm/antling/utils"
	"github.com/golang/glog"
	"gopkg.in/lxc/go-lxc.v2"
	"path"
)

var LOG_LEVELS = [...]lxc.LogLevel{
	lxc.CRIT,
	lxc.ERROR,
	lxc.WARN,
	lxc.NOTICE,
	lxc.INFO,
	lxc.TRACE,
}

func Start(job *structs.Job) error {
	// open temporary files for default pipes
	p, err := newPipes(job.Image.Archive)
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
	c.SetLogFile(path.Join(p.dir, "logs"))

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
