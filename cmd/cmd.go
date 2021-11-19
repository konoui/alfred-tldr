package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/konoui/alfred-tldr/pkg/tldr"
	"github.com/konoui/go-alfred"
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

var defaultPlatform = tldr.PlatformOSX

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

// NewRootCmd create a new cmd for root
func NewRootCmd(cfg *Config) *cobra.Command {

	var ptString string
	rootCmd := &cobra.Command{
		Use:   "tldr <cmd>",
		Short: "show cmd examples",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := awf.OnInitialize(); err != nil {
				return err
			}

			// normalize args since alfred workflow pass query as NFD
			for i, arg := range args {
				args[i] = alfred.Normalize(arg)
			}

			if err := cfg.setPlatform(ptString); err != nil {
				awf.SetEmptyWarning(strings.Title(err.Error()),
					"supported are linux/osx/sunos/windows").
					Output()
				return nil
			}

			if err := cfg.initTldr(); err != nil {
				return err
			}

			if cfg.version {
				return cfg.printVersion(version, revision)
			}
			if cfg.updateWorkflow {
				return cfg.updateTLDRWorkflow()
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
		cmd.Flag(longPlatformFlag),
		cmd.Flag(longUpdateFlag),
		cmd.Flag(longVersionFlag),
		cmd.Flag(longLanguageFlag),
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
	complete := fmt.Sprintf("--%s", p.Name)
	return alfred.NewItem().
		Title(title).
		Subtitle(p.Usage).
		Autocomplete(complete).
		Valid(false)
}

func (cfg *Config) setPlatform(ptString string) error {
	platforms := []tldr.Platform{
		tldr.PlatformCommon,
		tldr.PlatformLinux,
		tldr.PlatformOSX,
		tldr.PlatformWindows,
		tldr.PlatformSunos,
	}
	for _, pt := range platforms {
		if ptString == pt.String() {
			cfg.platform = pt
			return nil
		}
	}
	return fmt.Errorf("%s is unsupported platform", ptString)
}

func (cfg *Config) initTldr() error {
	path := filepath.Join(awf.GetDataDir(), "data")

	cfg.opts = append(cfg.opts,
		tldr.WithPlatform(cfg.platform),
		tldr.WithLanguage(cfg.language),
	)
	cfg.tldrClient = tldr.New(path,
		cfg.opts...,
	)

	ctx, cancel := context.WithTimeout(context.Background(), updateDBTimeout)
	defer cancel()
	return cfg.tldrClient.OnInitialize(ctx)
}

// Execute Execute root cmd
func Execute() {
	cfg := NewConfig()
	rootCmd := NewRootCmd(cfg)
	fatalIfErr(rootCmd.Execute())
}
