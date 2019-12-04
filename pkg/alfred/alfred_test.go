package alfred

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewScriptFilter(t *testing.T) {
	tests := []struct {
		description string
		want        ScriptFilter
	}{
		{
			description: "create new workflow",
			want:        ScriptFilter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := NewScriptFilter()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("unexpected response: want: %+v, got: %+v", tt.want, got)
			}

		})
	}
}

func TestScriptFilterMarshal(t *testing.T) {
	tests := []struct {
		description string
		expectErr   bool
		filepath    string
		items       Items
	}{
		{
			description: "create new scriptfilter",
			filepath:    "./test_scriptfilter_marshal.json",
			items: Items{
				Item{
					Title:    "title1",
					Subtitle: "subtitle1",
				},
				Item{
					Title:    "title2",
					Subtitle: "subtitle2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			f, err := os.Open(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			want, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			wf := NewScriptFilter()
			for _, item := range tt.items {
				wf.Append(item)
			}

			got, err := wf.Marshal()
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error want: got: %+v", err.Error())
			}

			if !EqualJSON(want, got) {
				t.Errorf("unexpected response: want: \n%+v, got: \n%+v", string(want), string(got))
			}
		})
	}
}

func TestNewWorkflow(t *testing.T) {
	tests := []struct {
		description string
		want        *Workflow
	}{
		{
			description: "create new workflow",
			want: &Workflow{
				std:  NewScriptFilter(),
				warn: NewScriptFilter(),
				err:  NewScriptFilter(),
				streams: streams{
					out: os.Stdout,
					err: os.Stdout,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := NewWorkflow()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("unexpected response: want: %+v, got: %+v", tt.want, got)
			}

		})
	}
}

func TestWorfkflowMarshal(t *testing.T) {
	tests := []struct {
		description string
		filepath    string
		items       Items
		emptyItem   Item
	}{
		{
			description: "output standard items",
			filepath:    "./test_scriptfilter_marshal.json",
			items: Items{
				Item{
					Title:    "title1",
					Subtitle: "subtitle1",
				},
				Item{
					Title:    "title2",
					Subtitle: "subtitle2",
				},
			},
			emptyItem: Item{
				Title:    "emptyTitle1",
				Subtitle: "emptySubtitle",
			},
		},
		{
			description: "output empty warning",
			filepath:    "./test_scriptfilter_empty_warning_marshal.json",
			items:       Items{},
			emptyItem: Item{
				Title:    "emptyTitle1",
				Subtitle: "emptySubtitle1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			f, err := os.Open(tt.filepath)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			want, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			awf := NewWorkflow()
			awf.EmptyWarning(tt.emptyItem.Title, tt.emptyItem.Subtitle)
			for _, item := range tt.items {
				awf.Append(item)
			}

			got := awf.Marshal()
			if !EqualJSON(want, got) {
				t.Errorf("unexpected response: want: \n%+v, got: \n%+v", string(want), string(got))
			}
		})
	}
}

func EqualJSON(a, b []byte) bool {
	var ao interface{}
	var bo interface{}

	if err := json.Unmarshal(a, &ao); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &bo); err != nil {
		return false
	}

	return reflect.DeepEqual(ao, bo)
}
