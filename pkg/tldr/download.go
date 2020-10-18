package tldr

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// download data from `url` to `dstDir` as `filename`
func download(url, dstDir, filename string) (string, error) {
	path := filepath.Join(dstDir, filename)
	if pathExists(path) {
		if err := os.RemoveAll(path); err != nil {
			return "", err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if _, err := io.Copy(f, res.Body); err != nil {
		return "", err
	}

	return path, nil
}

// unzip to `dstDir`
func unzip(zipPath, dstDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		if f.FileInfo().IsDir() {
			path := filepath.Join(dstDir, f.Name)
			if err = os.MkdirAll(path, f.Mode()); err != nil {
				return err
			}
		} else {
			buf := make([]byte, f.UncompressedSize)
			_, err = io.ReadFull(rc, buf)
			if err != nil {
				return err
			}

			path := filepath.Join(dstDir, f.Name)
			if err = ioutil.WriteFile(path, buf, f.Mode()); err != nil {
				return err
			}
		}
	}

	return nil
}
