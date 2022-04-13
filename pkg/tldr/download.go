package tldr

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// download data from `url` to `dstDir` as `filename`
func download(ctx context.Context, url, dstDir, filename string) (_ string, reterr error) {
	path := filepath.Join(dstDir, filename)
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if ferr := f.Close(); ferr != nil && reterr == nil {
			reterr = ferr
		}
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; code != http.StatusOK {
		return "", fmt.Errorf("http response code was %d for downloading from %s", code, url)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}

	return path, nil
}

// unzip to `dstDir`
func unzip(ctx context.Context, zipPath, dstDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	if len(r.File) == 0 {
		return errors.New("no files in a zip")
	}

	for _, f := range r.File {
		fn := func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				break
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			if f.FileInfo().IsDir() {
				path := filepath.Join(dstDir, f.Name)
				if err := os.MkdirAll(path, f.Mode()); err != nil {
					return err
				}
			} else {
				buf := make([]byte, f.UncompressedSize)
				_, err := io.ReadFull(rc, buf)
				if err != nil {
					return err
				}

				path := filepath.Join(dstDir, f.Name)
				if err := os.WriteFile(path, buf, f.Mode()); err != nil {
					return err
				}
			}
			return nil
		}

		if reterr := fn(); reterr != nil {
			return reterr
		}
	}

	return nil
}
