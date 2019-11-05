package tldr

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

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
		cmds = append(cmds, c[r.Index])
	}

	return cmds
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
