package lxc

import (
	"bufio"
	"bytes"
	"github.com/alienantfarm/antling/utils"
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

func pipeReading(path string, splitFunc bufio.SplitFunc) {
	glog.Infof("opening %s for reading", path)
	if pipe, err := os.OpenFile(path, os.O_RDONLY, 0600); err != nil {
		glog.Errorf("opening %s failed, %s", path, err)
	} else {
		defer pipe.Close()
		scanner := bufio.NewScanner(pipe)
		if splitFunc == nil {
			splitFunc = bufio.ScanLines
		}
		scanner.Split(splitFunc)
		for scanner.Scan() {
			glog.Infoln(scanner.Text()) // Println will add back the final '\n'
		}
		glog.Infof("closing reading scanner on %s", path)
	}
}

func initPipe(path string, splitFunc bufio.SplitFunc) (*os.File, error) {
	if err := syscall.Mkfifo(path, 0600); err != nil {
		return nil, err
	}
	go pipeReading(path, splitFunc)
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

	if p.stderr, err = initPipe(path.Join(p.dir, "stderr"), nil); err != nil {
		return
	}
	if p.stdout, err = initPipe(path.Join(p.dir, "stdout"), nil); err != nil {
		return
	}

	// set up a named pipe for logs also
	if _, err = initPipe(path.Join(p.dir, "logs"), ScanLXCLogs); err != nil {
		return
	}

	return
}
