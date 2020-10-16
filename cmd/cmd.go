package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/konoui/alfred-tldr/pkg/tldr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
	version             = "*"
	revision            = "*"
	op        tldr.Options
)

const (
	tldrDir = ".tldr"
)

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
	versionFlag  = "version"
)

// NewRootCmd create a new cmd for root
func NewRootCmd() *cobra.Command {
	var (
		enableFuzzy bool
		v           bool
	)
	rootCmd := &cobra.Command{
		Use:   "tldr <cmd>",
		Short: "show cmd examples",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if v {
				return printVersion(version, revision)
			}
			return run(args, op, enableFuzzy)
		},
		SilenceErrors:      true,
		SilenceUsage:       true,
		DisableSuggestions: true,
	}
	rootCmd.PersistentFlags().BoolVarP(&v, versionFlag, string(versionFlag[0]), false, "select platform")
	rootCmd.PersistentFlags().BoolVarP(&op.Update, updateFlag, string(updateFlag[0]), false, "update tldr repository")
	rootCmd.PersistentFlags().BoolVarP(&enableFuzzy, fuzzyFlag, string(fuzzyFlag[0]), false, "use fuzzy search")
	rootCmd.PersistentFlags().StringVarP(&op.Platform, platformFlag, string(platformFlag[0]), op.Platform, "select platform, supported are linux/osx/sunos/windows")

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

func makeUsageMap(cmd *cobra.Command) (m map[string]string) {
	pf := cmd.Flag(platformFlag)
	uf := cmd.Flag(updateFlag)
	m = make(map[string]string, 2)
	m[pf.Name] = makeDescription(pf)
	m[uf.Name] = makeDescription(uf)
	return
}

func makeDescription(p *pflag.Flag) string {
	return fmt.Sprintf("Usage: -%s, --%s %s", p.Shorthand, p.Name, p.Usage)
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

	workflowOutput(t, cmds, enableFuzzy)
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
