package tldr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	pageSourceURL  = "https://tldr.sh/assets/tldr.zip"
	languageCodeEN = "en"
)

// Options are tldr functions
type Options struct {
	Platform string
	Language string
	Update   bool
}

// Tldr Repository of tldir pages
type Tldr struct {
	path          string
	pageSourceURL string
	platforms     []string
	languages     []string
	update        bool
}

var ErrNotFoundPage = errors.New("no page found")

// New create a instance of tldr repository
func New(tldrPath string, opt *Options) *Tldr {
	if opt == nil {
		opt = new(Options)
	}

	return &Tldr{
		path:          tldrPath,
		pageSourceURL: pageSourceURL,
		platforms:     []string{opt.Platform, "common"},
		languages:     getLanguages(opt.Language),
		update:        opt.Update,
	}
}

// OnInitialize create and update tldr directory
func (t *Tldr) OnInitialize() error {
	initUpdate := false
	if !pathExists(t.path) {
		if err := os.MkdirAll(t.path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create tldr dir: %w", err)
		}
		// automatically updated if indexfile does not exist
		initUpdate = true
	}

	if t.update || initUpdate {
		if err := t.Update(); err != nil {
			return fmt.Errorf("failed to update tldr repository: %w", err)
		}
	}

	if f := t.indexFilePath(); !pathExists(f) {
		return fmt.Errorf("tldr database is broken %s", f)
	}

	return nil
}

// Update tldr pages from remote zip file
func (t *Tldr) Update() error {
	zipPath, err := download(t.pageSourceURL, t.path, filepath.Base(t.pageSourceURL))
	if err != nil {
		return fmt.Errorf("failed to download a tldr repository: %w", err)
	}
	defer os.Remove(zipPath)

	err = unzip(zipPath, t.path)
	if err != nil {
		return fmt.Errorf("failed to unzip a tldr repository: %w", err)
	}

	return nil
}

// FindPage find tldr page by `cmds`
func (t *Tldr) FindPage(cmds []string) (*Page, error) {
	page := strings.Join(cmds, "-") + ".md"
	for _, ptDir := range t.platforms {
		for _, lang := range t.languages {
			path := filepath.Join(t.path, getLangDir(lang), ptDir, page)
			if !pathExists(path) {
				// if cmd does not exist, try to find it in next platform/language
				continue
			}

			f, err := os.Open(path)
			if err != nil {
				return &Page{}, fmt.Errorf("failed to open the page (%s): %w", f.Name(), err)
			}
			defer f.Close()

			return parsePage(f)
		}
	}

	return &Page{}, fmt.Errorf("failed to find %s: %w", page, ErrNotFoundPage)
}

// Expired return true if tldr repository have passed `ttl`
func (t *Tldr) Expired(ttl time.Duration) bool {
	age, err := age(t.indexFilePath())
	if err != nil {
		return true
	}

	return age > ttl
}

func (t *Tldr) indexFilePath() string {
	return filepath.Join(t.path, "index.json")
}

// age return the time since the data exist at
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
