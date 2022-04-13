package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/konoui/alfred-tldr/pkg/tldr"
	"github.com/konoui/go-alfred"
	"github.com/konoui/go-alfred/initialize"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	version                    = "*"
	revision                   = "*"
	twoWeeks                   = 2 * 7 * 24 * time.Hour
	updateDBTimeout            = 30 * time.Second
	updateWorkflowTimeout      = 1 * time.Minute
	updateWorkflowCheckTimeout = 5 * time.Second
)

var (
	defaultPlatform = tldr.PlatformOSX
	defaultOpts     = []alfred.Option{
		alfred.WithMaxResults(30),
		alfred.WithGitHubUpdater(
			"konoui",
			"alfred-tldr",
			version,
			getUpdateWorkflowInterval(twoWeeks),
		),
		alfred.WithOutWriter(os.Stdout),
		alfred.WithLogWriter(os.Stderr),
		alfred.WithInitializers(
			initialize.NewEmbedAssets(),
		),
	}
)

const (
	longPlatformFlag   = "platform"
	longUpdateFlag     = "update"
	longVersionFlag    = "version"
	longLanguageFlag   = "language"
	confirmFlag        = "confirm"
	fuzzyFlag          = "fuzzy"
	updateWorkflowFlag = "update-workflow"
)

var (
	platformFlag = string(longPlatformFlag[0])
	updateFlag   = string(longUpdateFlag[0])
	versionFlag  = string(longVersionFlag[0])
	languageFlag = strings.ToUpper(string(longLanguageFlag[0]))
)

type client struct {
	cfg *Config
	*alfred.Workflow
	tldrClient *tldr.Tldr
}

// NewRootCmd create a new cmd for root
func NewRootCmd(cfg *Config, awf *alfred.Workflow) *cobra.Command {

	var ptString string
	c := &client{
		cfg:      cfg,
		Workflow: awf,
	}
	rootCmd := &cobra.Command{
		Use:   "tldr <cmd>",
		Short: "show cmd examples",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := cfg.setPlatform(ptString); err != nil {
				awf.SetEmptyWarning(err.Error(),
					"supported are linux/osx/sunos/windows").
					Output()
				return nil
			}

			tc, err := newTldrClient(cfg, awf)
			if err != nil {
				return err
			}

			c.tldrClient = tc
			// update and normalize
			args = c.UpdateOpts(alfred.WithArguments(args...)).Args()

			switch {
			case cfg.version:
				return printVersion(c, version, revision)
			case cfg.updateWorkflow:
				return updateTLDRWorkflow(c)
			case cfg.update:
				return updateDB(c)
			default:
				return printPage(c, args)
			}
		},
		SilenceErrors:      true,
		SilenceUsage:       true,
		DisableSuggestions: true,
	}
	rootCmd.PersistentFlags().BoolVarP(&cfg.version, longVersionFlag, versionFlag,
		false, "show the client version")
	rootCmd.PersistentFlags().BoolVarP(&cfg.update, longUpdateFlag, updateFlag,
		false, "update tldr database")
	rootCmd.PersistentFlags().StringVarP(&ptString, longPlatformFlag, platformFlag,
		defaultPlatform.String(), "select from linux/osx/sunos/windows")
	rootCmd.PersistentFlags().StringVarP(&cfg.language, longLanguageFlag, languageFlag, "", "select language e.g.) en")

	// internal flag
	rootCmd.PersistentFlags().BoolVar(&cfg.confirm, confirmFlag, false, "confirmation for update")
	rootCmd.PersistentFlags().BoolVar(&cfg.fuzzy, fuzzyFlag, false, "use fuzzy search")
	rootCmd.PersistentFlags().BoolVar(&cfg.updateWorkflow, updateWorkflowFlag, false, "update tldr workflow if possible")

	rootCmd.SetUsageFunc(getUsageFunc(c))
	rootCmd.SetHelpFunc(getHelpFunc(c))
	rootCmd.SetFlagErrorFunc(getFlagErrorFunc(c))
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	return rootCmd
}

func getHelpFunc(c *client) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		usageFunc := getUsageFunc(c)
		_ = usageFunc(cmd)
	}
}

func getFlagErrorFunc(c *client) func(*cobra.Command, error) error {
	return func(cmd *cobra.Command, err error) error {
		usageFunc := getUsageFunc(c)
		return usageFunc(cmd)
	}
}

func getUsageFunc(c *client) func(*cobra.Command) error {
	return func(cmd *cobra.Command) error {
		pflags := []*pflag.Flag{
			cmd.Flag(longPlatformFlag),
			cmd.Flag(longUpdateFlag),
			cmd.Flag(longVersionFlag),
			cmd.Flag(longLanguageFlag),
		}

		for _, p := range pflags {
			c.Append(
				makeUsageItem(p),
			)
		}
		c.Output()
		return nil
	}
}

func makeUsageItem(p *pflag.Flag) *alfred.Item {
	title := fmt.Sprintf("-%s, --%s %s", p.Shorthand, p.Name, p.Usage)
	complete := fmt.Sprintf("--%s", p.Name)
	return alfred.NewItem().
		Title(title).
		Subtitle(p.Usage).
		Autocomplete(complete).
		Valid(false)
}

// Execute executes root cmd
func Execute() {
	cfg := NewConfig()
	awf := alfred.NewWorkflow(defaultOpts...)
	rootCmd := NewRootCmd(cfg, awf)
	os.Exit(awf.RunSimple(rootCmd.Execute))
}
