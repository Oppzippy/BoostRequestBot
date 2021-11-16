package messages

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BackendAdvertiserChosenMessage struct {
	localizer    *i18n.Localizer
	boostRequest *repository.BoostRequest
	userProvider userProvider
	embedPartial *partials.BoostRequestEmbedPartial
}

func NewBackendAdvertiserChosenMessage(
	localizer *i18n.Localizer, up userProvider, df *partials.DiscountFormatter, br *repository.BoostRequest,
) *BackendAdvertiserChosenMessage {
	return &BackendAdvertiserChosenMessage{
		localizer:    localizer,
		boostRequest: br,
		userProvider: up,
		embedPartial: partials.NewBoostRequestEmbedPartial(localizer, df, br),
	}
}

func (m *BackendAdvertiserChosenMessage) Message() (*discordgo.MessageSend, error) {
	advertiser, err := m.userProvider.User(m.boostRequest.AdvertiserID)
	if err != nil {
		return nil, err
	}

	description := m.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "AdvertiserWillHandleBoostRequest",
			Other: "{{.AdvertiserMention}} will handle the following boost request.",
		},
		TemplateData: map[string]string{
			"AdvertiserMention": fmt.Sprintf("<@%s>", m.boostRequest.AdvertiserID),
		},
	})

	embed, err := m.embedPartial.Embed(partials.BoostRequestEmbedConfiguration{
		Description:    description,
		Price:          true,
		AdvertiserCut:  true,
		Discount:       true,
		DiscountTotals: true,
	})
	if err != nil {
		return nil, err
	}
	embed.Title = m.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "AdvertiserSelected",
			Other: "An advertiser has been selected.",
		},
	})
	embed.Color = 0xFF0000
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: advertiser.AvatarURL(""),
	}

	return &discordgo.MessageSend{
		Embed: embed,
	}, nil
}
