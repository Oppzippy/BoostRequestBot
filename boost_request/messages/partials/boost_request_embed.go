package partials

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/message_utils"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestEmbedPartial struct {
	boostRequest      *repository.BoostRequest
	localizer         *i18n.Localizer
	discountFormatter *DiscountFormatter
}

type BoostRequestEmbedConfiguration struct {
	PreferredAdvertisers bool
	Description          string
	Price                bool
	AdvertiserCut        bool
	Discount             bool
	DiscountTotals       bool
	ID                   bool
}

func NewBoostRequestEmbedPartial(
	localizer *i18n.Localizer, df *DiscountFormatter, br *repository.BoostRequest,
) *BoostRequestEmbedPartial {
	return &BoostRequestEmbedPartial{
		boostRequest:      br,
		localizer:         localizer,
		discountFormatter: df,
	}
}

func (m *BoostRequestEmbedPartial) Embed(config BoostRequestEmbedConfiguration) (*discordgo.MessageEmbed, error) {
	embed := &discordgo.MessageEmbed{
		Title: m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "BoostRequest",
				One:   "Boost Request",
				Other: "Boost Requests",
			},
			PluralCount: 1,
		}),
		Description: m.boostRequest.Message,
		Fields:      make([]*discordgo.MessageEmbedField, 0, 10),
		Color:       0x0000FF,
	}

	if preferredAdvertisers := m.preferredAdvertisersField(); config.PreferredAdvertisers && preferredAdvertisers != nil {
		embed.Fields = append(embed.Fields, preferredAdvertisers)
	}
	if config.Description != "" {
		embed.Description = config.Description
		if message := m.messageField(); message != nil {
			embed.Fields = append(embed.Fields, message)
		}
	}
	if price := m.priceField(); config.Price && price != nil {
		embed.Fields = append(embed.Fields, price)
	}
	if advertiserCut := m.advertiserCutField(); config.AdvertiserCut && advertiserCut != nil {
		embed.Fields = append(embed.Fields, advertiserCut)
	}
	if rd := m.roleDiscountFields(); config.Discount && rd != nil {
		embed.Fields = append(embed.Fields, rd)
	}
	if totals := m.discountTotalsFields(); config.DiscountTotals && totals != nil {
		embed.Fields = append(embed.Fields, totals...)
	}
	if len(m.boostRequest.EmbedFields) != 0 {
		embed.Fields = append(embed.Fields, repository.ToDiscordEmbedFields(m.boostRequest.EmbedFields)...)
	}
	if config.ID {
		embed.Footer = m.idFooter()
	}

	if len(embed.Fields) == 0 {
		embed.Fields = nil
	}

	return embed, nil
}

func (m *BoostRequestEmbedPartial) preferredAdvertisersField() *discordgo.MessageEmbedField {
	if len(m.boostRequest.PreferredAdvertiserIDs) > 0 {
		mentions := make([]string, 0, len(m.boostRequest.PreferredAdvertiserIDs))
		for id := range m.boostRequest.PreferredAdvertiserIDs {
			mentions = append(mentions, fmt.Sprintf("<@%s>", id))
		}
		return &discordgo.MessageEmbedField{
			Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "PreferredAdvertiser",
					One:   "Preferred Advertiser",
					Other: "Preferred Advertisers",
				},
				PluralCount: len(mentions),
			}),
			Value: strings.Join(mentions, " "),
		}
	}
	return nil
}

func (m *BoostRequestEmbedPartial) messageField() *discordgo.MessageEmbedField {
	if m.boostRequest.Message != "" {
		return &discordgo.MessageEmbedField{
			Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "Message",
					One:   "Message",
					Other: "Messages",
				},
				PluralCount: 1,
			}),
			Value: m.boostRequest.Message,
		}
	}
	return nil
}

func (m *BoostRequestEmbedPartial) priceField() *discordgo.MessageEmbedField {
	if m.boostRequest.Price != 0 {
		return &discordgo.MessageEmbedField{
			Name:   "Price",
			Value:  message_utils.FormatCopper(m.localizer, m.boostRequest.Price),
			Inline: true,
		}
	}
	return nil
}

func (m *BoostRequestEmbedPartial) advertiserCutField() *discordgo.MessageEmbedField {
	if m.boostRequest.AdvertiserCut != 0 {
		return &discordgo.MessageEmbedField{
			Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "BaseAdvertiserCut",
					One:   "Base Advertiser Cut",
					Other: "Base Advertiser Cuts",
				},
				PluralCount: 1,
			}),
			Value:  message_utils.FormatCopper(m.localizer, m.boostRequest.AdvertiserCut),
			Inline: true,
		}
	}
	return nil
}

func (m *BoostRequestEmbedPartial) roleDiscountFields() *discordgo.MessageEmbedField {
	if m.boostRequest.Price != 0 {
		return &discordgo.MessageEmbedField{
			Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "Discount",
					One:   "Discount",
					Other: "Discounts",
				},
				PluralCount: 1,
			}),
			Value:  message_utils.FormatCopper(m.localizer, m.boostRequest.Discount),
			Inline: true,
		}
	} else if len(m.boostRequest.RoleDiscounts) != 0 {
		return &discordgo.MessageEmbedField{
			Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "Discount",
					One:   "Discount",
					Other: "Discounts",
				},
				PluralCount: 10,
			}),
			Value: m.discountFormatter.FormatDiscounts(m.boostRequest.RoleDiscounts),
		}
	}
	return nil
}

func (m *BoostRequestEmbedPartial) discountTotalsFields() []*discordgo.MessageEmbedField {
	if m.boostRequest.Discount != 0 && m.boostRequest.Price != 0 {
		fields := []*discordgo.MessageEmbedField{
			{
				Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "DiscountedPrice",
						One:   "Discounted Price",
						Other: "Discounted Prices",
					},
					PluralCount: 1,
				}),
				Inline: true,
				Value:  message_utils.FormatCopper(m.localizer, m.boostRequest.Price-m.boostRequest.Discount),
			},
		}
		if m.boostRequest.AdvertiserCut != 0 {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "DiscountedBaseAdvertiserCut",
						One:   "Discounted Base Advertiser Cut",
						Other: "Discounted Base Advertiser Cuts",
					},
					PluralCount: 1,
				}),
				Value:  message_utils.FormatCopper(m.localizer, m.boostRequest.AdvertiserCut-m.boostRequest.Discount),
				Inline: true,
			})
		}
		return fields
	}
	return nil
}

func (m *BoostRequestEmbedPartial) idFooter() *discordgo.MessageEmbedFooter {
	return &discordgo.MessageEmbedFooter{
		Text: m.boostRequest.ExternalID.String(),
	}
}
