package tldr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	pageSourceURL  = "https://tldr.sh/assets/tldr.zip"
	indexSourceURL = "https://tldr.sh/assets/index.json"
)

// Options are tldr functions
type Options struct {
	Platform string
	Language string
	Update   bool
}

// Tldr Repository of tldir pages
type Tldr struct {
	path           string
	pageSourceURL  string
	indexSourceURL string
	indexFile      string
	zipFile        string
	platformDirs   []string
	langDir        string
	update         bool
	cacheMaxAge    time.Duration
}

// New create a instance of tldr repository
func New(tldrPath string, op Options) *Tldr {
	return &Tldr{
		path:           tldrPath,
		pageSourceURL:  pageSourceURL,
		indexSourceURL: indexSourceURL,
		indexFile:      filepath.Base(indexSourceURL),
		zipFile:        filepath.Base(pageSourceURL),
		platformDirs:   []string{op.Platform, "common"},
		langDir:        convertToLangDir(op.Language),
		update:         op.Update,
	}
}

func convertToLangDir(lang string) string {
	// TODO return multi language dirs Depending on lang
	return "pages"
}

// OnInitialize create tldr directory and check tldr pages whether cache expired or not
func (t *Tldr) OnInitialize() error {
	initUpdate := false
	if !pathExists(t.path) {
		if err := os.Mkdir(t.path, 0755); err != nil {
			return err
		}
		// automatically updated if indexfile does not exist
		initUpdate = true
	}

	if t.update || initUpdate {
		if err := t.Update(); err != nil {
			return err
		}
	}

	return nil
}

// Update tldr pages from remote zip file
func (t *Tldr) Update() error {
	_, err := download(t.indexSourceURL, t.path, t.indexFile)
	if err != nil {
		return err
	}

	zipPath, err := download(t.pageSourceURL, t.path, t.zipFile)
	if err != nil {
		return err
	}

	if err := unzip(zipPath, t.path); err != nil {
		return err
	}

	return nil
}

// FindPage find tldr page by `cmds`
func (t *Tldr) FindPage(cmds []string) (*Page, error) {
	for _, pt := range t.platformDirs {
		page := strings.Join(cmds, "-") + ".md"
		path := filepath.Join(t.path, t.langDir, pt, page)
		if !pathExists(path) {
			// if cmd does not exist, try to find it in next platform
			continue
		}

		f, err := os.Open(path)
		if err != nil {
			return &Page{}, err
		}
		defer f.Close()

		return parsePage(f)
	}

	return &Page{}, fmt.Errorf("not found %s page", cmds)
}

// Expired return true if cache is expired
func (t *Tldr) Expired(maxAge time.Duration) bool {
	indexPath := filepath.Join(t.path, t.indexFile)
	age, err := age(indexPath)
	if err != nil {
		return true
	}

	return age > maxAge
}

// age return the time since the data is cached at
func age(path string) (time.Duration, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return time.Duration(0), err
	}

	return time.Since(fi.ModTime()), nil
}

// pathExists return true if path exists
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
