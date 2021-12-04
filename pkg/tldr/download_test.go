package tldr

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		want      string
		expectErr bool
	}{
		{
			name:      "download example1",
			url:       testServer.TldrZipURL(),
			want:      filepath.Base(testServer.TldrZipURL()),
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// since download on current directory, got is only filename
			got, err := download(context.TODO(), tt.url, "", tt.want)
			t.Cleanup(func() { os.RemoveAll(got) })
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err)
			}

			if tt.want != got {
				t.Errorf("want: %+v, got: %+v", tt.want, got)
			}
		})
	}
}

func TestUnzip(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		expectErr bool
	}{
		{
			name:      "download and unzip 1",
			url:       testServer.TldrZipURL(),
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			path, err := download(context.TODO(), tt.url, tmpDir, filepath.Base(tt.url))
			if err != nil {
				t.Fatalf("faltal error: %+v", err)
			}
			defer os.RemoveAll(path)

			if err := unzip(context.TODO(), path, tmpDir); !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err)
			}
		})
	}
}
