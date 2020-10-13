package cmd

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/konoui/tldr/pkg/tldr"
	"github.com/spf13/cobra"
)

var (
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
)

var (
	op         tldr.Options
	tldrMaxAge time.Duration = 24 * 7 * time.Hour
)

const tldrDir = ".tldr"

func init() {
	platform := runtime.GOOS
	if platform == "darwin" {
		platform = "osx"
	}
	op = tldr.Options{
		Platform: platform,
		Language: "",
		Update:   false,
	}
}

const (
	platformFlag = "platform"
	updateFlag   = "update"
	fuzzyFlag    = "fuzzy"
)

// NewRootCmd create a new cmd for root
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tldr <cmd>",
		Short: "show cmd examples",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			enableFuzzy := getBoolFlag(cmd, fuzzyFlag)
			platform := getStringFlag(cmd, platformFlag)
			if platform != "" {
				op.Platform = platform
			}
			return run(args, op, enableFuzzy)
		},
		SilenceErrors:      true,
		SilenceUsage:       true,
		DisableSuggestions: true,
	}
	rootCmd.PersistentFlags().StringP(platformFlag, string(platformFlag[0]), "", "select platform")
	rootCmd.PersistentFlags().BoolP(updateFlag, string(updateFlag[0]), false, "update tldr repository")
	rootCmd.PersistentFlags().BoolP(fuzzyFlag, string(fuzzyFlag[0]), false, "use fuzzy search")
	rootCmd.SetUsageFunc(usageFunc)
	rootCmd.SetHelpFunc(helpFunc)
	rootCmd.SetFlagErrorFunc(flagErrorFunc)

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	return rootCmd
}

func usageFunc(cmd *cobra.Command) error {
	showWorkflowUsage(makeUsageMap(cmd))
	return nil
}

func helpFunc(cmd *cobra.Command, args []string) {
	showWorkflowUsage(makeUsageMap(cmd))
}

func flagErrorFunc(cmd *cobra.Command, err error) error {
	showWorkflowUsage(makeUsageMap(cmd))
	return nil
}

// Execute Execute root cmd
func Execute(rootCmd *cobra.Command) {
	rootCmd.SetOut(outStream)
	rootCmd.SetErr(errStream)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmds []string, op tldr.Options, enableFuzzy bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, tldrDir)

	t := tldr.New(path, op)

	err = t.OnInitialize()
	if err != nil {
		return err
	}

	renderToWorkflow(t, cmds, enableFuzzy)
	return nil
}
