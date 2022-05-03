package messages

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BackendAdvertiserChosenMessage struct {
	localizer    *i18n.Localizer
	boostRequest *repository.BoostRequest
	embedPartial *partials.BoostRequestEmbedTemplate
}

func NewBackendAdvertiserChosenMessage(
	localizer *i18n.Localizer, up userProvider, br *repository.BoostRequest,
) *BackendAdvertiserChosenMessage {
	return &BackendAdvertiserChosenMessage{
		localizer:    localizer,
		boostRequest: br,
		embedPartial: partials.NewBoostRequestEmbedTemplate(localizer, br),
	}
}

func (m *BackendAdvertiserChosenMessage) Message() (*discordgo.MessageSend, error) {
	embed, err := m.embedPartial.Embed(partials.BoostRequestEmbedConfiguration{
		Price: true,
		ID:    true,
	})
	if err != nil {
		return nil, err
	}
	embed.Title = m.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "BoostRequestClaimed",
			Other: "This boost request has been claimed.",
		},
	})
	embed.Color = 0xFF0000

	return &discordgo.MessageSend{
		Embed: embed,
	}, nil
}
