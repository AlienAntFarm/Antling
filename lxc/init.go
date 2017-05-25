package lxc

import (
	"archive/tar"
	"io"
	"os"
	"path"
)

func DeflateLXC(reader io.ReadCloser) error {
	defer reader.Close()
	tarReader := tar.NewReader(reader)
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
		if info.Name() == "manifest.json" { // generate the config
			if err := deflateConfig(tarReader); err != nil {
				return err
			}
			continue
		}

		target := path.Join(".", header.Name)
		if info.IsDir() {
			if err = os.MkdirAll(target, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func deflateConfig(reader io.Reader) error {
	return nil
}
