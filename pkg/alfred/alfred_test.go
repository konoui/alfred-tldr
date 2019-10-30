package alfred

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewScriptFilter(t *testing.T) {
	tests := []struct {
		description string
		want        *ScriptFilter
	}{
		{
			description: "create new workflow",
			want:        &ScriptFilter{},
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
			description: "create new workflow",
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

			v, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			wf := NewScriptFilter()
			for _, item := range tt.items {
				wf.Append(item)
			}

			res, err := wf.Marshal()
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error want: got: %+v", err.Error())
			}

			want := string(v)
			got := string(res)
			ret, err := AreEqualJSON(want, got)
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error want: got: %+v", err.Error())
			}
			if !ret {
				t.Errorf("unexpected response: want: \n%+v, got: \n%+v", want, got)
			}
		})
	}
}

func TestNewWorkflow(t *testing.T) {
	tests := []struct {
		description string
		want        Workflow
	}{
		{
			description: "create new workflow",
			want: Workflow{
				std:  NewScriptFilter(),
				warn: NewScriptFilter(),
				err:  NewScriptFilter(),
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

func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
