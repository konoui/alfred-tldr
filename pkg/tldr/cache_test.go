package tldr

import (
	"os"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	tests := []struct {
		description string
		dir         string
		want        string
		expectErr   bool
	}{
		{description: "exists cache dir", dir: os.TempDir(), want: "", expectErr: false},
		{description: "no exists cache dir", dir: "/unk", want: "/unk directory does not exist", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			defer func() {
				err := recover()
				if tt.expectErr && err == nil {
					t.Errorf("expect error happens, but got response")
				}
				if !tt.expectErr && err != nil {
					t.Errorf("unexpected error want: %+v, got: %+v", tt.want, err)
				}
				if err != nil && err.(error).Error() != tt.want {
					t.Errorf("unexpected response: want: %+v, got: %+v", tt.want, err)
				}
			}()
			NewCache(tt.dir, "", 1*time.Minute)
		})
	}
}

func TestCacheExpired(t *testing.T) {
	tests := []struct {
		description string
		dir         string
		file        string
		expiredTime time.Duration
		expired     bool
	}{
		{description: "test for not expired cache", dir: "./", file: "test1", expiredTime: 3 * time.Minute, expired: false},
		{description: "test for expired cache", dir: "./", file: "test1", expiredTime: 0 * time.Minute, expired: true},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if _, err := os.Create(tt.file); err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tt.file)

			cache := NewCache(tt.dir, tt.file, tt.expiredTime)
			if !tt.expired && cache.Expired() {
				t.Errorf("cache should not be expired but cache expired.")
			}
			if tt.expired && cache.NotExpired() {
				t.Errorf("cache should be expired but not expired.")
			}
		})
	}
}
