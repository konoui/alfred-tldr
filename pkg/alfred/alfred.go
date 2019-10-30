package alfred

import (
	"encoding/json"
	"fmt"
)

// see https://www.alfredapp.com/help/workflows/inputs/script-filter/json/

// Icon displayed in the result row
type Icon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

// Mod element gives you control over how the modifier keys react
type Mod struct {
	Variables map[string]string `json:"variables,omitempty"`
	Valid     *bool             `json:"valid,omitempty"`
	Arg       string            `json:"arg,omitempty"`
	Subtitle  string            `json:"subtitle,omitempty"`
	Icon      *Icon             `json:"icon,omitempty"`
}

// Text element defines the text the user will get when copying the selected result row
type Text struct {
	Copy      string `json:"copy,omitempty"`
	Largetype string `json:"largetype,omitempty"`
}

// Item a workflow object
type Item struct {
	Variables    map[string]string `json:"variables,omitempty"`
	UID          string            `json:"uid,omitempty"`
	Title        string            `json:"title"`
	Subtitle     string            `json:"subtitle,omitempty"`
	Arg          string            `json:"arg,omitempty"`
	Icon         *Icon             `json:"icon,omitempty"`
	Autocomplete string            `json:"autocomplete,omitempty"`
	Type         string            `json:"type,omitempty"`
	Valid        *bool             `json:"valid,omitempty"`
	Match        string            `json:"match,omitempty"`
	Mods         map[string]Mod    `json:"mods,omitempty"`
	Text         *Text             `json:"text,omitempty"`
	QuicklookURL string            `json:"quicklookurl,omitempty"`
}

// Rerun re-run automatically after an interval
type Rerun float64

// Variables passed out of the script filter within a `variables` object.
type Variables map[string]string

// Items array of `item`
type Items []Item

// ScriptFilter JSON Format
type ScriptFilter struct {
	rerun     Rerun
	variables Variables
	items     Items
}

// NewScriptFilter creates a new ScriptFilter
func NewScriptFilter() *ScriptFilter {
	return &ScriptFilter{}
}

// Append a new Item to Items
func (w *ScriptFilter) append(item Item) {
	w.items = append(w.items, item)
}

type out struct {
	Rerun     Rerun     `json:"rerun,omitempty"`
	Variables Variables `json:"variables,omitempty"`
	Items     Items     `json:"items"`
}

// Marshal as ScriptFilter as Json
func (w *ScriptFilter) Marshal() ([]byte, error) {
	return json.MarshalIndent(
		out{
			Rerun:     w.rerun,
			Variables: w.variables,
			Items:     w.items,
		}, "", "	")
}

type wfKey string

const (
	std  wfKey = "std"
	warn wfKey = "warn"
	err  wfKey = "err"
)

// Workflow has standard and error Workflow
type Workflow map[wfKey]*ScriptFilter

// NewWorkflow creates a new AlfredWorkflow
func NewWorkflow() Workflow {
	return Workflow{
		std:  NewScriptFilter(),
		warn: NewScriptFilter(),
		err:  NewScriptFilter(),
	}
}

// Append a new Item to standard ScriptFilter
func (awf Workflow) Append(item Item) {
	awf[std].append(item)
}

// Warning append a new Item to error ScriptFilter
func (awf Workflow) Warning(title, subtitle string) {
	awf[warn].append(
		Item{
			Title:    title,
			Subtitle: subtitle,
		})
}

// Marshal WorkFlow results
func (awf Workflow) Marshal() []byte {
	wf := awf[std]
	if len(wf.items) == 0 {
		warnRes, err := awf[warn].Marshal()
		if err != nil {
			return []byte("")
		}
		return warnRes
	}

	res, err := wf.Marshal()
	if err != nil {
		// FIXME should create a new workflow and use it to set Error
		awf.Warning("An Error Occurs... ", fmt.Sprintf("items length: %d, items: %v", len(wf.items), wf.items))
		warnRes, err := awf[warn].Marshal()
		if err != nil {
			return []byte("")
		}
		return warnRes
	}

	return res
}
