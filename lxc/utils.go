package lxc

import (
	"bufio"
	"bytes"
	"github.com/alienantfarm/antling/utils"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"syscall"
)

type pipes struct {
	stdin, stderr, stdout *os.File
	dir                   string
	closing               bool
	wg                    *sync.WaitGroup
}

func extractLog(data []byte) []byte {
	splitted := bytes.Split(data, []byte{' ', ' ', ' ', ' '})
	data = (splitted[len(splitted)-1])
	data = bytes.TrimPrefix(data, []byte{' '})
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// reuse https://golang.org/src/bufio/scan.go?s=11488:11566#L329
func ScanLXCLogs(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, extractLog(data[0:i]), nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), extractLog(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func (p *pipes) pipeReading(path string, splitFunc bufio.SplitFunc) {
	defer p.wg.Done()

	glog.Infof("opening %s for reading", path)
	if pipe, err := os.OpenFile(path, os.O_RDONLY, os.ModeNamedPipe); err != nil {
		glog.Errorf("opening %s failed, %s", path, err)
	} else {
		defer pipe.Close()
		scanner := bufio.NewScanner(pipe)
		scanner.Split(splitFunc)
		for scanner.Scan() {
			if t := scanner.Text(); t != "" {
				glog.Infoln(t)
			}
			if p.closing { // when  finished state is reached switch to non-blocking
				syscall.SetNonblock(int(pipe.Fd()), true)
			}
		}
		glog.Infof("closing reading scanner on %s", path)
	}
}

func (p *pipes) initPipe(file string, splitFunc bufio.SplitFunc) (*os.File, error) {
	file = path.Join(p.dir, file)
	if err := syscall.Mkfifo(file, 0600); err != nil {
		return nil, err
	}
	p.wg.Add(1)
	go p.pipeReading(file, splitFunc)
	glog.Infof("opening for %s writing", file)
	return os.OpenFile(file, os.O_WRONLY, os.ModeNamedPipe)
}

func (p *pipes) Close() {
	p.closing = true
	for _, pipe := range []*os.File{p.stdin, p.stdout, p.stderr} {
		pipe.Write([]byte{'\n'})
		pipe.Close()
	}
	p.wg.Wait()
	glog.Infof("all pipes are closed, removing temp dir %s", p.dir)
	if err := os.RemoveAll(p.dir); err != nil {
		glog.Errorf("%s", err)
	}
}

func newPipes(prefix string) (p *pipes, err error) {
	p = &pipes{wg: &sync.WaitGroup{}}

	if p.dir, err = ioutil.TempDir("", "antling."+prefix); err != nil {
		return
	}
	defer utils.RemoveOnFail(p.dir, err)

	if p.stdin, err = os.Open(os.DevNull); err != nil {
		return
	}
	if p.stdout, err = p.initPipe("stdout", bufio.ScanLines); err != nil {
		return
	}
	if p.stderr, err = p.initPipe("stderr", bufio.ScanLines); err != nil {
		return
	}
	// set up a named pipe for logs also
	if _, err = p.initPipe("logs", ScanLXCLogs); err != nil {
		return
	}

	return
}
