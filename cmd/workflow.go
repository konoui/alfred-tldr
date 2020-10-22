package cmd

import (
	"os"

	"github.com/konoui/go-alfred"
)

// decide next action for workflow filter
const (
	nextActionKey   = "nextAction"
	nextActionCmd   = "cmd"
	nextActionShell = "shell"
)

var awf = alfred.NewWorkflow()

func init() {
	awf.SetOut(outStream)
	awf.SetErr(errStream)
}

func getDataDir() (string, error) {
	base, err := alfred.GetDataDir()
	if err != nil {
		// Note fallback to home directory
		base, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}
	return base, nil
}

func shouldUpdateInShell(updateFlag string) bool {
	// nextCmd is defined alfred workflow
	v := os.Getenv("nextCmd")
	return v == updateFlag
}

func fatal(err error) {
	awf.Fatal("Fatal errors occur", err.Error())
}
