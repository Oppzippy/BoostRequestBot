package messages

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type AdvertiserChosenDMToRequester struct {
	localizer    *i18n.Localizer
	userProvider userProvider
	boostRequest *repository.BoostRequest
	embedPartial *partials.BoostRequestEmbedTemplate
}

func NewAdvertiserChosenDMToRequester(
	localizer *i18n.Localizer, up userProvider, br *repository.BoostRequest,
) *AdvertiserChosenDMToRequester {
	return &AdvertiserChosenDMToRequester{
		localizer:    localizer,
		userProvider: up,
		boostRequest: br,
		embedPartial: partials.NewBoostRequestEmbedTemplate(localizer, br),
	}
}

func (m *AdvertiserChosenDMToRequester) Message() (*discordgo.MessageSend, error) {
	advertiser, err := m.userProvider.User(m.boostRequest.AdvertiserID)
	if err != nil {
		return nil, err
	}

	var content string
	if m.boostRequest.NameVisibility == repository.NameVisibilityHide {
		content = m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "ClaimerNameHidden",
				Other: "The claimer's name is hidden.",
			},
		})
	} else {
		content = m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "AdvertiserChosenDMToRequester",
				Other: "{{.AdvertiserMention}} {{.AdvertiserTag}} will reach out to you shortly.",
			},
			TemplateData: map[string]string{
				"AdvertiserMention": advertiser.Mention(),
				"AdvertiserTag":     advertiser.String(),
			},
		})
	}

	embed, err := m.embedPartial.Embed(partials.BoostRequestEmbedConfiguration{
		Description: content,
		Price:       true,
	})
	if err != nil {
		return nil, err
	}
	embed.Color = 0x00FF00
	if m.boostRequest.NameVisibility != repository.NameVisibilityHide {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: advertiser.AvatarURL(""),
		}
	}

	return &discordgo.MessageSend{
		Embed: embed,
	}, nil
}
