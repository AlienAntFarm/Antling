package lxc

import (
	"archive/tar"
	"github.com/golang/glog"
	"io"
	"os"
	"path"
)

func DeflateLXC(reader io.ReadCloser) error {
	defer reader.Close()
	tarReader := tar.NewReader(reader)
	flag := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		info := header.FileInfo()
		if info.Name() == ".dockerenv" { // ignore .dockerenv file
			continue
		}

		if glog.V(5) { // only for debug, output too much informations
			glog.Infof("Uncompressing %s, with metadata %q", header.Name, info)
		}
		target := path.Join(".", header.Name)
		mode := info.Mode()
		switch {
		case mode&os.ModeDir != 0:
			if err = os.MkdirAll(target, mode); err != nil {
				return err
			}
		case mode&os.ModeSymlink != 0:
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		default:
			if file, err := os.OpenFile(target, flag, mode); err != nil {
				return err
			} else if _, err = io.Copy(file, tarReader); err != nil {
				return err
			} else if err := file.Close(); err != nil {
				return err
			}
		}

	}
	return nil
}
