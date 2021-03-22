package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	return "osx"
}

const (
	platformFlag = "platform"
	updateFlag   = "update"
	confirmFlag  = "confirm"
	fuzzyFlag    = "fuzzy"
	versionFlag  = "version"
	languageFlag = "language"
)

type config struct {
	platform   string
	language   string
	update     bool
	confirm    bool
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
	rootCmd.PersistentFlags().StringVarP(&cfg.platform, platformFlag, string(platformFlag[0]), initPlatform(), "select from linux/osx/sunos/windows")
	rootCmd.PersistentFlags().StringVarP(&cfg.language, languageFlag, "L", "", "select language e.g.) en")

	// internal flag
	rootCmd.PersistentFlags().BoolVarP(&cfg.confirm, confirmFlag, string(confirmFlag[0]), false, "confirmation for update")
	rootCmd.PersistentFlags().BoolVar(&cfg.fuzzy, fuzzyFlag, false, "use fuzzy search")

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
		cmd.Flag(languageFlag),
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
		Title(title).
		Subtitle(p.Usage).
		Valid(false)
}

func (cfg *config) initTldr() error {
	base, err := alfred.GetDataDir()
	if err != nil {
		return err
	}
	path := filepath.Join(base, "data")
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
	// insert update recommendation first
	if isUpdateRecommendEnabled() && cfg.tldrClient.Expired(twoWeeks) {
		awf.Append(
			alfred.NewItem().
				Title("Please Enter! Tldr database is older than 2 weeks").
				Arg(fmt.Sprintf("--%s --%s", updateFlag, confirmFlag)).
				Icon(alfred.IconAlertNote),
		).Variable(nextActionKey, nextActionShell)
	}

	// no input case
	if len(cmds) == 0 {
		awf.Append(
			alfred.NewItem().
				Title("Please input a command").
				Subtitle("e.g.) tldr tar e.g.) tldr --help").
				Valid(false),
		).Output()
		return nil
	}

	awf.SetEmptyWarning("No matching query", "Try a different query")
	p, err := cfg.tldrClient.FindPage(cmds)
	if err != nil {
		if errors.Is(err, tldr.ErrNotFoundPage) {
			if cfg.language != "" {
				awf.Clear().SetEmptyWarning(
					"Not found the command in selected language",
					"Try not to specify language option",
				).Output()
				return nil
			}
			if cfg.fuzzy {
				// list suggestions
				return cfg.printFuzzyPages(cmds)
			}
			awf.Output()
			return nil
		}
		return err
	}

	// TODO change icon for usage
	// descriptions has one line at least
	// see: https://github.com/tldr-pages/tldr/blob/master/contributing-guides/style-guide.md
	title := p.CmdDescriptions[0]
	subtitle := ""
	if len(p.CmdDescriptions) >= 2 {
		subtitle = p.CmdDescriptions[1]
	}
	awf.Append(
		alfred.NewItem().
			Title(title).
			Subtitle(subtitle).
			Valid(false).
			Icon(
				alfred.NewIcon().
					Path("description.png"),
			),
	)
	for _, cmd := range p.CmdExamples {
		awf.Append(
			alfred.NewItem().
				Title(cmd.Cmd).
				Subtitle(cmd.Description).
				Arg(cmd.Cmd),
		).Variable(nextActionKey, nextActionCopy)
	}

	awf.Output()
	return nil
}

func (cfg *config) printFuzzyPages(cmds []string) error {
	index, err := cfg.tldrClient.LoadIndexFile()
	if err != nil {
		return err
	}

	suggestions := index.Commands.Search(cmds)
	for _, cmd := range suggestions {
		complete := cmd.Name
		if cmd.Platform[0] != "common" {
			complete = fmt.Sprintf("-%s %s %s",
				string(platformFlag[0]),
				cmd.Platform[0],
				cmd.Name,
			)
		}
		awf.Append(
			alfred.NewItem().
				Title(cmd.Name).
				Subtitle(fmt.Sprintf("Platforms: %s", strings.Join(cmd.Platform, ","))).
				Valid(false).
				Autocomplete(complete).
				Icon(
					alfred.NewIcon().
						Path("candidate.png"),
				),
		)
	}

	awf.Output()
	return nil
}

func (cfg *config) printVersion(v, r string) (_ error) {
	title := fmt.Sprintf("alfred-tldr %v(%s)", v, r)
	awf.Append(
		alfred.NewItem().Title(title),
	).Output()
	return
}

func (cfg *config) updateDB() error {
	if !cfg.update {
		return errors.New("update is called even though update flag is not specified")
	}

	if cfg.confirm {
		// update explicitly
		awf.Logger().Infoln("updating tldr database...")
		return cfg.tldrClient.Update()
	}

	awf.Append(
		alfred.NewItem().
			Title("Please Enter if update tldr database").
			Arg(fmt.Sprintf("--%s --%s", updateFlag, confirmFlag)),
	).
		Variable(nextActionKey, nextActionShell).
		Output()

	return nil
}

// Execute Execute root cmd
func Execute(rootCmd *cobra.Command) {
	rootCmd.SetOut(outStream)
	rootCmd.SetErr(errStream)
	if err := rootCmd.Execute(); err != nil {
		fatal(err)
	}
}
