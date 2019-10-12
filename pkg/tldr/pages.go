package tldr

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	pageSource      = "https://tldr.sh/assets/tldr.zip"
	indexSource     = "https://tldr.sh/assets/index.json"
	cacheExpiredErr = "more than a week passed, should update tldr using --update"
	// CacheExpiredMsg a message to tell user should update
	CacheExpiredMsg = cacheExpiredErr
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
}

// NewTldr create a instance of tldr repository
func NewTldr(tldrPath string, op Options) *Tldr {
	return &Tldr{
		path:           tldrPath,
		pageSourceURL:  pageSource,
		indexSourceURL: indexSource,
		indexFile:      filepath.Base(indexSource),
		zipFile:        filepath.Base(pageSource),
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
	if !pathExists(t.path) {
		if err := os.Mkdir(t.path, 0755); err != nil {
			return err
		}
	}

	if t.update {
		if err := t.Update(); err != nil {
			return err
		}
	}

	cache := NewCache(t.path, t.indexFile, 24*7*time.Hour)
	if cache.Expired() {
		return fmt.Errorf(cacheExpiredErr)
	}

	return nil
}

// IsCacheExpired return true if `err` means expired cache
func IsCacheExpired(err error) bool {
	if err != nil && err.Error() == cacheExpiredErr {
		return true
	}
	return false
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

// FindPage find tldr page by `cmd`
func (t *Tldr) FindPage(cmd string) (*Page, error) {
	for _, pt := range t.platformDirs {
		path := filepath.Join(t.path, t.langDir, pt, cmd+".md")
		if !pathExists(path) {
			// try to find it in next platform if cmd does not exists
			continue
		}

		f, err := os.Open(path)
		if err != nil {
			return &Page{}, err
		}
		defer f.Close()

		return parsePage(f)
	}

	return &Page{}, fmt.Errorf("not found %s page", cmd)
}
