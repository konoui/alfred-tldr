package tldr

import (
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
				Options{Update: true},
			)
			if err := tldr.OnInitialize(); err != nil {
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
			description: "similar ls cmd",
			query:       []string{"lsa"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tldr := New(
				filepath.Join(os.TempDir(), ".tldr"),
				Options{},
			)
			if err := tldr.OnInitialize(); err != nil {
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

		})
	}
}
