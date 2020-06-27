package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/konoui/go-alfred"
	"github.com/konoui/tldr/pkg/tldr"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
)

var (
	op         tldr.Options
	tldrMaxAge time.Duration = 24 * 7 * time.Hour
	updater                  = map[alfred.ModKey]alfred.Mod{
		alfred.ModCtrl: alfred.Mod{
			Subtitle: "update tldr repository",
			Arg:      "tldr --update",
			Variables: map[string]string{
				nextActionKey: nextActionShell,
			},
		},
	}
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

// NewRootCmd create a new cmd for root
func NewRootCmd() *cobra.Command {
	var isWorkflow bool
	var enableFuzzy bool
	rootCmd := &cobra.Command{
		Use:   "tldr <cmd>",
		Short: "show cmd examples",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args, op, isWorkflow, enableFuzzy)
		},
		SilenceErrors: true,
		SilenceUsage:  false,
	}
	rootCmd.PersistentFlags().StringVarP(&op.Platform, "platform", "p", op.Platform, "platform")
	//rootCmd.PersistentFlags().StringVarP(&op.Language, "language", "l", op.Language, "language")
	rootCmd.PersistentFlags().BoolVarP(&op.Update, "update", "u", op.Update, "update")
	rootCmd.PersistentFlags().BoolVarP(&isWorkflow, "workflow", "w", false, "rendering for alfred workflow")
	rootCmd.PersistentFlags().BoolVarP(&enableFuzzy, "fuzzy", "f", false, "enable fuzzy search for cmds")

	return rootCmd
}

// Execute Execute root cmd
func Execute(rootCmd *cobra.Command) {
	rootCmd.SetOut(outStream)
	rootCmd.SetErr(errStream)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmds []string, op tldr.Options, isWorkflow, enableFuzzy bool) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, tldrDir)

	t := tldr.New(path, op)

	err = t.OnInitialize()
	if err != nil {
		return err
	}

	if isWorkflow {
		renderToWorkflow(t, cmds, enableFuzzy)
		return nil
	}

	return renderToOut(t, cmds)
}

const (
	bold  = "\x1b[1m"
	blue  = "\x1b[34m"
	green = "\x1b[32m"
	red   = "\x1b[31m"
	reset = "\x1b[33;0m"
)

func renderToOut(t *tldr.Tldr, cmds []string) error {
	if len(cmds) == 0 {
		return fmt.Errorf("no argument")
	}
	if t.Expired(tldrMaxAge) {
		fmt.Fprintln(errStream, "more than a week passed, should update tldr using --update")
	}

	p, err := t.FindPage(cmds)
	if err != nil {
		fmt.Fprintln(errStream, "This page doesn't exist yet!\nSubmit new pages here: https://github.com/tldr-pages/tldr")
		return nil
	}

	coloredCmdName := bold + p.CmdName + reset
	fmt.Fprintln(outStream, coloredCmdName)
	fmt.Fprintln(outStream)
	fmt.Fprintln(outStream, p.CmdDescription)
	for _, cmd := range p.CmdExamples {
		coloredDescription := "- " + green + cmd.Description + reset
		fmt.Fprintln(outStream, coloredDescription)
		line := strings.Replace(cmd.Cmd, "{{", blue, -1)
		line = strings.Replace(line, "}}", red, -1)
		coloredCmd := red + line + reset
		fmt.Fprintln(outStream, coloredCmd)
		fmt.Fprintln(outStream)
	}
	return nil
}

// decide next action for workflow filter
const (
	nextActionKey   = "nextAction"
	nextActionCmd   = "cmd"
	nextActionShell = "shell"
)

func renderToWorkflow(t *tldr.Tldr, cmds []string, enableFuzzy bool) {
	awf := alfred.NewWorkflow()
	awf.SetOut(outStream)
	awf.SetErr(errStream)
	awf.EmptyWarning("No matching query", "Try a different query")

	p, _ := t.FindPage(cmds)
	for _, cmd := range p.CmdExamples {
		awf.Append(&alfred.Item{
			Title:    cmd.Cmd,
			Subtitle: cmd.Description,
			Arg:      cmd.Cmd,
			Mods:     updater,
		})
	}

	if enableFuzzy && len(p.CmdExamples) == 0 {
		index, err := t.LoadIndexFile()
		if err != nil {
			awf.Fatal("fatal error occurs", err.Error())
			return
		}

		suggestions := index.Commands.Search(cmds)
		for _, cmd := range suggestions {
			awf.Append(&alfred.Item{
				Title:        cmd.Name,
				Subtitle:     fmt.Sprintf("Platforms: %s", strings.Join(cmd.Platform, ",")),
				Autocomplete: cmd.Name,
				Variables: map[string]string{
					nextActionKey: nextActionCmd,
				},
				Arg: fmt.Sprintf("%s --platform %s", cmd.Name, cmd.Platform[0]),
				Icon: &alfred.Icon{
					Path: "candidate.png",
				},
				Mods: updater,
			})
		}
	}

	awf.Output()
}
