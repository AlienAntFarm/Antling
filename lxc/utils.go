package lxc

import (
	"bufio"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

type pipes struct {
	stdin, stderr, stdout *os.File
	dir                   string
}

func pipeReading(path string) {
	glog.Infof("opening %s for reading", path)
	if pipe, err := os.OpenFile(path, os.O_RDONLY, 0600); err != nil {
		glog.Errorf("opening %s failed, %s", path, err)
	} else {
		defer pipe.Close()
		scanner := bufio.NewScanner(pipe)
		for scanner.Scan() {
			glog.Infoln(scanner.Text()) // Println will add back the final '\n'
		}
		glog.Infof("closing reading scanner on %s", path)
	}
}

func initPipe(path string) (*os.File, error) {
	if err := syscall.Mkfifo(path, 0600); err != nil {
		return nil, err
	}
	go pipeReading(path)
	glog.Infof("opening for %s writing", path)
	return os.OpenFile(path, os.O_WRONLY, 0600)
}

func (p *pipes) Close() {
	for _, pipe := range []*os.File{p.stdin, p.stderr, p.stdout} {
		pipe.Close()
	}
	os.RemoveAll(p.dir)
}

func newPipes(prefix string) (p *pipes, err error) {
	p = &pipes{}

	if p.stdin, err = os.Open(os.DevNull); err != nil {
		return
	}
	if p.dir, err = ioutil.TempDir("", "antling."+prefix); err != nil {
		return
	}
	defer utils.RemoveOnFail(p.dir, err)

	if p.stderr, err = initPipe(path.Join(p.dir, "stderr")); err != nil {
		return
	}
	if p.stdout, err = initPipe(path.Join(p.dir, "stdout")); err != nil {
		return
	}

	// set up a named pipe for logs also
	if _, err = initPipe(path.Join(p.dir, "logs")); err != nil {
		return
	}

	return
}
