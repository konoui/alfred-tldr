package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
			name: "lsof",
			args: args{
				command:  "lsof",
				filepath: testdataPath("test_output_lsof.json"),
			},
			update: true,
		},
		{
			name: "sub command git checkout",
			args: args{
				command:  "git checkout",
				filepath: testdataPath("test_output_git-checkout.json"),
			},
		},
		{
			name: "fuzzy search",
			args: args{
				command:  "gitchec --fuzzy",
				filepath: testdataPath("test_output_git-checkout_with_fuzzy.json"),
			},
		},
		{
			name: "show no error when cache expired",
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
			name: "specify language flag but not found",
			args: args{
				command:  "-L ja tar",
				filepath: testdataPath("test_output_language-empty-result.json"),
			},
		},
		{
			name: "empty result",
			args: args{
				command:  "ta",
				filepath: testdataPath("test_output_empty-result.json"),
			},
		},
		{
			name: "fuzzy but empty result",
			args: args{
				command:  "--fuzzy aaaaa",
				filepath: testdataPath("test_output_empty-result.json"),
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
				if err := writeFile(tt.args.filepath, outGotData); err != nil {
					t.Fatal(err)
				}
			}

			if diff := alfred.DiffScriptFilter(wantData, outGotData); diff != "" {
				t.Errorf("-want +got\n%+v", diff)
			}
		})
	}
}

func writeFile(filename string, data []byte) error {
	pretty := new(bytes.Buffer)
	if err := json.Indent(pretty, data, "", "  "); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, pretty.Bytes(), 0600); err != nil {
		return err
	}
	return nil
}
