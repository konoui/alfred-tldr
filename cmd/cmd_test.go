package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/konoui/tldr/pkg/tldr"
	"github.com/mattn/go-shellwords"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		description string
		expectErr   bool
		filepath    string
		command     string
		cacheMaxAge time.Duration
		errMsg      string
	}{
		{
			description: "text output tests with lsof",
			expectErr:   false,
			command:     "lsof --update",
			filepath:    "./test_output_lsof.txt",
			cacheMaxAge: tldr.CacheMaxAge,
			errMsg:      "",
		},
		{
			description: "alfred workflow tests with lsof",
			expectErr:   false,
			command:     "lsof --update --workflow",
			filepath:    "./test_output_lsof.json",
			cacheMaxAge: tldr.CacheMaxAge,
			errMsg:      "",
		},
		{
			description: "text output tests sub cmd tests with git checkout",
			expectErr:   false,
			command:     "git checkout --update",
			filepath:    "./test_output_git-checkout.txt",
			cacheMaxAge: tldr.CacheMaxAge,
			errMsg:      "",
		},
		{
			description: "alfred workflow sub cmd tests with git checkout",
			expectErr:   false,
			command:     "git checkout --update --workflow",
			filepath:    "./test_output_git-checkout.json",
			cacheMaxAge: tldr.CacheMaxAge,
			errMsg:      "",
		},
		{
			description: "text output tests with expired cache with lsof.",
			expectErr:   false,
			command:     "lsof",
			filepath:    "./test_output_lsof.txt",
			cacheMaxAge: 0 * time.Hour,
			errMsg:      fmt.Sprintln("should update tldr using --update"),
		},
		{
			description: "alfred workflow tests with expired cache with lsof. alfred workflow show no error",
			expectErr:   false,
			command:     "lsof --workflow",
			filepath:    "./test_output_lsof.json",
			cacheMaxAge: 0 * time.Hour,
			errMsg:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// set cache max age
			tldr.CacheMaxAge = tt.cacheMaxAge

			f, err := os.Open(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			v, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
			outStream, errStream = outBuf, errBuf
			cmdArgs, err := shellwords.Parse(tt.command)
			if err != nil {
				t.Fatalf("args parse error: %+v", err)
			}
			rootCmd := NewRootCmd()
			rootCmd.SetOutput(outStream)
			rootCmd.SetArgs(cmdArgs)

			err = rootCmd.Execute()
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error want: got: %+v", err.Error())
			}

			got := outBuf.String()
			if want := string(v); want != got {
				t.Errorf("unexpected response: want: \n%+v, got: \n%+v", want, got)
			}

			got = errBuf.String()
			if want := tt.errMsg; want != got {
				t.Errorf("unexpected response: want: \n%+v, got: \n%+v", want, got)
			}

		})
	}
}
