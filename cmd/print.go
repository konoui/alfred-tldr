package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/konoui/alfred-tldr/pkg/tldr"
	"github.com/konoui/go-alfred"
)

func (cfg *Config) printPage(cmds []string) error {
	// insert update recommendation first
	ctx, cancel := context.WithTimeout(context.Background(), updateWorkflowCheckTimeout)
	defer cancel()
	if cfg.fromEnv.isUpdateWorkflowRecommendEnabled && awf.Updater().IsNewVersionAvailable(ctx) {
		awf.Append(
			alfred.NewItem().
				Title("Newer tldr wrokflow is available!").
				Subtitle("Please Enter!").
				Arg(fmt.Sprintf("--%s --%s", updateWorkflowFlag, confirmFlag)).
				Variable(nextActionKey, nextActionShell).
				Icon(awf.Assets().IconAlertNote()),
		)
	}

	if cfg.fromEnv.isUpdateDBRecommendEnabled && cfg.tldrClient.Expired(twoWeeks) {
		awf.Append(
			alfred.NewItem().
				Title("Tldr database is older than 2 weeks").
				Subtitle("Please Enter!").
				Arg(fmt.Sprintf("--%s --%s", longUpdateFlag, confirmFlag)).
				Icon(awf.Assets().IconAlertNote()).
				Variable(nextActionKey, nextActionShell),
		)
	}

	// no input case
	if len(cmds) == 0 {
		awf.Append(
			alfred.NewItem().
				Title("Please input a command").
				Subtitle("e.g.) tldr tar e.g.) tldr --help").
				Valid(false),
		).Output()
		return nil
	}

	awf.SetEmptyWarning("No matching query", "Try a different query")
	p, err := cfg.tldrClient.FindPage(cmds)
	if err != nil {
		if errors.Is(err, tldr.ErrNotFoundPage) {
			if cfg.language != "" {
				awf.Clear().SetEmptyWarning(
					"Not found the command in selected language",
					"Try not to specify language option",
				).Output()
				return nil
			}
			if cfg.fuzzy {
				// list suggestions
				return cfg.printFuzzyPages(cmds)
			}
			awf.Output()
			return nil
		}
		return err
	}

	awf.Append(
		makeDescriptionItem(p, cfg.fromEnv.modKeyOpenURL),
	)
	for _, cmd := range p.CmdExamples {
		command := cfg.fromEnv.formatFunc(cmd.Cmd)
		awf.Append(
			alfred.NewItem().
				Title(command).
				Subtitle(cmd.Description).
				Arg(command),
		).Variable(nextActionKey, nextActionCopy)
	}

	awf.Output()
	return nil
}

func makeDescriptionItem(p *tldr.Page, modKey alfred.ModKey) *alfred.Item {
	// Note descriptions has one line at least
	// see: https://github.com/tldr-pages/tldr/blob/master/contributing-guides/style-guide.md
	title := p.CmdDescriptions[0]
	subtitle := ""
	if len(p.CmdDescriptions) >= 2 {
		subtitle = p.CmdDescriptions[1]
	}

	openMod := alfred.NewMod()
	u, err := parseDetailURL(p.CmdDescriptions)
	if err != nil {
		awf.Logger().Warnln(err)
		openMod.
			Valid(false).
			Subtitle("no action")
	} else {
		openMod.
			Arg(u).
			Subtitle("open more information url").
			Variable(nextActionKey, nextActionOpenURL)
	}

	return alfred.NewItem().
		Title(title).
		Subtitle(subtitle).
		Valid(false).
		Icon(
			alfred.NewIcon().
				Path("description.png"),
		).
		Mod(modKey, openMod)
}

func (cfg *Config) printFuzzyPages(cmds []string) error {
	index, err := cfg.tldrClient.LoadIndexFile()
	if err != nil {
		return err
	}

	suggestions := index.Commands.Search(cmds)
	for _, cmd := range suggestions {
		complete := cmd.Name
		pt := choicePlatform(cmd.Platforms, cfg.platform)
		if pt != tldr.PlatformCommon && pt != defaultPlatform {
			complete = fmt.Sprintf("-%s %s %s",
				platformFlag,
				pt,
				cmd.Name,
			)
		}
		awf.Append(
			alfred.NewItem().
				Title(cmd.Name).
				Subtitle(fmt.Sprintf("Platforms: %s", fmt.Sprintf("%s", cmd.Platforms))).
				Valid(false).
				Autocomplete(complete).
				Icon(
					alfred.NewIcon().
						Path("candidate.png"),
				),
		)
	}

	awf.Output()
	return nil
}

func (cfg *Config) printVersion(v, r string) (_ error) {
	title := fmt.Sprintf("alfred-tldr %v(%s)", v, r)
	awf.Append(
		alfred.NewItem().Title(title),
	).Output()
	return
}

func choicePlatform(pts []tldr.Platform, selected tldr.Platform) tldr.Platform {
	if len(pts) >= 2 {
		// if there are more than two platforms,
		// priority are follows
		// selected pt, common, others
		for _, pt := range pts {
			if pt == selected {
				return selected
			}
		}

		for _, pt := range pts {
			if pt == tldr.PlatformCommon {
				return tldr.PlatformCommon
			}
		}
	}

	// Note: unexpected case, we assume one platform in pts at least but if empty, return common
	if len(pts) == 0 {
		return tldr.PlatformCommon
	}

	return pts[0]
}

// see format https://github.com/tldr-pages/tldr/blob/main/contributing-guides/style-guide.md
// > Short, snappy description.
// > Preferably one line; two are acceptable if necessary.
// > More information: <https://example.com>.
func parseDetailURL(descriptions []string) (string, error) {
	d := descriptions[len(descriptions)-1]
	for _, scheme := range []string{"https://", "http://"} {
		lastIndex := strings.LastIndex(d, ">.")
		firstIndex := strings.Index(d, "<"+scheme)
		if lastIndex < 0 || firstIndex < 0 {
			continue
		}

		if lastIndex < firstIndex {
			return "", fmt.Errorf("found URL in descriptions but something wrong %s", d)
		}

		u, err := url.Parse(d[firstIndex+1 : lastIndex])
		if err != nil {
			return "", fmt.Errorf("found URL in descriptions but failed to parse: %w", err)
		}
		return u.String(), nil
	}
	return "", fmt.Errorf("not found URL in descriptions")
}
