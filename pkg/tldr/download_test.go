package tldr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		description string
		url         string
		want        string
		expectErr   bool
	}{
		{
			description: "download example1",
			url:         "https://raw.githubusercontent.com/tldr-pages/tldr/master/pages/common/lsof.md",
			want:        "lsof.md",
			expectErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// download on current dir
			got, err := download(tt.url, "", tt.want)
			defer os.RemoveAll(got)
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error want: %+v, got: %+v", tt.want, err.Error())
			}

			if tt.want != got {
				t.Errorf("unexpected response: want: %+v, got: %+v", tt.want, got)
			}
		})
	}
}

func TestUnzip(t *testing.T) {
	tests := []struct {
		description string
		url         string
		expectErr   bool
	}{
		{
			description: "download and unzip 1",
			url:         "https://tldr.sh/assets/tldr.zip",
			expectErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tmpDir := os.TempDir()
			path, err := download(tt.url, tmpDir, filepath.Base(tt.url))
			if err != nil {
				t.Fatalf("faltal error: %+v", err)
			}
			defer os.RemoveAll(path)

			if err := unzip(path, tmpDir); !tt.expectErr && err != nil {
				t.Errorf("unexpected error want: %+v", err.Error())
			}
		})
	}
}
