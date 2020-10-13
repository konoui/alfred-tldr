package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func getStringFlag(cmd *cobra.Command, name string) string {
	v, err := cmd.PersistentFlags().GetString(name)
	if err != nil {
		return ""
	}
	return v
}

func getBoolFlag(cmd *cobra.Command, name string) bool {
	v, err := cmd.PersistentFlags().GetBool(name)
	if err != nil {
		return false
	}
	return v
}

func makeUsageMap(cmd *cobra.Command) (m map[string]string) {
	pf := cmd.Flag(platformFlag)
	uf := cmd.Flag(updateFlag)
	m = make(map[string]string, 2)
	m[pf.Name] = makeDescription(pf)
	m[uf.Name] = makeDescription(uf)
	return
}

func makeDescription(p *pflag.Flag) string {
	return fmt.Sprintf("Usage: -%s, --%s %s", p.Shorthand, p.Name, p.Usage)
}
