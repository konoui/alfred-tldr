package tldr

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFindPage(t *testing.T) {
	tests := []struct {
		description string
		want        string
		expectErr   bool
		cmds        []string
	}{
		{
			description: "valid cmd",
			expectErr:   false,
			want:        "lsof",
			cmds:        []string{"lsof"},
		},
		{
			description: "valid cmd and sub cmd",
			expectErr:   false,
			want:        "git checkout",
			cmds:        []string{"git", "checkout"},
		},
		{
			description: "invalid cmd, response will be empty Page struct",
			expectErr:   true,
			want:        "",
			cmds:        []string{"lsofaaaaaaaaaaaaaaa"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tldr := New(
				filepath.Join(os.TempDir(), ".tldr"),
				Options{Update: true},
			)
			err := tldr.OnInitialize()
			if err != nil {
				t.Fatal(err)
			}

			page, err := tldr.FindPage(tt.cmds)
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err.Error())
			}
			if got := page.CmdName; got != tt.want {
				t.Errorf("want: %+v, got: %+v", tt.want, got)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		description string
		tldr        Tldr
		expectErr   bool
	}{
		{
			description: "success test for expected",
			expectErr:   false,
			tldr: Tldr{
				path:           filepath.Join(os.TempDir(), ".tldr"),
				pageSourceURL:  pageSourceURL,
				indexSourceURL: indexSourceURL,
				indexFile:      filepath.Base(indexSourceURL),
				zipFile:        filepath.Base(pageSourceURL),
				update:         true,
			},
		},
		{
			description: "failed test due to permission deny",
			expectErr:   true,
			tldr: Tldr{
				path:           "/.tldr",
				pageSourceURL:  pageSourceURL,
				indexSourceURL: indexSourceURL,
				indexFile:      filepath.Base(indexSourceURL),
				zipFile:        filepath.Base(pageSourceURL),
			},
		},
		{
			description: "failed test due to invalid url",
			expectErr:   true,
			tldr: Tldr{
				path:           filepath.Join(os.TempDir(), ".tldr"),
				pageSourceURL:  "https://google.com/index.html",
				indexSourceURL: indexSourceURL,
				indexFile:      filepath.Base(indexSourceURL),
				zipFile:        filepath.Base(pageSourceURL),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := tt.tldr.Update()
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err.Error())
			}
		})
	}
}

func TestOnInitialize(t *testing.T) {
	tests := []struct {
		description string
		tldr        Tldr
		expectErr   bool
	}{
		{
			description: "success test for expected",
			expectErr:   false,
			tldr: Tldr{
				path:           filepath.Join(os.TempDir(), ".tldr"),
				pageSourceURL:  pageSourceURL,
				indexSourceURL: indexSourceURL,
				indexFile:      filepath.Base(indexSourceURL),
				zipFile:        filepath.Base(pageSourceURL),
				update:         true,
				cacheMaxAge:    24 * 7 * time.Hour,
			},
		},
		{
			description: "failed test due to permission deny",
			expectErr:   true,
			tldr: Tldr{
				path:           "/.tldr",
				pageSourceURL:  pageSourceURL,
				indexSourceURL: indexSourceURL,
				indexFile:      filepath.Base(indexSourceURL),
				zipFile:        filepath.Base(pageSourceURL),
				platformDirs:   []string{"linux", "common"},
				langDir:        "pages",
				update:         false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := tt.tldr.OnInitialize()
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err.Error())
			}
		})
	}
}

func TestExpired(t *testing.T) {
	tests := []struct {
		description string
		tldr        Tldr
		want        bool
		tldrMaxAge  time.Duration
	}{
		{
			description: "failed test due to expired cache",
			tldr: Tldr{
				path:           filepath.Join(os.TempDir(), ".tldr"),
				pageSourceURL:  pageSourceURL,
				indexSourceURL: indexSourceURL,
				indexFile:      filepath.Base(indexSourceURL),
				zipFile:        filepath.Base(pageSourceURL),
				update:         false,
				cacheMaxAge:    0 * time.Hour,
			},
			tldrMaxAge: 0 * time.Hour,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := tt.tldr.OnInitialize()
			if err != nil {
				t.Fatal(err)
			}

			if got := tt.tldr.Expired(tt.tldrMaxAge); got != tt.want {
				t.Errorf("want: %+v, got: %+v", tt.want, got)
			}
		})
	}
}
