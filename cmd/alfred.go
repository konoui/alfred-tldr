package cmd

import (
	"fmt"
	"strings"

	"github.com/konoui/go-alfred"
	"github.com/konoui/tldr/pkg/tldr"
)

// decide next action for workflow filter
const (
	nextActionKey = "nextAction"
	nextActionCmd = "cmd"
)

func showWorkflowUsage(usageMap map[string]string) {
	awf := alfred.NewWorkflow()
	awf.SetOut(outStream)
	awf.SetErr(errStream)
	for _, u := range usageMap {
		awf.Append(&alfred.Item{
			Title: u,
		})
	}
	awf.Output()
}

func renderToWorkflow(t *tldr.Tldr, cmds []string, enableFuzzy bool) {
	awf := alfred.NewWorkflow()
	awf.SetOut(outStream)
	awf.SetErr(errStream)
	awf.EmptyWarning("No matching query", "Try a different query")

	p, _ := t.FindPage(cmds)
	for _, cmd := range p.CmdExamples {
		awf.Append(
			alfred.NewItem().
				SetTitle(cmd.Cmd).
				SetSubtitle(cmd.Description).
				SetArg(cmd.Cmd),
		)
	}

	if enableFuzzy && len(p.CmdExamples) == 0 {
		index, err := t.LoadIndexFile()
		if err != nil {
			awf.Fatal("Fatal errors occur", err.Error())
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
	}

	awf.Output()
}
