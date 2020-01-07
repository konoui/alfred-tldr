package tldr

import (
	"os"
	"testing"
	"time"
)

func TestCacheExpired(t *testing.T) {
	tests := []struct {
		description string
		filepath    string
		expiredTime time.Duration
		expired     bool
	}{
		{
			description: "test for not expired cache",
			filepath:    "./test1",
			expiredTime: 3 * time.Minute,
			expired:     false,
		},
		{
			description: "test for expired cache",
			filepath:    "./test1",
			expiredTime: 0 * time.Minute,
			expired:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if _, err := os.Create(tt.filepath); err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tt.filepath)

			age, err := age(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}

			if !tt.expired && age > tt.expiredTime {
				t.Errorf("cache should not be expired but cache expired.")
			}
		})
	}
}
