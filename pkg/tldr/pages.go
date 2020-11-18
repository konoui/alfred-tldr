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
	indexSourceURL = "https://tldr.sh/assets/index.json"
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
	path           string
	pageSourceURL  string
	indexSourceURL string
	indexFile      string
	zipFile        string
	platformDirs   []string
	languages      []string
	update         bool
}

var ErrNotFoundPage = errors.New("no page found")

// New create a instance of tldr repository
func New(tldrPath string, opt *Options) *Tldr {
	if opt == nil {
		opt = new(Options)
	}

	return &Tldr{
		path:           tldrPath,
		pageSourceURL:  pageSourceURL,
		indexSourceURL: indexSourceURL,
		indexFile:      filepath.Base(indexSourceURL),
		zipFile:        filepath.Base(pageSourceURL),
		platformDirs:   []string{opt.Platform, "common"},
		languages:      getLanguages(opt.Language),
		update:         opt.Update,
	}
}

func getLanguages(optionLang string) (priorities []string) {
	if optionLang != "" {
		return []string{optionLang}
	}

	// see https://github.com/tldr-pages/tldr/blob/master/CLIENT-SPECIFICATION.md#language
	langCode := getLanguageCode(os.Getenv("LANG"))
	if langCode == "" {
		return []string{languageCodeEN}
	}
	defer func() {
		if !contains(langCode, priorities) {
			priorities = append(priorities, langCode)
		}
		if !contains(languageCodeEN, priorities) {
			priorities = append(priorities, languageCodeEN)
		}
	}()

	languageEnv := os.Getenv("LANGUAGE")
	if languageEnv == "" {
		return
	}

	for _, language := range strings.Split(languageEnv, ":") {
		code := getLanguageCode(language)
		if contains(code, priorities) {
			continue
		}
		priorities = append(priorities, code)
	}
	return
}

func contains(target string, list []string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func getLanguageCode(language string) string {
	code := strings.SplitN(language, ".", 2)[0]
	if code == "C" || code == "POSIX" {
		return ""
	}

	for _, v := range []string{"pt_PT", "pt_BR", "zh_TW"} {
		if v == code {
			return code
		}
	}

	if code == "pt" {
		return "pt_PT"
	}

	return strings.SplitN(code, "_", 2)[0]
}

func getLangDir(lang string) string {
	if lang == "en" {
		return "pages"
	}

	return "pages" + "." + lang
}

// OnInitialize create and update tldr directory
func (t *Tldr) OnInitialize() error {
	initUpdate := false
	if !pathExists(t.path) {
		if err := os.MkdirAll(t.path, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create tldr dir %w", err)
		}
		// automatically updated if indexfile does not exist
		initUpdate = true
	}

	if t.update || initUpdate {
		if err := t.Update(); err != nil {
			return fmt.Errorf("failed to update tldr repository %w", err)
		}
	}

	return nil
}

// Update tldr pages from remote zip file
func (t *Tldr) Update() error {
	_, err := download(t.indexSourceURL, t.path, t.indexFile)
	if err != nil {
		return fmt.Errorf("failed to download a index file %w", err)
	}

	zipPath, err := download(t.pageSourceURL, t.path, t.zipFile)
	if err != nil {
		return fmt.Errorf("failed to download a tldr repository %w", err)
	}

	if err := unzip(zipPath, t.path); err != nil {
		return fmt.Errorf("failed to unzip a tldr repository %w", err)
	}

	return nil
}

// FindPage find tldr page by `cmds`
func (t *Tldr) FindPage(cmds []string) (*Page, error) {
	page := strings.Join(cmds, "-") + ".md"
	for _, pt := range t.platformDirs {
		for _, lang := range t.languages {
			path := filepath.Join(t.path, getLangDir(lang), pt, page)
			if !pathExists(path) {
				// if cmd does not exist, try to find it in next platform
				continue
			}

			f, err := os.Open(path)
			if err != nil {
				return &Page{}, fmt.Errorf("failed to open the page (%s) %w", f.Name(), err)
			}
			defer f.Close()

			return parsePage(f)
		}
	}

	return &Page{}, fmt.Errorf("failed to find %s %w", page, ErrNotFoundPage)
}

// Expired return true if tldr repository have passed `ttl`
func (t *Tldr) Expired(ttl time.Duration) bool {
	indexPath := filepath.Join(t.path, t.indexFile)
	age, err := age(indexPath)
	if err != nil {
		return true
	}

	return age > ttl
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
