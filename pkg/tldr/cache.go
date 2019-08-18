package tldr

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Cache a simple API
type Cache struct {
	dir    string
	file   string
	maxAge time.Duration
}

// NewCache creates a new cache Instance
func NewCache(dir, file string, maxAge time.Duration) *Cache {
	if !pathExists(dir) {
		panic(fmt.Errorf("%s directory does not exist", dir))
	}
	return &Cache{dir: dir, file: file, maxAge: maxAge}
}

// Clear remove cache file if exists
func (c Cache) Clear() error {
	p := c.path()
	if pathExists(p) {
		return os.RemoveAll(p)
	}
	return nil
}

// NotExpired return true if cache is not expired
func (c Cache) NotExpired() bool {
	return !c.Expired()
}

// Expired return true if cache is expired
func (c Cache) Expired() bool {
	age, err := c.Age()
	if err != nil {
		return true
	}
	return age > c.maxAge
}

// Age return the time since the data is cached at
func (c Cache) Age() (time.Duration, error) {
	p := c.path()
	fi, err := os.Stat(p)
	if err != nil {
		return time.Duration(0), err
	}
	return time.Since(fi.ModTime()), nil
}

// Exists return true if the cache file exists
func (c Cache) Exists() bool {
	return pathExists(c.path())
}

// path return the path of cache file
func (c Cache) path() string {
	return filepath.Join(c.dir, c.file)
}

// pathExists return true if path exists
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
