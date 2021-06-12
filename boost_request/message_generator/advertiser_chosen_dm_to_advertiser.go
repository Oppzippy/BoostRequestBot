package message_generator

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type AdvertiserChosenDMToAdvertiser struct {
	localizer         *i18n.Localizer
	boostRequest      *repository.BoostRequest
	userProvider      userProvider
	discountFormatter *DiscountFormatter
}

func NewAdvertiserChosenDMToAdvertiser(
	localizer *i18n.Localizer, up userProvider, df *DiscountFormatter, br *repository.BoostRequest,
) *AdvertiserChosenDMToAdvertiser {
	return &AdvertiserChosenDMToAdvertiser{
		localizer:         localizer,
		boostRequest:      br,
		userProvider:      up,
		discountFormatter: df,
	}
}

func (m *AdvertiserChosenDMToAdvertiser) Message() (*discordgo.MessageSend, error) {
	if m.boostRequest.EmbedFields == nil {
		return m.humanMessage()
	} else {
		return m.botMessage()
	}
}

func (m *AdvertiserChosenDMToAdvertiser) humanMessage() (*discordgo.MessageSend, error) {
	requester, err := m.userProvider.User(m.boostRequest.RequesterID)
	if err != nil {
		return nil, err
	}

	fields := formatBoostRequest(m.localizer, m.boostRequest)

	if len(m.boostRequest.RoleDiscounts) != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "RequesterEligibleForDiscounts",
					Other: "The requester is eligible for discounts",
				},
			}),
			Value: m.discountFormatter.FormatDiscounts(m.boostRequest.RoleDiscounts),
		})
	}

	return &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: 0xFF0000,
			Title: "You have been selected to handle a boost request.",
			Description: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "PleaseMessage",
					Other: "Please message {{.RequesterMention}} {{.RequesterTag}}.",
				},
				TemplateData: map[string]string{
					"RequesterMention": requester.Mention(),
					"RequesterTag":     requester.String(),
				},
			}),
			Fields: fields,
		},
	}, nil
}

func (m *AdvertiserChosenDMToAdvertiser) botMessage() (*discordgo.MessageSend, error) {
	return &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: 0xFF0000,
			Title: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "YouHandleBoostRequest",
					Other: "You have been selected to handle a boost request.",
				},
			}),
			Description: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "MessageUserBelow",
					Other: "Please message the user listed below.",
				},
			}),
			Fields: formatBoostRequest(m.localizer, m.boostRequest),
		},
	}, nil
}
