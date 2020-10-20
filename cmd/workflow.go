package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/konoui/alfred-tldr/pkg/tldr"
	"github.com/konoui/go-alfred"
)

// decide next action for workflow filter
const (
	nextActionKey = "nextAction"
	nextActionCmd = "cmd"
)

var awf = alfred.NewWorkflow()

func init() {
	awf.SetOut(outStream)
	awf.SetErr(errStream)
}

func getDataDir() (string, error) {
	return alfred.GetDataDir()
}

func showWorkflowUsage(usageMap map[string]*alfred.Item) {
	for _, i := range usageMap {
		awf.Append(
			i,
		)
	}
	awf.Output()
}

func workflowOutput(t *tldr.Tldr, cmds []string, enableFuzzy bool) {
	if len(cmds) == 0 {
		awf.Append(
			alfred.NewItem().
				SetTitle("Please input a command").
				SetSubtitle("e.g.) tldr tar e.g.) tldr --help"),
		).Output()
		return
	}

	awf.EmptyWarning("No matching query", "Try a different query")
	p, err := t.FindPage(cmds)
	if err != nil {
		if errors.Is(err, tldr.ErrNoPage) && enableFuzzy {
			fuzzyOutput(t, cmds)
		} else {
			awf.Output()
		}
		return
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
}

func fuzzyOutput(t *tldr.Tldr, cmds []string) {
	index, err := t.LoadIndexFile()
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

func printVersion(v, r string) (_ error) {
	title := fmt.Sprintf("alfred-tldr %v(%s)", v, r)
	awf.Append(
		alfred.NewItem().SetTitle(title),
	).Output()
	return
}

func fatal(err error) {
	awf.Fatal("Fatal errors occur", err.Error())
}
