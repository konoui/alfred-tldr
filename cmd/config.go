package cmd

import (
	"github.com/konoui/alfred-tldr/pkg/tldr"
	"github.com/konoui/go-alfred"
)

type envs struct {
	formatFunc                       func(string) string
	modKeyOpenURL                    alfred.ModKey
	isUpdateWorkflowRecommendEnabled bool
	isUpdateDBRecommendEnabled       bool
}

type config struct {
	platform       tldr.Platform
	language       string
	update         bool
	updateWorkflow bool
	confirm        bool
	fuzzy          bool
	version        bool
	tldrClient     *tldr.Tldr
	fromEnv        envs
}

func newConfig() *config {
	cfg := new(config)
	cfg.fromEnv.formatFunc = getCommandFormatFunc()
	cfg.fromEnv.modKeyOpenURL = getModKeyOpenURL()
	cfg.fromEnv.isUpdateDBRecommendEnabled = isUpdateDBRecommendEnabled()
	cfg.fromEnv.isUpdateWorkflowRecommendEnabled = isUpdateWorkflowRecommendEnabled()
	return cfg
}
