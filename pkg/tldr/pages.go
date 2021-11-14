package tldr

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Platform string

const (
	PlatformCommon  Platform = "common"
	PlatformWindows Platform = "windows"
	PlatformLinux   Platform = "linux"
	PlatformSunos   Platform = "sunos"
	PlatformOSX     Platform = "osx"
)

const (
	PageSourceURL  = "https://tldr.sh/assets/tldr.zip"
	languageCodeEN = "en"
)

var (
	ErrNotFoundPage = errors.New("no page found")
)

func (pt Platform) String() string {
	return string(pt)
}

type Option func(*Tldr)

func WithPlatform(pt Platform) Option {
	return func(t *Tldr) {
		// Note specified platform is highest priority
		t.platforms = append([]Platform{pt}, t.platforms...)
	}
}

func WithLanguage(lang string) Option {
	return func(t *Tldr) {
		t.languages = getLanguages(lang)
	}
}

func WithForceUpdate() Option {
	return func(t *Tldr) {
		t.update = true
	}
}

// WithRepositoryURL replaces default tldr remote url
// This is useful for local test
func WithRepositoryURL(u string) Option {
	return func(t *Tldr) {
		t.pageSourceURL = u
	}
}

// Tldr Repository of tldir pages
type Tldr struct {
	path          string
	pageSourceURL string
	platforms     []Platform
	languages     []string
	update        bool
}

// New create a instance of tldr repository
func New(tldrPath string, opts ...Option) *Tldr {

	t := &Tldr{
		path:          tldrPath,
		pageSourceURL: PageSourceURL,
		platforms:     []Platform{PlatformCommon},
		languages:     getLanguages(""),
		update:        false,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// OnInitialize create and update tldr directory
func (t *Tldr) OnInitialize(ctx context.Context) error {
	initUpdate := false
	if !pathExists(t.path) {
		if err := os.MkdirAll(t.path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create tldr dir: %w", err)
		}
		// automatically updated if indexfile does not exist
		initUpdate = true
	}

	if t.update || initUpdate {
		if err := t.Update(ctx); err != nil {
			return fmt.Errorf("failed to update tldr repository: %w", err)
		}
	}

	if f := t.indexFilePath(); !pathExists(f) {
		return fmt.Errorf("tldr database is broken %s", f)
	}

	return nil
}

// Update tldr pages from remote zip file
func (t *Tldr) Update(ctx context.Context) error {
	zipPath, err := download(ctx, t.pageSourceURL, t.path, filepath.Base(t.pageSourceURL))
	if err != nil {
		return fmt.Errorf("failed to download a tldr repository: %w", err)
	}
	defer os.Remove(zipPath)

	err = unzip(ctx, zipPath, t.path)
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
			path := filepath.Join(t.path, getLangDir(lang), ptDir.String(), page)
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
