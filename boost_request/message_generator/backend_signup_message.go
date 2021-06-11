package message_generator

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BackendSignupMessage struct {
	boostRequest      *repository.BoostRequest
	localizer         *i18n.Localizer
	discountFormatter *DiscountFormatter
}

func NewBackendSignupMessage(
	localizer *i18n.Localizer, df *DiscountFormatter, br *repository.BoostRequest,
) *BackendSignupMessage {
	return &BackendSignupMessage{
		boostRequest:      br,
		localizer:         localizer,
		discountFormatter: df,
	}
}

func (m *BackendSignupMessage) Message() (*discordgo.MessageSend, error) {
	br := m.boostRequest
	var fields []*discordgo.MessageEmbedField
	if rd := m.roleDiscountField(); rd != nil {
		fields = make([]*discordgo.MessageEmbedField, 1)
		fields[0] = rd
	}

	return &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: 0x0000FF,
			Title: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:  "NewBoostRequest",
					One: "New Boost Request",
				},
				PluralCount: 1,
			}),
			Description: br.Message,
			Fields:      fields,
		},
	}, nil
}

func (m *BackendSignupMessage) roleDiscountField() *discordgo.MessageEmbedField {
	if len(m.boostRequest.RoleDiscounts) != 0 {
		return &discordgo.MessageEmbedField{
			Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "RequesterEligibleForDiscounts",
					Other: "The requester is eligible for discounts",
				},
			}),
			Value: m.discountFormatter.FormatDiscounts(m.boostRequest),
		}
	}
	return nil
}
