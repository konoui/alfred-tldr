package tldr

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParsePage(t *testing.T) {
	tests := []struct {
		description string
		url         string
		want        Page
		expectErr   bool
	}{
		{
			description: "lsof",
			url:         "https://raw.githubusercontent.com/tldr-pages/tldr/master/pages/common/lsof.md",
			expectErr:   false,
			want: Page{
				CmdName:        "lsof",
				CmdDescription: "Lists open files and the corresponding processes.\nNote: Root privileges (or sudo) is required to list files opened by others.\n",
				CmdExamples: []*CmdExample{
					&CmdExample{
						Description: "Find the processes that have a given file open:",
						Cmd:         "lsof {{path/to/file}}",
					},
					&CmdExample{
						Description: "Find the process that opened a local internet port:",
						Cmd:         "lsof -i :{{port}}",
					},
					&CmdExample{
						Description: "Only output the process ID (PID):",
						Cmd:         "lsof -t {{path/to/file}}",
					},
					&CmdExample{
						Description: "List files opened by the given user:",
						Cmd:         "lsof -u {{username}}",
					},
					&CmdExample{
						Description: "List files opened by the given command or process:",
						Cmd:         "lsof -c {{process_or_command_name}}",
					},
					&CmdExample{
						Description: "List files opened by a specific process, given its PID:",
						Cmd:         "lsof -p {{PID}}",
					},
					&CmdExample{
						Description: "List open files in a directory:",
						Cmd:         "lsof +D {{path/to/directory}}",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			path, err := download(tt.url, os.TempDir(), filepath.Base(tt.url))
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(path)

			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			page, err := parsePage(f)
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error want: %+v, got: %+v", tt.want, err.Error())
			}

			want := tt.want.CmdName
			got := page.CmdName
			if want != got {
				t.Errorf("unexpected response: want: %+v, got: %+v", want, got)
			}

			want = tt.want.CmdDescription
			got = page.CmdDescription
			if want != got {
				t.Errorf("unexpected response: want: %+v, got: %+v", want, got)
			}

			if !reflect.DeepEqual(tt.want.CmdExamples, tt.want.CmdExamples) {
				t.Errorf("unexpected response: want: %+v, got: %+v", tt.want.CmdExamples, tt.want.CmdExamples)
			}
		})
	}
}
