package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/konoui/alfred-tldr/pkg/tldr"
	tldrtest "github.com/konoui/alfred-tldr/pkg/tldr/test"
	"github.com/konoui/go-alfred"
	"github.com/konoui/go-alfred/env"
	"github.com/konoui/go-alfred/update"
	mock "github.com/konoui/go-alfred/update/mock_update"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

// global test server instance
var testServer = tldrtest.NewServer().Start()

func testdataPath(file string) string {
	return filepath.Join("testdata", file)
}

func setup(t *testing.T, command string,
	alfredOpts ...alfred.Option) (awf *alfred.Workflow, cmd *cobra.Command, outBuf, errBuf *bytes.Buffer) {
	t.Helper()

	tmpDir := t.TempDir()
	for k, v := range map[string]string{
		env.KeyWorkflowData:     "/tmp",
		env.KeyWorkflowCache:    tmpDir,
		env.KeyWorkflowBundleID: "test-bundle-id",
	} {
		t.Setenv(k, v)
	}

	cmdArgs, err := shellwords.Parse(command)
	if err != nil {
		t.Fatalf("args parse error: %+v", err)
	}

	outBuf, errBuf = new(bytes.Buffer), new(bytes.Buffer)
	awf = alfred.NewWorkflow(append(alfredOpts,
		alfred.WithOutWriter(outBuf),
		alfred.WithLogWriter(errBuf),
		alfred.WithArguments(cmdArgs...),
		alfred.WithInitializers(nil),
	)...)
	cfg := NewConfig()
	// set dummy url for local test
	cfg.tldrOpts = append(cfg.tldrOpts, tldr.WithRepositoryURL(testServer.TldrZipURL()))
	cmd = NewRootCmd(cfg, awf)

	cmd.SetOut(outBuf)
	cmd.SetErr(errBuf)
	cmd.SetArgs(cmdArgs)

	return
}

func execute(t *testing.T, awf *alfred.Workflow, rootCmd *cobra.Command, wantExitCode int) {
	exitCode := awf.RunSimple(rootCmd.Execute)
	if exitCode != wantExitCode {
		t.Errorf("unexpected exit code want %d, got %d", wantExitCode, exitCode)
	}
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
		up     func()
		down   func()
	}{
		{
			name: "lsof",
			args: args{
				command:  "lsof",
				filepath: "output-lsof.json",
			},
		},
		{
			name: "sub command git checkout",
			args: args{
				command:  "git checkout",
				filepath: "output-git-checkout.json",
			},
		},
		{
			name: "tar with uppercase output format",
			args: args{
				command:  "tar",
				filepath: "output-tar-with-uppercase-format.json",
			},
			up:   func() { os.Setenv(envKeyCommandFormat, "uppercase") },
			down: func() { os.Unsetenv(envKeyCommandFormat) },
		},
		{
			name: "tar with single output format",
			args: args{
				command:  "tar",
				filepath: "output-tar-with-single-format.json",
			},
			up:   func() { os.Setenv(envKeyCommandFormat, "single") },
			down: func() { os.Unsetenv(envKeyCommandFormat) },
		},
		{
			name: "tar with original output format",
			args: args{
				command:  "tar",
				filepath: "output-tar-with-original-format.json",
			},
			up:   func() { os.Setenv(envKeyCommandFormat, "original") },
			down: func() { os.Unsetenv(envKeyCommandFormat) },
		},
		{
			name: "tar with remove output format",
			args: args{
				command:  "tar",
				filepath: "output-tar-with-remove-format.json",
			},
			up:   func() { os.Setenv(envKeyCommandFormat, "remove") },
			down: func() { os.Unsetenv(envKeyCommandFormat) },
		},
		{
			name: "tar and open url with ctrl mod key",
			args: args{
				command:  "tar",
				filepath: "output-tar-with-ctrl-mod-key.json",
			},
			up:   func() { os.Setenv(envKeyOpenURLMod, "ctrl") },
			down: func() { os.Unsetenv(envKeyOpenURLMod) },
		},
		{
			name: "fuzzy search returns git checkout",
			args: args{
				command:  "gitchec --fuzzy",
				filepath: "output-git-checkout-with-fuzzy.json",
			},
		},
		{
			name: "fuzzy search returns non-common platform",
			args: args{
				command:  "pstree --fuzzy",
				filepath: "output-pstree-with-fuzzy.json",
			},
		},
		{
			name: "multiple platform command archey without -p flag returns computer platform OSX for autocomplete",
			args: args{
				command:  "arche --fuzzy",
				filepath: "output-archey-without-pflag-fuzzy.json",
			},
		},
		{
			name: "multiple platform command archey with -p flag returns specified platform for autocomplete",
			args: args{
				command:  "arche -p linux --fuzzy",
				filepath: "output-archey-with-pflag-with-fuzzy.json",
			},
		},
		{
			name: "version flag is the highest priority",
			args: args{
				command:  "-v -u tldr -L -a",
				filepath: "output-version.json",
			},
		},
		{
			name: "input something but commands not found",
			args: args{
				command:  "dummy-empty-result",
				filepath: "output-empty-result.json",
			},
		},
		{
			name: "fuzzy search but commands not found",
			args: args{
				command:  "--fuzzy dummy-empty-result",
				filepath: "output-empty-result.json",
			},
		},
		{
			name: "specify language flag but commands not found",
			args: args{
				command:  "-L invalid-lang tar",
				filepath: "output-language-empty-result.json",
			},
		},
		{
			name: "no input returns no-input message",
			args: args{
				command:  "",
				filepath: "output-no-input.json",
			},
		},
		{
			name: "specify help flag returns usages",
			args: args{
				command:  "--help",
				filepath: "output-usage.json",
			},
		},
		{
			name: "string flag -L without any values and invalid flag return no-input message",
			args: args{
				command:  "-L -a",
				filepath: "output-no-input.json",
			},
		},
		{
			name: "string flag -L without any value returns usage",
			args: args{
				command:  "--fuzzy -L",
				filepath: "output-usage.json",
			},
		},
		{
			name: "invalid bool flag returns usage",
			args: args{
				command:  "lsof -a",
				filepath: "output-usage.json",
			},
		},
		{
			name: "invalid bool flag and valid flag return usage",
			args: args{
				command:  "-a -u",
				filepath: "output-usage.json",
			},
		},
		{
			name: "invalid platform value returns platform error message",
			args: args{
				command:  "-p a",
				filepath: "output-platform-error.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.up != nil || tt.down != nil {
				t.Cleanup(tt.down)
				tt.up()
			}

			testpath := testdataPath(tt.args.filepath)
			wantData, err := os.ReadFile(testpath)
			if err != nil {
				t.Fatal(err)
			}

			awf, cmd, outBuf, _ := setup(t, tt.args.command)
			execute(t, awf, cmd, 0)
			outGotData := outBuf.Bytes()

			// automatically update test data
			if tt.update {
				if err := writeFile(testpath, outGotData); err != nil {
					t.Fatal(err)
				}
			}

			if diff := alfred.DiffOutput(wantData, outGotData); diff != "" {
				t.Errorf("-want +got\n%+v", diff)
			}
		})
	}
}

func TestUpdateConfirmation(t *testing.T) {
	type args struct {
		filepath              string
		command               string
		dbTTL                 time.Duration
		newerVersionAvailable bool
	}
	tests := []struct {
		name         string
		args         args
		update       bool
		wantExitCode int
	}{
		{
			name: "no input and update recommendations",
			args: args{
				dbTTL:                 0,
				newerVersionAvailable: true,
				command:               "",
				filepath:              "output-update-recommendations.json",
			},
		},
		{
			name: "lsof with workflow update recommendation",
			args: args{
				dbTTL:                 1000 * time.Hour,
				newerVersionAvailable: true,
				command:               "lsof",
				filepath:              "output-lsof-with-update-workflow-recommendation.json",
			},
		},
		{
			name: "lsof with db recommendation",
			args: args{
				dbTTL:                 0,
				newerVersionAvailable: false,
				command:               "lsof",
				filepath:              "output-lsof-with-update-db-recommendation.json",
			},
		},
		{
			name: "update db confirmation",
			args: args{
				dbTTL:                 0,
				newerVersionAvailable: false,
				command:               "--update",
				filepath:              "output-update-db-confirmation.json",
			},
		},
		{
			name: "prints update confirmation when specified --update flag and ignore argument",
			args: args{
				dbTTL:                 0,
				newerVersionAvailable: false,
				command:               "--update tldr",
				filepath:              "output-update-confirmation.json",
			},
		},
		{
			name: "update-workflow confirmation does not support",
			args: args{
				command:  "--update-workflow",
				filepath: "output-update-workflow-flag-error.json",
			},
			wantExitCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Setenv(envKeyUpdateDBRecommendation, "true"); err != nil {
				t.Fatal(err)
			}
			if err := os.Setenv(envKeyUpdateWorkflowRecommendation, "true"); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := os.Unsetenv(envKeyUpdateDBRecommendation); err != nil {
					t.Fatal(err)
				}
				if err := os.Unsetenv(envKeyUpdateWorkflowRecommendation); err != nil {
					t.Fatal(err)
				}
			})

			testpath := testdataPath(tt.args.filepath)
			wantData, err := os.ReadFile(testpath)
			if err != nil {
				t.Fatal(err)
			}

			mockSource, teardown := newMockUpdaterSource(t, tt.args.newerVersionAvailable)
			defer teardown()

			// disable ttl
			twoWeeks = tt.args.dbTTL
			awf, cmd, outBuf, errBuf := setup(t, tt.args.command, alfred.WithUpdater(mockSource))
			execute(t, awf, cmd, tt.wantExitCode)
			outGotData := outBuf.Bytes()

			// automatically update test data
			if tt.update {
				if err := writeFile(testpath, outGotData); err != nil {
					t.Fatal(err)
				}
			}

			if diff := alfred.DiffOutput(wantData, outGotData); diff != "" {
				t.Errorf("-want +got\n%+v\n", diff)
				t.Errorf("error output: %s", errBuf.String())
			}
		})
	}
}

func newMockUpdaterSource(t *testing.T, newerVersionAvailable bool) (_ update.UpdaterSource, _ func()) {
	ctrl := gomock.NewController(t)
	mockSource := mock.NewMockUpdaterSource(ctrl)
	mockUpdater := mock.NewMockUpdater(ctrl)
	mockUpdater.EXPECT().Update(gomock.Any()).Return(nil).AnyTimes()
	mockSource.EXPECT().IsNewVersionAvailable(gomock.Any()).Return(newerVersionAvailable, nil).AnyTimes()
	mockSource.EXPECT().IfNewVersionAvailable().Return(mockUpdater).AnyTimes()
	return mockSource, ctrl.Finish
}

func TestUpdateExecution(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name        string
		args        args
		expectedErr bool
		errMsg      string
		wantMsg     string
	}{
		{
			name: "update db returns succeeded message",
			args: args{
				command: "--update --confirm",
			},
			wantMsg:     "update succeeded",
			expectedErr: false,
		},
		{
			name: "when update-workflow without updater, nil updater returns error. update execution outputs message to stdout",
			args: args{
				command: "--update-workflow --confirm",
			},
			expectedErr: false,
			wantMsg:     "update failed due to no implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			awf, cmd, outBuf, errBuf := setup(t, tt.args.command)
			exitCode := awf.RunSimple(cmd.Execute)
			if tt.expectedErr && exitCode == 0 {
				t.Errorf("unexpected results")
			}
			if tt.expectedErr && exitCode != 0 {
				if !strings.Contains(errBuf.String(), tt.errMsg) {
					t.Errorf("want: %v\n got: %v", tt.errMsg, errBuf)
				}
			}

			got := outBuf.String()
			if !strings.Contains(got, tt.wantMsg) {
				t.Errorf("want: %v\n got: %v", tt.wantMsg, got)
			}
		})
	}
}

func Test_choicePlatform(t *testing.T) {
	type args struct {
		pts      []tldr.Platform
		selected tldr.Platform
	}
	tests := []struct {
		name string
		args args
		want tldr.Platform
	}{
		{
			name: "return common when platforms do not contain selected",
			args: args{
				pts: []tldr.Platform{
					tldr.PlatformCommon,
				},
				selected: tldr.PlatformLinux,
			},
			want: tldr.PlatformCommon,
		},
		{
			name: "return selected insted of common when platforms contain selected",
			args: args{
				pts: []tldr.Platform{
					tldr.PlatformCommon,
					tldr.PlatformLinux,
				},
				selected: tldr.PlatformLinux,
			},
			want: tldr.PlatformLinux,
		},
		{
			name: "returns common when platforms do not contain selected",
			args: args{
				pts: []tldr.Platform{
					tldr.PlatformCommon,
					tldr.PlatformLinux,
				},
				selected: tldr.PlatformOSX,
			},
			want: tldr.PlatformCommon,
		},
		{
			name: "returns commond when empty pts",
			args: args{
				pts:      []tldr.Platform{},
				selected: tldr.PlatformCommon,
			},
			want: tldr.PlatformCommon,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := choicePlatform(tt.args.pts, tt.args.selected); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("choicePlatform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func writeFile(filename string, data []byte) error {
	pretty := new(bytes.Buffer)
	if err := json.Indent(pretty, data, "", "  "); err != nil {
		return err
	}

	return os.WriteFile(filename, pretty.Bytes(), 0o600)
}
