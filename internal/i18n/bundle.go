package i18n

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

// Supported languages
const (
	TR = "tr"
	EN = "en"
)

func InitBundle(languagesPath string) {
	bundle = i18n.NewBundle(language.Turkish)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	var languages = []string{
		languagesPath + "/en.json",
		languagesPath + "/tr.json",
	}

	for _, language := range languages {
		bundle.MustLoadMessageFile(language)
	}

}
