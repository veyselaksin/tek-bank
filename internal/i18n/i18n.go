package i18n

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"tek-bank/cmd/config"
)

type localize struct {
	messageId    string
	templateData map[string]string
	pluralCount  int
}

type localizeBuilder struct {
	localize *localize
}

func (b *localizeBuilder) WithTemplateData(templateData map[string]string) *localizeBuilder {
	b.localize.templateData = templateData
	return b
}

func (b *localizeBuilder) WithPluralCount(pluralCount int) *localizeBuilder {
	b.localize.pluralCount = pluralCount
	return b
}

func (b *localizeBuilder) Build(loc *i18n.Localizer) string {

	var pluralCount int
	if b.localize.pluralCount == 0 && b.localize.templateData != nil {
		pluralCount = 1
	} else {
		pluralCount = b.localize.pluralCount
	}

	message := loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    b.localize.messageId,
		TemplateData: b.localize.templateData,
		PluralCount:  pluralCount,
	})

	return message
}

func (b *localizeBuilder) BuildWithContext(c *fiber.Ctx) string {

	var pluralCount int
	if b.localize.pluralCount == 0 && b.localize.templateData != nil {
		pluralCount = 1
	} else {
		pluralCount = b.localize.pluralCount
	}

	lang := config.GetLanguage(c)
	loc := i18n.NewLocalizer(bundle, lang)

	message := loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    b.localize.messageId,
		TemplateData: b.localize.templateData,
		PluralCount:  pluralCount,
	})

	return message
}

func (b *localizeBuilder) BuildWithLanguage(lang string) string {

	var pluralCount int
	if b.localize.pluralCount == 0 && b.localize.templateData != nil {
		pluralCount = 1
	} else {
		pluralCount = b.localize.pluralCount
	}

	loc := i18n.NewLocalizer(bundle, lang)

	message := loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    b.localize.messageId,
		TemplateData: b.localize.templateData,
		PluralCount:  pluralCount,
	})

	return message
}

// CreateMessageBuilder is a helper function for creating message builder
// func CreateMessageBuilder(messageId string) *localizeBuilder {
//
//		return &localizeBuilder{
//			localize: &localize{
//				messageId: messageId,
//			},
//		}
//	}

// CreateMsg is a helper function for creating message with context
func CreateMsg(ctx *fiber.Ctx, messageId string, templateData ...map[string]string) string {

	loc := i18n.NewLocalizer(bundle, config.GetLanguage(ctx))
	msg := loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID: messageId,
	})

	if templateData != nil {
		msg = loc.MustLocalize(&i18n.LocalizeConfig{
			MessageID:    messageId,
			TemplateData: templateData[0],
		})
	}

	return msg
}
