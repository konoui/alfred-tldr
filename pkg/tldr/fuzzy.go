package tldr

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sahilm/fuzzy"
)

// CmdInfo contains name, platform, language
type CmdInfo struct {
	Name     string   `json:"name"`
	Platform []string `json:"platform"`
	Language []string `json:"language"`
}

// Commands a slice of CmdInfo
type Commands []*CmdInfo

// CmdIndex structure of index.json
type CmdIndex struct {
	Commands Commands `json:"commands"`
}

func (c Commands) String(i int) string {
	return c[i].Name
}

// Len return length of Commands for fuzzy interface
func (c Commands) Len() int {
	return len(c)
}

// Filter fuzzy search commands by query
func (c Commands) Filter(query string) Commands {
	cmds := Commands{}
	results := fuzzy.FindFrom(query, c)
	for _, r := range results {
		// Note: replace highfun with space as command name in index file joined highfun
		// e.g.) git-checkout -> git checkout
		cmdName := strings.Replace(c[r.Index].Name, "-", " ", -1)
		c[r.Index].Name = cmdName
		cmds = append(cmds, c[r.Index])
	}

	return cmds
}

// Search fuzzy search commands by query. This is wrapped Filtter
func (c Commands) Search(args []string) Commands {
	// Note: We should replace space with highfun as a index file format is joined with highfun
	// e.g.) git checkout -> git-checkout.md
	query := strings.Join(args, "-")
	return c.Filter(query)
}

// LoadIndexFile load command index file
func (t *Tldr) LoadIndexFile() (*CmdIndex, error) {
	f, err := os.Open(filepath.Join(t.path, t.indexFile))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	v, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	cmdIndex := &CmdIndex{}
	return cmdIndex, json.Unmarshal(v, cmdIndex)
}
