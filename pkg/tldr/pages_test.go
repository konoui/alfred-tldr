package tldr

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const tldrZipFilename = "tldr.zip"

var testServer *httptest.Server

func tmpDir() string {
	return "/tmp"
}

func serverURL() string {
	return testServer.URL
}

func tldrZipURL() string {
	return serverURL() + "/" + tldrZipFilename
}

func init() {
	setupTldrRepositoryServer()
}

// global test server
func setupTldrRepositoryServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, tldrZipFilename) {
			fmt.Fprintf(w, "hello")
			return
		}

		zipPath := filepath.Join(tmpDir(), tldrZipFilename)
		if _, err := os.Stat(zipPath); err != nil {
			panic(err)
		}

		f, err := os.Open(zipPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := io.Copy(w, f); err != nil {
			panic(err)
		}
	})

	// set to global val
	testServer = httptest.NewUnstartedServer(mux)
	testServer.Start()
}

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
			description: "valid sub cmd",
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
				WithRepositoryURL(tldrZipURL()),
				WithLanguage("en"),
			)
			err := tldr.OnInitialize(context.TODO())
			if err != nil {
				t.Fatal(err)
			}

			page, err := tldr.FindPage(tt.cmds)
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err)
			}
			if got := page.CmdName; got != tt.want {
				t.Errorf("want: %+v, got: %+v", tt.want, got)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		opts     []Option
		tldrPath string
	}
	tests := []struct {
		name      string
		args      args
		expectErr bool
	}{
		{
			name:      "success test for expected",
			expectErr: false,
			args: args{
				tldrPath: filepath.Join(os.TempDir(), ".tldr"),
				opts: []Option{
					WithRepositoryURL(tldrZipURL()),
				},
			},
		},
		{
			name:      "failed test due to permission deny",
			expectErr: true,
			args: args{
				tldrPath: "/.tldr",
				opts: []Option{
					WithRepositoryURL(tldrZipURL()),
				},
			},
		},
		{
			name:      "failed test due to invalid url",
			expectErr: true,
			args: args{
				tldrPath: filepath.Join(os.TempDir(), ".tldr"),
				opts: []Option{
					WithRepositoryURL(serverURL() + "/" + "invalid-file"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := New(tt.args.tldrPath, tt.args.opts...)
			err := i.Update(context.TODO())
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err)
			}
		})
	}
}

func TestOnInitialize(t *testing.T) {
	type args struct {
		opts     []Option
		tldrPath string
	}
	tests := []struct {
		name      string
		args      args
		expectErr bool
	}{
		{
			name:      "success test for expected",
			expectErr: false,
			args: args{
				tldrPath: filepath.Join(os.TempDir(), ".tldr"),
				opts: []Option{
					WithRepositoryURL(tldrZipURL()),
					WithForceUpdate(),
				},
			},
		},
		{
			name:      "failed test due to permission deny",
			expectErr: true,
			args: args{
				tldrPath: "/.tldr",
				opts: []Option{
					WithPlatform(PlatformLinux),
					WithRepositoryURL(tldrZipURL()),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := New(tt.args.tldrPath, tt.args.opts...)
			err := i.OnInitialize(context.TODO())
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err)
			}
		})
	}
}

func TestExpired(t *testing.T) {
	type args struct {
		opts     []Option
		tldrPath string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		tldrTTL time.Duration
	}{
		{
			name: "failed test due to expired cache",
			args: args{
				tldrPath: filepath.Join(os.TempDir(), ".tldr"),
				opts: []Option{
					WithRepositoryURL(tldrZipURL()),
				},
			},
			tldrTTL: 0 * time.Hour,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := New(tt.args.tldrPath, tt.args.opts...)
			err := i.OnInitialize(context.TODO())
			if err != nil {
				t.Fatal(err)
			}

			if got := i.Expired(tt.tldrTTL); got != tt.want {
				t.Errorf("want: %+v, got: %+v", tt.want, got)
			}
		})
	}
}
