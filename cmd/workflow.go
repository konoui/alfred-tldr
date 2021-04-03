package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/konoui/go-alfred"
	"github.com/konoui/go-alfred/update"
)

// decide next action for workflow filter
const (
	nextActionKey   = "nextAction"
	nextActionCopy  = "copy"
	nextActionShell = "shell"
	// Note the key is also defined in workflow environment variable
	updateDBRecommendEnvKey       = "ALFRED_TLDR_DB_UPDATE_RECOMMEND"
	updateWorkflowRecommendEnvKey = "ALFRED_TLDR_WORKFLOW_UPDATE_RECOMMEND"
)

var awf *alfred.Workflow

func init() {
	awf = alfred.NewWorkflow(
		alfred.WithMaxResults(30),
		alfred.WithGitHubUpdater(
			"konoui",
			"alfred-tldr",
			version,
			update.WithCheckInterval(twoWeeks),
		),
	)
	awf.SetOut(outStream)
	awf.SetLog(errStream)
}

func isUpdateDBRecommendEnabled() bool {
	return parseBool(updateDBRecommendEnvKey)
}

func isUpdateWorkflowRecommendEnabled() bool {
	return parseBool(updateWorkflowRecommendEnvKey)
}

func parseBool(key string) bool {
	sv := os.Getenv(key)
	bv, err := strconv.ParseBool(sv)
	if err != nil {
		return false
	}
	return bv
}

func fatal(err error) {
	awf.Fatal("a fatal error occurred", fmt.Sprint(err))
}
