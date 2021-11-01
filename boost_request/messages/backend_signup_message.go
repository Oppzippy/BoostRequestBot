package messages

import (
	"fmt"
	"strings"

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
	fields := make([]*discordgo.MessageEmbedField, 0, 3)
	if price := m.priceField(); price != nil {
		fields = append(fields, price)
	}
	if advertiserCut := m.advertiserCutField(); advertiserCut != nil {
		fields = append(fields, advertiserCut)
	}
	if rd := m.roleDiscountField(); rd != nil {
		fields = append(fields, rd)
	}

	if len(fields) == 0 {
		fields = nil
	}

	var preferredAdvertiserMentions string
	if len(br.PreferredAdvertiserIDs) > 0 {
		mentions := make([]string, len(br.PreferredAdvertiserIDs))
		for i, id := range br.PreferredAdvertiserIDs {
			mentions[i] = fmt.Sprintf("<@%s>", id)
		}
		preferredAdvertiserMentions = strings.Join(mentions, " ")
	}

	return &discordgo.MessageSend{
		Content: preferredAdvertiserMentions,
		Embed: &discordgo.MessageEmbed{
			Color: 0x0000FF,
			Title: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "NewBoostRequest",
					One:   "New Boost Request",
					Other: "New Boost Requests",
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
			Value: m.discountFormatter.FormatDiscounts(m.boostRequest.RoleDiscounts),
		}
	}
	return nil
}

func (m *BackendSignupMessage) priceField() *discordgo.MessageEmbedField {
	if m.boostRequest.Price != 0 {
		return &discordgo.MessageEmbedField{
			Name:  "Price",
			Value: formatCopper(m.localizer, m.boostRequest.Price),
		}
	}
	return nil
}

func (m *BackendSignupMessage) advertiserCutField() *discordgo.MessageEmbedField {
	if m.boostRequest.AdvertiserCut != 0 {
		return &discordgo.MessageEmbedField{
			Name:  "Advertiser Cut",
			Value: formatCopper(m.localizer, m.boostRequest.AdvertiserCut),
		}
	}
	return nil
}
