package partials

import (
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
	Description    string
	Price          bool
	AdvertiserCut  bool
	Discount       bool
	DiscountTotals bool
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
		Fields:      make([]*discordgo.MessageEmbedField, 0, 6),
		Color:       0x0000FF,
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
	if totals := m.DiscountTotalsFields(); config.DiscountTotals && totals != nil {
		embed.Fields = append(embed.Fields, totals...)
	}
	if len(m.boostRequest.EmbedFields) != 0 {
		embed.Fields = append(embed.Fields, repository.ToDiscordEmbedFields(m.boostRequest.EmbedFields)...)
	}

	if len(embed.Fields) == 0 {
		embed.Fields = nil
	}

	return embed, nil
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
			Name:   "Advertiser Cut",
			Value:  message_utils.FormatCopper(m.localizer, m.boostRequest.AdvertiserCut),
			Inline: true,
		}
	}
	return nil
}

func (m *BoostRequestEmbedPartial) roleDiscountFields() *discordgo.MessageEmbedField {
	if m.boostRequest.Discount != 0 && m.boostRequest.Price != 0 {
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
					ID:    "RequesterEligibleForDiscounts",
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

func (m *BoostRequestEmbedPartial) DiscountTotalsFields() []*discordgo.MessageEmbedField {
	if m.boostRequest.Discount != 0 && m.boostRequest.Price != 0 {
		return []*discordgo.MessageEmbedField{
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
			{
				Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "DiscountedAdvertiserCut",
						One:   "Discounted Advertiser Cut",
						Other: "Discounted Advertiser Cuts",
					},
					PluralCount: 1,
				}),
				Value:  message_utils.FormatCopper(m.localizer, m.boostRequest.AdvertiserCut-m.boostRequest.Discount),
				Inline: true,
			},
		}
	}
	return nil
}
