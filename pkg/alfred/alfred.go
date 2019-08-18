package alfred

import (
	"encoding/json"
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

// Alfred Script Filter JSON Format
type Alfred struct {
	rerun     Rerun
	variables Variables
	items     Items
}

// New creates a new Workflow
func New() *Alfred {
	return &Alfred{}
}

// Append a new Item to alfred
func (a *Alfred) Append(item Item) {
	a.items = append(a.items, item)
}

type out struct {
	Rerun     Rerun     `json:"rerun,omitempty"`
	Variables Variables `json:"variables,omitempty"`
	Items     Items     `json:"items"`
}

const errJSON = `{
    "items": [{
	    "title": "no matching entry",
	    "subtitle": "Try a different query?"
    }]
}`

// Marshal as Workflow as Json
func (a *Alfred) Marshal() string {
	if len(a.items) == 0 {
		return errJSON
	}
	obj, err := json.MarshalIndent(
		out{
			Rerun:     a.rerun,
			Variables: a.variables,
			Items:     a.items,
		}, "", "	")
	if err != nil {
		return errJSON
	}

	return string(obj)
}
