package tldr

import (
	"fmt"
	"os"
	"time"
)

type errorCacheExpired struct {
	error
}

// CacheMaxAge tldr page cache age. default is a week.
// you should not override.
var CacheMaxAge time.Duration = 24 * 7 * time.Hour

func (e *errorCacheExpired) Error() string {
	return e.error.Error()
}

// IsCacheExpired return true if `err` means expired cache
func IsCacheExpired(err error) bool {
	if _, ok := err.(*errorCacheExpired); !ok {
		return false
	}

	return true
}

// pathExists return true if path exists
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// age return the time since the data is cached at
func age(path string) (time.Duration, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return time.Duration(0), err
	}

	return time.Since(fi.ModTime()), nil
}

func expiredMsg(t time.Duration) string {
	const msg = "should update tldr using --update"
	week := t / (24 * time.Hour * 7)
	if week > 0 {
		return fmt.Sprintf("more than a week passed, %s", msg)
	}

	return msg
}
