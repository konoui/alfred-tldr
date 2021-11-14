package tldr

import (
	"context"
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
				WithForceUpdate(),
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
	tests := []struct {
		description string
		tldr        *Tldr
		expectErr   bool
	}{
		{
			description: "success test for expected",
			expectErr:   false,
			tldr: New(
				filepath.Join(os.TempDir(), ".tldr"),
				WithForceUpdate(),
			),
		},
		{
			description: "failed test due to permission deny",
			expectErr:   true,
			tldr: New(
				"/.tldr",
			),
		},
		{
			description: "failed test due to invalid url",
			expectErr:   true,
			tldr: New(
				filepath.Join(os.TempDir(), ".tldr"),
				WithRepositoryURL("https://google.com/index.html"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := tt.tldr.Update(context.TODO())
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
	tests := []struct {
		description string
		tldr        *Tldr
		expectErr   bool
	}{
		{
			description: "success test for expected",
			expectErr:   false,
			tldr: New(
				filepath.Join(os.TempDir(), ".tldr"),
				WithForceUpdate(),
			),
		},
		{
			description: "failed test due to permission deny",
			expectErr:   true,
			tldr: New(
				"/.tldr",
				WithPlatform(PlatformLinux),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := tt.tldr.OnInitialize(context.TODO())
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
	tests := []struct {
		description string
		tldr        *Tldr
		want        bool
		tldrTTL     time.Duration
	}{
		{
			description: "failed test due to expired cache",
			tldr: New(
				filepath.Join(os.TempDir(), ".tldr"),
			),
			tldrTTL: 0 * time.Hour,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := tt.tldr.OnInitialize(context.TODO())
			if err != nil {
				t.Fatal(err)
			}

			if got := tt.tldr.Expired(tt.tldrTTL); got != tt.want {
				t.Errorf("want: %+v, got: %+v", tt.want, got)
			}
		})
	}
}
