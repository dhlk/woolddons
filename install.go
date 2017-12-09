package main

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func ZipDownload(url, target string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func zipAcross(source string, fn func(file *zip.File) error) error {
	dat, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	zr, err := zip.NewReader(bytes.NewReader(dat), int64(len(dat)))
	if err != nil {
		return err
	}

	for _, file := range zr.File {
		err = fn(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func modZipFile(target string, install bool) func(file *zip.File) error {
	return func(file *zip.File) error {
		dest := target + file.Name

		if !install {
			return os.Remove(dest)
		}

		zfr, err := file.Open()
		if err != nil {
			return err
		}
		os.MkdirAll(path.Dir(dest), os.ModePerm)
		d, _ := os.Create(dest)
		io.Copy(d, zfr)
		d.Close()

		return nil
	}
}

func ZipInstall(source, target string) error {
	return zipAcross(source, modZipFile(target, true))
}

func ZipUninstall(source, target string) error {
	// TODO cleanup directories?
	return zipAcross(source, modZipFile(target, false))
}
