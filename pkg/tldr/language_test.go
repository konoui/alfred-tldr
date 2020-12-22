package tldr

import (
	"os"
	"reflect"
	"testing"
)

func Test_getLanguages(t *testing.T) {
	type args struct {
		optionLang  string
		langEnv     string
		languageEnv string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "specify lang flag",
			want: []string{"nl"},
			args: args{
				optionLang:  "nl",
				langEnv:     "en",
				languageEnv: "en",
			},
		},
		{
			name: "from LANG env",
			want: []string{"nl", "en"},
			args: args{
				langEnv: "nl",
			},
		},
		{
			name: "from both env",
			want: []string{"pt_BR", "nl", "en"},
			args: args{
				langEnv:     "nl",
				languageEnv: "pt_BR",
			},
		},
		{
			name: "should not duplicate en",
			want: []string{"en"},
			args: args{
				langEnv: "en",
			},
		},
		{
			name: "should not duplicate en with both env",
			want: []string{"en"},
			args: args{
				langEnv:     "en",
				languageEnv: "en",
			},
		},
		{
			name: "should not duplicate en and add language env",
			want: []string{"nl", "en"},
			args: args{
				langEnv:     "en",
				languageEnv: "nl",
			},
		},
		{
			name: "LANGAGE env has some values and duplications",
			want: []string{"pt_BR", "nl", "ja", "it", "pt_PT", "en"},
			args: args{
				langEnv:     "nl",
				languageEnv: "pt_BR:nl:ja:it:pt:pt",
			},
		},
		{
			name: "ignore LANGUAGE env if LANG env is empty",
			want: []string{"en"},
			args: args{
				languageEnv: "pt_BR",
			},
		},
		{
			name: "regard POSIX as empty",
			want: []string{"en"},
			args: args{
				langEnv:     "POSIX",
				languageEnv: "pt_BR",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			if err := os.Setenv("LANG", tt.args.langEnv); err != nil {
				t.Fatal(err)
			}
			if err := os.Setenv("LANGUAGE", tt.args.languageEnv); err != nil {
				t.Fatal(err)
			}
			// cleanup
			defer func() {
				if err := os.Unsetenv("LANG"); err != nil {
					t.Fatal(err)
				}
				if err := os.Unsetenv("LANGUAGE"); err != nil {
					t.Fatal(err)
				}
			}()

			if got := getLanguages(tt.args.optionLang); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLanguages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLanguageCode(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "split with dot and underbar",
			input: "ja_JP.UTF-8",
			want:  "ja",
		},
		{
			name:  "not split with underbar but split dot",
			input: "pt_PT.UTF-8",
			want:  "pt_PT",
		},
		{
			name:  "should filter value",
			input: "POSIX",
			want:  "",
		},
		{
			name:  "should filter value",
			input: "C",
			want:  "",
		},
		{
			name:  "not lang value return the input value",
			input: "invalid",
			want:  "invalid",
		},
		{
			name:  "no input",
			input: "",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLanguageCode(tt.input); got != tt.want {
				t.Errorf("getLanguageCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
