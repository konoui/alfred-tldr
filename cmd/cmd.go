package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/konoui/tldr/pkg/alfred"
	"github.com/konoui/tldr/pkg/tldr"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	outStream io.Writer = os.Stdout
	errStream io.Writer = os.Stderr
)

var (
	op tldr.Options
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

// NewRootCmd create a new cmd for root
func NewRootCmd() *cobra.Command {
	var isWorkflow bool
	rootCmd := &cobra.Command{
		Use:   "tldr <cmd>",
		Short: "show cmd examples",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if isWorkflow {
				return run(args[0], op, renderToWorkflow)
			}
			return run(args[0], op, renderToOut)
		},
		SilenceUsage: true,
	}
	rootCmd.PersistentFlags().StringVarP(&op.Platform, "platform", "p", op.Platform, "platform")
	//rootCmd.PersistentFlags().StringVarP(&op.Language, "language", "l", op.Language, "language")
	rootCmd.PersistentFlags().BoolVarP(&op.Update, "update", "u", op.Update, "update")
	rootCmd.PersistentFlags().BoolVarP(&isWorkflow, "workflow", "w", false, "rendering for alfred workflow")

	return rootCmd
}

// Execute Execute root cmd
func Execute(rootCmd *cobra.Command) {
	rootCmd.SetOutput(os.Stdout)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(errStream, err)
		os.Exit(1)
	}
}

type renderFunc func(*tldr.Page, bool, error)

func run(cmd string, op tldr.Options, f renderFunc) error {
	const tldrDir = ".tldr"
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, tldrDir)

	t := tldr.NewTldr(path, op)

	err = t.OnInitialize()
	isCacheExpired := tldr.IsCacheExpired(err)
	if !isCacheExpired && err != nil {
		return err
	}

	p, err := t.FindPage(cmd)
	f(p, isCacheExpired, err)

	return nil
}

const (
	bold  = "\x1b[1m"
	blue  = "\x1b[34m"
	green = "\x1b[32m"
	red   = "\x1b[31m"
	reset = "\x1b[33;0m"
)

func renderToOut(p *tldr.Page, isCacheExpired bool, pageErr error) {
	if isCacheExpired {
		fmt.Fprintf(errStream, "%s\n", tldr.CacheExpiredMsg)
	}

	if pageErr != nil {
		fmt.Fprintln(errStream, "This page doesn't exist yet!\nSubmit new pages here: https://github.com/tldr-pages/tldr")
		return
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
}

func renderToWorkflow(p *tldr.Page, isCacheExpired bool, pageErr error) {
	wf := alfred.New()
	for _, cmd := range p.CmdExamples {
		wf.Append(alfred.Item{
			Title:        cmd.Cmd,
			Subtitle:     cmd.Description,
			Autocomplete: cmd.Cmd,
			Arg:          cmd.Cmd,
		})
	}

	fmt.Fprintln(outStream, wf.Marshal())
}
