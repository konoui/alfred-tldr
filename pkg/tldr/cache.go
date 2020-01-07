package tldr

import (
	"os"
	"time"
)

type cacheExpiredError struct {
	msg string
}

// CacheMaxAge tldr page cache age. default is a week.
// you should not override.
var CacheMaxAge time.Duration = 24 * 7 * time.Hour

func (e *cacheExpiredError) Error() string {
	return e.msg
}

// IsCacheExpired return true if `err` means expired cache
func IsCacheExpired(err error) bool {
	if _, ok := err.(*cacheExpiredError); ok {
		return true
	}

	return false
}

// expired return true if cache is expired
func expired(path string, maxAge time.Duration) bool {
	age, err := age(path)
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
