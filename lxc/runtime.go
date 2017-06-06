package lxc

import (
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/alienantfarm/antling/utils"
	"github.com/golang/glog"
	"gopkg.in/lxc/go-lxc.v2"
	"os"
	"path"
	"runtime"
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

	options, err := setupOptions(job, p)
	if err != nil {
		return err
	}
	if err := c.Start(); err != nil {
		return err
	}

	glog.Infof("RunCommand")
	if _, err := c.RunCommand(job.SanitizeCmd(), *options); err != nil {
		return err
	}

	glog.Infof("Done")
	if err := c.Stop(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Hostname string
	Env      []string
	RootFS   string
	Arch     string
}

func setupOptions(job *structs.Job, pipes *pipes) (*lxc.AttachOptions, error) {

	config := &Config{
		Env:      job.SanitizeEnv(),
		Hostname: job.Image.Hostname,
		RootFS:   path.Join(utils.Config.Paths.LXC, job.Image.Archive),
		Arch:     runtime.GOARCH,
	}
	glog.Infof(utils.MarshalJSON(config))

	// create file and write it through a template
	if file, err := os.Create(path.Join(config.RootFS, "config")); err != nil {
		return nil, err
	} else if err := utils.Config.Templates.ConfLXC.Execute(file, config); err != nil {
		return nil, err
	} else {
		file.Close()
	}

	glog.Infof(job.SanitizeCwd())
	return &lxc.AttachOptions{
		Namespaces: -1,
		Arch:       -1,
		Cwd:        job.SanitizeCwd(),
		UID:        -1,
		GID:        -1,
		ClearEnv:   false,
		Env:        nil,
		EnvToKeep:  nil,
		StdinFd:    pipes.stdin.Fd(),
		StdoutFd:   pipes.stdout.Fd(),
		StderrFd:   pipes.stderr.Fd(),
	}, nil
}
