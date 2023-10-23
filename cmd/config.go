package cmd

import (
	"context"
	"fmt"
	"path/filepath"

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
	fromEnv        envs
	tldrOpts       []tldr.Option
	update         bool
	updateWorkflow bool
	confirm        bool
	fuzzy          bool
	version        bool
}

func NewConfig() *Config {
	cfg := new(Config)
	cfg.fromEnv.formatFunc = getCommandFormatFunc()
	cfg.fromEnv.modKeyOpenURL = getModKeyOpenURL()
	cfg.fromEnv.isUpdateDBRecommendEnabled = isUpdateDBRecommendEnabled()
	cfg.fromEnv.isUpdateWorkflowRecommendEnabled = isUpdateWorkflowRecommendEnabled()
	return cfg
}

func (cfg *Config) setPlatform(ptString string) error {
	platforms := []tldr.Platform{
		tldr.PlatformCommon,
		tldr.PlatformLinux,
		tldr.PlatformOSX,
		tldr.PlatformWindows,
		tldr.PlatformSunos,
	}
	for _, pt := range platforms {
		if ptString == pt.String() {
			cfg.platform = pt
			return nil
		}
	}
	return fmt.Errorf("%s is unsupported platform", ptString)
}

func newTldrClient(cfg *Config, awf *alfred.Workflow) (*tldr.Tldr, error) {
	path := filepath.Join(awf.GetDataDir(), "data")

	opts := append([]tldr.Option{
		tldr.WithPlatform(cfg.platform),
		tldr.WithLanguage(cfg.language),
	}, cfg.tldrOpts...)

	tldrClient := tldr.New(path, opts...)
	ctx, cancel := context.WithTimeout(context.Background(), updateDBTimeout)
	defer cancel()
	if err := tldrClient.OnInitialize(ctx); err != nil {
		return nil, err
	}
	return tldrClient, nil
}
