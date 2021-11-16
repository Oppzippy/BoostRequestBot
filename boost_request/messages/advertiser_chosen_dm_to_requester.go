package messages

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type AdvertiserChosenDMToRequester struct {
	localizer         *i18n.Localizer
	userProvider      userProvider
	discountFormatter *partials.DiscountFormatter
	boostRequest      *repository.BoostRequest
	embedPartial      *partials.BoostRequestEmbedPartial
}

func NewAdvertiserChosenDMToRequester(
	localizer *i18n.Localizer, up userProvider, df *partials.DiscountFormatter, br *repository.BoostRequest,
) *AdvertiserChosenDMToRequester {
	return &AdvertiserChosenDMToRequester{
		localizer:         localizer,
		userProvider:      up,
		discountFormatter: df,
		boostRequest:      br,
		embedPartial:      partials.NewBoostRequestEmbedPartial(localizer, df, br),
	}
}

func (m *AdvertiserChosenDMToRequester) Message() (*discordgo.MessageSend, error) {
	advertiser, err := m.userProvider.User(m.boostRequest.AdvertiserID)
	if err != nil {
		return nil, err
	}

	content := m.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "AdvertiserChosenDMToRequester",
			Other: "{{.AdvertiserMention}} {{.AdvertiserTag}} will reach out to you shortly. Anyone else that messages you regarding this boost request is not from Huokan Boosting Community and may attempt to scam you.",
		},
		TemplateData: map[string]string{
			"AdvertiserMention": advertiser.Mention(),
			"AdvertiserTag":     advertiser.String(),
		},
	})

	embed, err := m.embedPartial.Embed(partials.BoostRequestEmbedConfiguration{
		Description: content,
		Price:       true,
		Discount:    true,
	})
	if err != nil {
		return nil, err
	}
	embed.Color = 0x00FF00
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: advertiser.AvatarURL(""),
	}

	return &discordgo.MessageSend{
		Embed: embed,
	}, nil
}
