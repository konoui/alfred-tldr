package tldr

import (
	"os"
	"strings"
)

func getLanguages(optionLang string) (priorities []string) {
	if optionLang != "" {
		return []string{optionLang}
	}

	// see https://github.com/tldr-pages/tldr/blob/master/CLIENT-SPECIFICATION.md#language
	langCode := getLanguageCode(os.Getenv("LANG"))
	if langCode == "" {
		return []string{languageCodeEN}
	}
	defer func() {
		if !contains(langCode, priorities) {
			priorities = append(priorities, langCode)
		}
		if !contains(languageCodeEN, priorities) {
			priorities = append(priorities, languageCodeEN)
		}
	}()

	languageEnv := os.Getenv("LANGUAGE")
	if languageEnv == "" {
		return
	}

	for _, language := range strings.Split(languageEnv, ":") {
		code := getLanguageCode(language)
		if contains(code, priorities) {
			continue
		}
		priorities = append(priorities, code)
	}
	return
}

func contains(target string, list []string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func getLanguageCode(language string) string {
	code := strings.SplitN(language, ".", 2)[0]
	if code == "C" || code == "POSIX" {
		return ""
	}

	for _, v := range []string{"pt_PT", "pt_BR", "zh_TW"} {
		if v == code {
			return code
		}
	}

	if code == "pt" {
		return "pt_PT"
	}

	return strings.SplitN(code, "_", 2)[0]
}

func getLangDir(lang string) string {
	if lang == languageCodeEN {
		return "pages"
	}

	return "pages" + "." + lang
}
