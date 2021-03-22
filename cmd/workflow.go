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
	nextActionShell = "shell"
	// Note the key is also defined in workflow environment variable
	autoUpdateEnvKey = "ALFRED_TLDR_AUTO_UPDATE"
)

var awf *alfred.Workflow

func init() {
	awf = alfred.NewWorkflow(
		alfred.WithMaxResults(30),
	)
	awf.SetOut(outStream)
	awf.SetLog(errStream)
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
	awf.Fatal("a fatal error occurred", fmt.Sprint(err))
}
