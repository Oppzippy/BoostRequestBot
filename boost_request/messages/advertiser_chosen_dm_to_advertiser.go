package messages

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type AdvertiserChosenDMToAdvertiser struct {
	localizer    *i18n.Localizer
	boostRequest *repository.BoostRequest
	userProvider userProvider
	embedPartial *partials.BoostRequestEmbedTemplate
}

func NewAdvertiserChosenDMToAdvertiser(
	localizer *i18n.Localizer, up userProvider, br *repository.BoostRequest,
) *AdvertiserChosenDMToAdvertiser {
	return &AdvertiserChosenDMToAdvertiser{
		localizer:    localizer,
		boostRequest: br,
		userProvider: up,
		embedPartial: partials.NewBoostRequestEmbedTemplate(localizer, br),
	}
}

func (m *AdvertiserChosenDMToAdvertiser) Message() (*discordgo.MessageSend, error) {
	requester, err := m.userProvider.User(m.boostRequest.RequesterID)
	if err != nil {
		restError, ok := err.(*discordgo.RESTError)
		if !(ok && restError.Message != nil && restError.Message.Code == discordgo.ErrCodeUnknownUser) {
			return nil, err
		}
	}

	var description string
	if m.boostRequest.NameVisibility == repository.NameVisibilityHide {
		description = m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "RequesterNameHidden",
				Other: "The requester's name is hidden.",
			},
		})
	} else if requester != nil {
		description = m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "PleaseMessage",
				Other: "Please message {{.RequesterMention}} {{.RequesterTag}}.",
			},
			TemplateData: map[string]string{
				"RequesterMention": requester.Mention(),
				"RequesterTag":     requester.String(),
			},
		})
	} else {
		description = m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "MessageUserBelow",
				One:   "Please message the user listed below.",
				Other: "Please message the users listed below.",
			},
			PluralCount: 1,
		})
	}

	embed, err := m.embedPartial.Embed(partials.BoostRequestEmbedConfiguration{
		Description: description,
		Price:       true,
		ID:          true,
	})
	if err != nil {
		return nil, err
	}
	embed.Title = m.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "YouHandleBoostRequest",
			One:   "You have been selected to handle a boost request.",
			Other: "You have been selected to handle boost requests.",
		},
		PluralCount: 1,
	})
	if requester != nil && m.boostRequest.NameVisibility != repository.NameVisibilityHide {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: requester.AvatarURL(""),
		}
	}

	return &discordgo.MessageSend{
		Content:         requester.Mention(),
		AllowedMentions: &discordgo.MessageAllowedMentions{},
		Embed:           embed,
	}, nil
}
