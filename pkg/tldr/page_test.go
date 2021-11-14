package tldr

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParsePage(t *testing.T) {
	tests := []struct {
		description string
		filepath    string
		want        *Page
		expectErr   bool
	}{
		{
			description: "lsof",
			filepath:    "testdata/lsof.md",
			expectErr:   false,
			want: &Page{
				CmdName: "lsof",
				CmdDescriptions: []string{
					"Lists open files and the corresponding processes.",
					"Note: Root privileges (or sudo) is required to list files opened by others.",
					"More information: <https://manned.org/lsof>.",
				},
				CmdExamples: []*CmdExample{
					{
						Description: "Find the processes that have a given file open:",
						Cmd:         "lsof {{path/to/file}}",
					},
					{
						Description: "Find the process that opened a local internet port:",
						Cmd:         "lsof -i :{{port}}",
					},
					{
						Description: "Only output the process ID (PID):",
						Cmd:         "lsof -t {{path/to/file}}",
					},
					{
						Description: "List files opened by the given user:",
						Cmd:         "lsof -u {{username}}",
					},
					{
						Description: "List files opened by the given command or process:",
						Cmd:         "lsof -c {{process_or_command_name}}",
					},
					{
						Description: "List files opened by a specific process, given its PID:",
						Cmd:         "lsof -p {{PID}}",
					},
					{
						Description: "List open files in a directory:",
						Cmd:         "lsof +D {{path/to/directory}}",
					},
					{
						Description: "Find the process that is listening on a local IPv6 TCP port and don't convert network or port numbers:",
						Cmd:         "lsof -i6TCP:{{port}} -sTCP:LISTEN -n -P",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			f, err := os.Open(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			got, err := parsePage(f)
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error got: %+v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
