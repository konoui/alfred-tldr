package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/konoui/go-alfred"
)

// decide next action for workflow filter
const (
	nextActionKey     = "nextAction"
	nextActionCopy    = "copy"
	nextActionShell   = "shell"
	nextActionOpenURL = "openURL"
	// Note the key is also defined in workflow environment variable
	envKeyUpdateDBRecommendation       = "TLDR_DB_UPDATE_RECOMMENDATION"
	envKeyUpdateWorkflowRecommendation = "TLDR_WORKFLOW_UPDATE_RECOMMENDATION"
	envKeyUpdateWorkflowIntervalDays   = "TLDR_WORKFLOW_UPDATE_INTERVAL_DAYS"
	envKeyCommandFormat                = "TLDR_COMMAND_FORMAT"
	envKeyOpenURLMod                   = "TLDR_MOD_KEY_OPEN_URL"
)

var awf = alfred.NewWorkflow(
	alfred.WithMaxResults(30),
	alfred.WithGitHubUpdater(
		"konoui",
		"alfred-tldr",
		version,
		getUpdateWorkflowInterval(twoWeeks),
	),
	alfred.WithOutWriter(os.Stdout),
	alfred.WithLogWriter(os.Stderr),
)

func getModKeyOpenURL() alfred.ModKey {
	v := os.Getenv(envKeyOpenURLMod)
	switch v {
	case "alt":
		return alfred.ModAlt
	case "cmd":
		return alfred.ModCmd
	case "ctrl":
		return alfred.ModCtrl
	case "fn":
		return alfred.ModFn
	case "shift":
		return alfred.ModShift
	default:
		// TODO should be empty
		return alfred.ModCmd
	}
}

func getCommandFormatFunc() func(string) string {
	v := os.Getenv(envKeyCommandFormat)
	switch v {
	case "original":
		return func(cmd string) string {
			return cmd
		}
	case "remove":
		return func(cmd string) string {
			tmp := strings.ReplaceAll(cmd, "{{", "")
			return strings.ReplaceAll(tmp, "}}", "")
		}
	case "uppercase":
		var wordRe = regexp.MustCompile("{{.+?}}")
		return func(cmd string) string {
			return wordRe.ReplaceAllStringFunc(cmd, func(w string) string {
				upper := strings.ToUpper(w)
				start := len("{{")
				end := len("}}")
				return upper[start : len(upper)-end]
			})
		}
	case "single":
	default:
	}

	return func(cmd string) string {
		tmp := strings.ReplaceAll(cmd, "{{", "{")
		return strings.ReplaceAll(tmp, "}}", "}")
	}
}

func isUpdateDBRecommendEnabled() bool {
	return parseBool(envKeyUpdateDBRecommendation)
}

func isUpdateWorkflowRecommendEnabled() bool {
	return parseBool(envKeyUpdateWorkflowRecommendation)
}

func getUpdateWorkflowInterval(defaultInterval time.Duration) time.Duration {
	v := os.Getenv(envKeyUpdateWorkflowIntervalDays)
	fv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return defaultInterval
	}
	tv := time.Duration(fv) * 24 * time.Hour
	return tv
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
