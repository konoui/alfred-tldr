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
	updateEnvKey    = "ALFRED_TLDR_UPDATE"
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

func shouldUpdateWithShell() bool {
	// will use auto update for self execution
	v := os.Getenv(updateEnvKey)
	return v != ""
}

func fatal(err error) {
	awf.Fatal("Fatal errors occur", err.Error())
}
