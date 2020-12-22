package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/konoui/go-alfred"
)

// decide next action for workflow filter
const (
	nextActionKey   = "nextAction"
	nextActionCopy  = "copy"
	nextActionCmd   = "cmd"
	nextActionShell = "shell"
	updateEnvKey    = "ALFRED_TLDR_UPDATE"
	// Note the key is also defined in workflow environment variable
	autoUpdateEnvKey = "ALFRED_TLDR_AUTO_UPDATE"
)

var awf = alfred.NewWorkflow(
	alfred.WithMaxResults(30),
	alfred.WithOutStream(outStream),
	alfred.WithLogStream(errStream),
)

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

func isAutoUpdateEnabled() bool {
	sv := os.Getenv(autoUpdateEnvKey)
	bv, err := strconv.ParseBool(sv)
	if err != nil {
		return false
	}
	return bv
}

func fatal(err error) {
	awf.Fatal("Fatal errors occur", fmt.Sprint(err))
}
