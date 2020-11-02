package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/konoui/go-alfred"
	"github.com/mattn/go-shellwords"
)

func testdataPath(file string) string {
	return filepath.Join("./testdata", file)
}

func TestExecute(t *testing.T) {
	type args struct {
		filepath string
		command  string
	}
	tests := []struct {
		name   string
		args   args
		update bool
	}{
		{
			name: "alfred workflow. lsof",
			args: args{
				command:  "lsof",
				filepath: testdataPath("test_output_lsof.json"),
			},
		},
		{
			name: "alfred workflow. sub command git checkout",
			args: args{
				command:  "git checkout",
				filepath: testdataPath("test_output_git-checkout.json"),
			},
		},
		{
			name: "alfred workflow. fuzzy search",
			args: args{
				command:  "gitchec --fuzzy",
				filepath: testdataPath("test_output_git-checkout_with_fuzzy.json"),
			},
		},
		{
			name: "alfred workflow. show no error when cache expired",
			args: args{
				command:  "lsof",
				filepath: testdataPath("test_output_lsof.json"),
			},
		},
		{
			name: "version flag is the highest priority",
			args: args{
				command:  "-v tldr -p -a",
				filepath: testdataPath("test_output_version.json"),
			},
		},
		{
			name: "print update confirmation when specified --update flag and ignore argument",
			args: args{
				command:  "--update tldr",
				filepath: testdataPath("test_output_update-confirmation.json"),
			},
		},
		{
			name: "no input",
			args: args{
				command:  "",
				filepath: testdataPath("test_output_no-input.json"),
			},
		},
		{
			name: "string flag but no value and invalid flag",
			args: args{
				command:  "-p -a",
				filepath: testdataPath("test_output_no-input.json"),
			},
		},
		{
			name: "string flag but no value",
			args: args{
				command:  "--fuzzy -p",
				filepath: testdataPath("test_output_usage.json"),
			},
		},
		{
			name: "bool invalid flag",
			args: args{
				command:  "lsof -a",
				filepath: testdataPath("test_output_usage.json"),
			},
		},
		{
			name: " bool invalid and valid flags",
			args: args{
				command:  "-a -u",
				filepath: testdataPath("test_output_usage.json"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Unsetenv(updateEnvKey); err != nil {
				t.Fatal(err)
			}
			wantData, err := ioutil.ReadFile(tt.args.filepath)
			if err != nil {
				t.Fatal(err)
			}

			// overwrite global awf
			awf = alfred.NewWorkflow()
			rootCmd := NewRootCmd()
			outBuf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
			outStream, errStream = outBuf, errBuf
			cmdArgs, err := shellwords.Parse(tt.args.command)
			if err != nil {
				t.Fatalf("args parse error: %+v", err)
			}
			awf.SetOut(outBuf)
			awf.SetErr(errBuf)
			rootCmd.SetOutput(outStream)
			rootCmd.SetArgs(cmdArgs)

			err = rootCmd.Execute()
			if err != nil {
				t.Errorf("unexpected error got: %+v", err)
			}

			outGotData := outBuf.Bytes()

			// automatically update test data
			if tt.update {
				if err := ioutil.WriteFile(tt.args.filepath, outGotData, 0644); err != nil {
					t.Fatal(err)
				}
			}

			if diff := alfred.DiffScriptFilter(wantData, outGotData); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
