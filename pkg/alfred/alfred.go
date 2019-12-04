package alfred

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// see https://www.alfredapp.com/help/workflows/inputs/script-filter/json/

// Icon displayed in the result row
type Icon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

// ModKey is a mod key pressed by the user to run an alternate
type ModKey string

// Valid attribute to mark if the result is valid based on the modifier selection and set a different arg to be passed out if actioned with the modifier.
const (
	ModCmd   ModKey = "cmd"   // Alternate action for ⌘↩
	ModAlt   ModKey = "alt"   // Alternate action for ⌥↩
	ModOpt   ModKey = "alt"   // Synonym for ModAlt
	ModCtrl  ModKey = "ctrl"  // Alternate action for ^↩
	ModShift ModKey = "shift" // Alternate action for ⇧↩
	ModFn    ModKey = "fn"    // Alternate action for fn↩
)

// Mod element gives you control over how the modifier keys react
type Mod struct {
	Variables map[ModKey]string `json:"variables,omitempty"`
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

type out struct {
	Rerun     Rerun     `json:"rerun,omitempty"`
	Variables Variables `json:"variables,omitempty"`
	Items     Items     `json:"items"`
}

// NewScriptFilter creates a new ScriptFilter
func NewScriptFilter() ScriptFilter {
	return ScriptFilter{}
}

// Append a new Item to Items
func (w *ScriptFilter) Append(item Item) {
	w.items = append(w.items, item)
}

// Marshal ScriptFilter as Json
func (w *ScriptFilter) Marshal() ([]byte, error) {
	return json.MarshalIndent(
		out{
			Rerun:     w.rerun,
			Variables: w.variables,
			Items:     w.items,
		}, "", "	")
}

// Workflow is map of ScriptFilters
type Workflow struct {
	std     ScriptFilter
	warn    ScriptFilter
	err     ScriptFilter
	streams streams
}

type streams struct {
	out io.Writer
	err io.Writer
}

// SetStdStream redirect stdout to s
func (awf *Workflow) SetStdStream(s io.Writer) {
	awf.streams.out = s
}

// SetErrStream redirect stderr to s
func (awf *Workflow) SetErrStream(s io.Writer) {
	awf.streams.out = s
}

// NewWorkflow has simple ScriptFilter api
func NewWorkflow() *Workflow {
	return &Workflow{
		std:  NewScriptFilter(),
		warn: NewScriptFilter(),
		err:  NewScriptFilter(),
		streams: streams{
			out: os.Stdout,
			err: os.Stdout,
		},
	}
}

// Append a new Item to standard ScriptFilter
func (awf *Workflow) Append(item Item) {
	awf.std.Append(item)
}

// EmptyWarning create a new Item to Marshal　when there are no standard items
func (awf *Workflow) EmptyWarning(title, subtitle string) {
	awf.warn = NewScriptFilter()
	awf.warn.Append(
		Item{
			Title:    title,
			Subtitle: subtitle,
		})
}

// Warning append a new Item to error ScriptFilter
func (awf *Workflow) error(title, subtitle string) {
	awf.err = NewScriptFilter()
	awf.err.Append(
		Item{
			Title:    title,
			Subtitle: subtitle,
		})
}

// Marshal WorkFlow results
func (awf *Workflow) Marshal() []byte {
	wf := awf.std
	if len(wf.items) == 0 {
		warnRes, err := awf.warn.Marshal()
		if err != nil {
			return []byte("")
		}
		return warnRes
	}

	res, err := wf.Marshal()
	if err != nil {
		awf.error(fmt.Sprintf("An Error Occurs: %s", err.Error()), fmt.Sprintf("items length: %d, items: %v", len(wf.items), wf.items))
		errRes, err := awf.err.Marshal()
		if err != nil {
			return []byte("")
		}
		return errRes
	}

	return res
}

// Fatal output error to io stream
func (awf *Workflow) Fatal(title, subtitle string) {
	awf.error(title, subtitle)
	res, _ := awf.err.Marshal()
	fmt.Fprintln(awf.streams.err, string(res))
}

// Output to io stream
func (awf *Workflow) Output() {
	res := awf.Marshal()
	fmt.Fprintln(awf.streams.out, string(res))
}
