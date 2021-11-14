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

type Config struct {
	platform       tldr.Platform
	language       string
	update         bool
	updateWorkflow bool
	confirm        bool
	fuzzy          bool
	version        bool
	tldrClient     *tldr.Tldr
	fromEnv        envs
	opts           []tldr.Option
}

func NewConfig() *Config {
	cfg := new(Config)
	cfg.fromEnv.formatFunc = getCommandFormatFunc()
	cfg.fromEnv.modKeyOpenURL = getModKeyOpenURL()
	cfg.fromEnv.isUpdateDBRecommendEnabled = isUpdateDBRecommendEnabled()
	cfg.fromEnv.isUpdateWorkflowRecommendEnabled = isUpdateWorkflowRecommendEnabled()
	return cfg
}
