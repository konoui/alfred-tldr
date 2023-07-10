package tldr

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sahilm/fuzzy"
)

// CmdInfo contains name, platform, language
type CmdInfo struct {
	Name      string     `json:"name"`
	Platforms []Platform `json:"platform"`
	Languages []string   `json:"language"`
}

// Cmds a slice of CmdInfo
type Cmds []*CmdInfo

// CmdsIndex structure of index.json
type CmdsIndex struct {
	Commands Cmds `json:"commands"`
}

func (c Cmds) String(i int) string {
	return c[i].Name
}

// Len return length of Commands for fuzzy interface
func (c Cmds) Len() int {
	return len(c)
}

// Search fuzzy search commands by query.
func (c Cmds) Search(args []string) Cmds {
	query := commandWithHyphen(args)
	results := fuzzy.FindFrom(query, c)
	cmds := make(Cmds, results.Len())
	for i, r := range results {
		if guessHyphenCommand(args) {
			// Note: if original args contain hyphen, it regards as command with hyphen.
			// e.g.)  apt-key -> apt-key not apt key
			cmds[i] = c[r.Index]
		} else {
			// Note: if original args does not have hyphen,	replacing highfun with space after search.
			// e.g.) git-checkout -> git checkout
			cmdName := strings.ReplaceAll(c[r.Index].Name, "-", " ")
			c[r.Index].Name = cmdName
			cmds[i] = c[r.Index]
		}
	}
	return cmds
}

func commandWithHyphen(args []string) (arg string) {
	// e.g.) git checkout -> git-checkout filename is git-checkout.md
	arg = strings.Join(args, "-")
	// e.g.)  apt- key -> apt-key not apt--key
	arg = strings.ReplaceAll(arg, "--", "-")
	return
}

func guessHyphenCommand(args []string) bool {
	for _, v := range args {
		if strings.Contains(v, "-") {
			return true
		}
	}
	return false
}

// LoadIndexFile load command index file
func (t *Tldr) LoadIndexFile() (*CmdsIndex, error) {
	f, err := os.Open(t.indexFilePath())
	if err != nil {
		return nil, fmt.Errorf("failed to open a index file: %w", err)
	}
	defer f.Close()

	cmdIndex := &CmdsIndex{}
	return cmdIndex, json.NewDecoder(f).Decode(cmdIndex)
}
