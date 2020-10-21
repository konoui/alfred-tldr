package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/konoui/alfred-tldr/pkg/tldr"
	"github.com/konoui/go-alfred"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
	version             = "*"
	revision            = "*"
	op        *tldr.Options
)

func init() {
	platform := runtime.GOOS
	if platform == "darwin" {
		platform = "osx"
	}
	op = &tldr.Options{
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
	rootCmd.PersistentFlags().BoolVarP(&v, versionFlag, string(versionFlag[0]), false, "show the client version")
	rootCmd.PersistentFlags().BoolVarP(&op.Update, updateFlag, string(updateFlag[0]), false, "update tldr database")
	rootCmd.PersistentFlags().BoolVarP(&enableFuzzy, fuzzyFlag, string(fuzzyFlag[0]), false, "use fuzzy search")
	rootCmd.PersistentFlags().StringVarP(&op.Platform, platformFlag, string(platformFlag[0]), op.Platform, "select from linux/osx/sunos/windows")

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

func makeUsageMap(cmd *cobra.Command) map[string]*alfred.Item {
	flags := []*pflag.Flag{
		cmd.Flag(platformFlag),
		cmd.Flag(updateFlag),
		cmd.Flag(versionFlag),
	}

	m := make(map[string]*alfred.Item, len(flags))
	for _, f := range flags {
		m[f.Name] = makeItem(f)
	}
	return m
}

func makeItem(p *pflag.Flag) *alfred.Item {
	title := fmt.Sprintf("-%s, --%s %s", p.Shorthand, p.Name, p.Usage)
	return alfred.NewItem().
		SetTitle(title).
		SetSubtitle(p.Usage)
}

func run(cmds []string, op *tldr.Options, enableFuzzy bool) (err error) {
	var base string
	base, err = getDataDir()
	if err != nil {
		// Note fallback to home directory
		base, err = os.UserHomeDir()
		if err != nil {
			return
		}
	}

	path := filepath.Join(base, ".alfred-tldr")
	t := tldr.New(path, op)

	err = t.OnInitialize()
	if err != nil {
		return
	}

	workflowOutput(t, cmds, enableFuzzy)
	return
}

// Execute Execute root cmd
func Execute(rootCmd *cobra.Command) {
	rootCmd.SetOut(outStream)
	rootCmd.SetErr(errStream)
	if err := rootCmd.Execute(); err != nil {
		fatal(err)
	}
}
