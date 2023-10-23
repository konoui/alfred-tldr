package tldr

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	tldrtest "github.com/konoui/alfred-tldr/pkg/tldr/test"
)

// global test server instance
var testServer = tldrtest.NewServer().Start()

// heler function for test
func WithTestZipURL() Option {
	return WithRepositoryURL(testServer.TldrZipURL())
}

func WithTestInvalidURL() Option {
	return WithRepositoryURL(testServer.ServerURL() + "/" + "invalid-file")
}

func TestFindPage(t *testing.T) {
	tests := []struct {
		description string
		want        string
		cmds        []string
		expectErr   bool
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
				WithRepositoryURL(testServer.TldrZipURL()),
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
		tldrPath string
		opts     []Option
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
					WithTestZipURL(),
				},
			},
		},
		{
			name:      "failed test due to permission deny",
			expectErr: true,
			args: args{
				tldrPath: "/.tldr",
				opts: []Option{
					WithTestZipURL(),
				},
			},
		},
		{
			name:      "failed test due to invalid url",
			expectErr: true,
			args: args{
				tldrPath: filepath.Join(os.TempDir(), ".tldr"),
				opts: []Option{
					WithTestInvalidURL(),
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
		tldrPath string
		opts     []Option
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
					WithForceUpdate(),
					WithTestZipURL(),
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
					WithTestZipURL(),
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
		tldrPath string
		opts     []Option
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
					WithTestZipURL(),
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
