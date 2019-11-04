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
			tldr := NewTldr(
				filepath.Join(os.TempDir(), ".tldr"),
				Options{Update: true},
			)
			index, err := tldr.LoadIndexFile()
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error want: got: %+v", err.Error())
			}

			if len(index.Commands) == 0 {
				t.Errorf("cannot load index file as commands length is 0")
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		description string
		query       string
		want        string
	}{
		{
			description: "similar ls cmd",
			query:       "lsa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tldr := NewTldr(
				filepath.Join(os.TempDir(), ".tldr"),
				Options{Update: true},
			)
			index, err := tldr.LoadIndexFile()
			if err != nil {
				t.Fatal(err)
			}

			got := index.Commands.Filter(tt.query)
			if len(got) == 0 {
				t.Errorf("filter result is 0")
			}

		})
	}
}
