package tldr

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadIndexFile(t *testing.T) {
	tests := []struct {
		description string
		expectErr   bool
	}{
		{
			description: "load index file correctly",
			expectErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tldr := New(
				filepath.Join(os.TempDir(), ".tldr"),
				WithTestZipURL(),
			)
			if err := tldr.OnInitialize(context.TODO()); err != nil {
				t.Fatal(err)
			}
			index, err := tldr.LoadIndexFile()
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err)
			}

			if index != nil && len(index.Commands) == 0 {
				t.Errorf("cannot load index file as commands length is 0")
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		description string
		query       []string
		want        string
	}{
		{
			description: "similar html5validator cmd",
			query:       []string{"html5"},
			want:        "html5validator",
		},
		{
			description: "no hyphen git check",
			query:       []string{"git", "check"},
			want:        "git checkout",
		},
		{
			description: "hyphen apt- key",
			query:       []string{"apt-", "get"},
			want:        "apt-get",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tldr := New(
				filepath.Join(os.TempDir(), ".tldr"),
				WithTestZipURL(),
			)
			if err := tldr.OnInitialize(context.TODO()); err != nil {
				t.Fatal(err)
			}
			index, err := tldr.LoadIndexFile()
			if err != nil {
				t.Fatal(err)
			}

			got := index.Commands.Search(tt.query)
			if len(got) == 0 {
				t.Errorf("filter result is 0")
			}
			for _, cmd := range got {
				if cmd.Name == tt.want {
					return
				}
			}
			t.Errorf("not found %s in cmds", tt.want)
		})
	}
}
