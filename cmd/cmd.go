package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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
	twoWeeks            = 2 * 7 * 24 * time.Hour
)

func initPlatform() string {
	platform := runtime.GOOS
	if platform == "darwin" {
		platform = "osx"
	}
	return platform
}

const (
	platformFlag = "platform"
	updateFlag   = "update"
	fuzzyFlag    = "fuzzy"
	versionFlag  = "version"
)

type config struct {
	platform   string
	language   string
	update     bool
	fuzzy      bool
	version    bool
	tldrClient *tldr.Tldr
}

// NewRootCmd create a new cmd for root
func NewRootCmd() *cobra.Command {
	cfg := new(config)
	rootCmd := &cobra.Command{
		Use:   "tldr <cmd>",
		Short: "show cmd examples",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cfg.initTldr(); err != nil {
				return err
			}
			if cfg.version {
				return cfg.printVersion(version, revision)
			}
			if cfg.update {
				return cfg.updateDB()
			}
			return cfg.printPage(args)
		},
		SilenceErrors:      true,
		SilenceUsage:       true,
		DisableSuggestions: true,
	}
	rootCmd.PersistentFlags().BoolVarP(&cfg.version, versionFlag, string(versionFlag[0]), false, "show the client version")
	rootCmd.PersistentFlags().BoolVarP(&cfg.update, updateFlag, string(updateFlag[0]), false, "update tldr database")
	rootCmd.PersistentFlags().BoolVarP(&cfg.fuzzy, fuzzyFlag, string(fuzzyFlag[0]), false, "use fuzzy search")
	rootCmd.PersistentFlags().StringVarP(&cfg.platform, platformFlag, string(platformFlag[0]), initPlatform(), "select from linux/osx/sunos/windows")

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
	showWorkflowUsage(cmd)
	return nil
}

func helpFunc(cmd *cobra.Command, args []string) {
	showWorkflowUsage(cmd)
}

func flagErrorFunc(cmd *cobra.Command, err error) error {
	showWorkflowUsage(cmd)
	return nil
}

func showWorkflowUsage(cmd *cobra.Command) {
	pflags := []*pflag.Flag{
		cmd.Flag(platformFlag),
		cmd.Flag(updateFlag),
		cmd.Flag(versionFlag),
	}

	for _, p := range pflags {
		awf.Append(
			makeUsageItem(p),
		)
	}
	awf.Output()
}

func makeUsageItem(p *pflag.Flag) *alfred.Item {
	title := fmt.Sprintf("-%s, --%s %s", p.Shorthand, p.Name, p.Usage)
	return alfred.NewItem().
		SetTitle(title).
		SetSubtitle(p.Usage)
}

func (cfg *config) initTldr() error {
	base, err := getDataDir()
	if err != nil {
		return err
	}
	path := filepath.Join(base, ".alfred-tldr")
	// Note turn off update option as we update database explicitly
	opt := &tldr.Options{
		Update:   false,
		Platform: cfg.platform,
		Language: cfg.language,
	}

	cfg.tldrClient = tldr.New(path, opt)
	return cfg.tldrClient.OnInitialize()
}

func (cfg *config) printPage(cmds []string) error {
	defer func() {
		_ = cfg.updateDBInBackground()
	}()
	if len(cmds) == 0 {
		awf.Append(
			alfred.NewItem().
				SetTitle("Please input a command").
				SetSubtitle("e.g.) tldr tar e.g.) tldr --help"),
		).Output()
		return nil
	}

	awf.EmptyWarning("No matching query", "Try a different query")
	p, err := cfg.tldrClient.FindPage(cmds)
	if err != nil {
		if errors.Is(err, tldr.ErrNoPage) && cfg.fuzzy {
			cfg.printFuzzyPages(cmds)
		} else {
			awf.Output()
		}
		return nil
	}

	for _, cmd := range p.CmdExamples {
		awf.Append(
			alfred.NewItem().
				SetTitle(cmd.Cmd).
				SetSubtitle(cmd.Description).
				SetArg(cmd.Cmd),
		)
	}

	awf.Output()
	return nil
}

func (cfg *config) printFuzzyPages(cmds []string) {
	index, err := cfg.tldrClient.LoadIndexFile()
	if err != nil {
		fatal(err)
	}

	suggestions := index.Commands.Search(cmds)
	for _, cmd := range suggestions {
		awf.Append(
			alfred.NewItem().
				SetTitle(cmd.Name).
				SetSubtitle(fmt.Sprintf("Platforms: %s", strings.Join(cmd.Platform, ","))).
				SetAutocomplete(cmd.Name).
				SetArg(fmt.Sprintf("%s --%s %s", cmd.Name, platformFlag, cmd.Platform[0])).
				SetVariable(nextActionKey, nextActionCmd).
				SetIcon(
					alfred.NewIcon().
						SetPath("candidate.png"),
				),
		)
	}

	awf.Output()
}

func (cfg *config) printVersion(v, r string) (_ error) {
	title := fmt.Sprintf("alfred-tldr %v(%s)", v, r)
	awf.Append(
		alfred.NewItem().SetTitle(title),
	).Output()
	return
}

func (cfg *config) updateDB() error {
	if !cfg.update {
		return errors.New("update is called even though update flag is not specified")
	}

	if shouldUpdateWithShell() {
		// update explicitly
		awf.Logf("updating tldr database...\n")
		return cfg.tldrClient.Update()
	}

	subtitle := ""
	if cfg.tldrClient.Expired(twoWeeks) {
		subtitle = "tldr database is older than 2 weeks"
	}
	awf.Append(
		alfred.NewItem().
			SetTitle("Please Enter if update tldr database").
			SetSubtitle(subtitle).
			SetVariable(nextActionKey, nextActionShell).
			SetArg("--"+updateFlag),
	// pass key/value as environment variable to next worfklow action
	).SetVariable(updateEnvKey, "true").Output()
	return nil
}

func (cfg *config) updateDBInBackground() error {
	if !isAutoUpdateEnabled() {
		awf.Logf("skip auto-update check as auto-update env is not enabled\n")
		return nil
	}

	if !cfg.tldrClient.Expired(twoWeeks) {
		return nil
	}

	jobName := "update"
	if awf.Job(jobName).IsRunning() {
		return nil
	}

	if err := os.Setenv(updateEnvKey, "true"); err != nil {
		return err
	}
	_, err := awf.Job(jobName).Start(os.Args[0], "--"+updateFlag)
	return err
}

// Execute Execute root cmd
func Execute(rootCmd *cobra.Command) {
	rootCmd.SetOut(outStream)
	rootCmd.SetErr(errStream)
	if err := rootCmd.Execute(); err != nil {
		fatal(err)
	}
}
