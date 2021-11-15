package messages

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type AdvertiserChosenDMToRequester struct {
	localizer         *i18n.Localizer
	userProvider      userProvider
	discountFormatter *DiscountFormatter
	boostRequest      *repository.BoostRequest
}

func NewAdvertiserChosenDMToRequester(
	localizer *i18n.Localizer, up userProvider, df *DiscountFormatter, br *repository.BoostRequest,
) *AdvertiserChosenDMToRequester {
	return &AdvertiserChosenDMToRequester{
		localizer:         localizer,
		userProvider:      up,
		discountFormatter: df,
		boostRequest:      br,
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

	var fields []*discordgo.MessageEmbedField
	if m.boostRequest.Discount != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "Discount",
					Other: formatCopper(m.localizer, m.boostRequest.Discount),
				},
			}),
		})
	} else if len(m.boostRequest.RoleDiscounts) != 0 {
		fields = make([]*discordgo.MessageEmbedField, 1)
		fields[0] = &discordgo.MessageEmbedField{
			Name: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "YouEligibleForDiscounts",
					Other: "You are eligible for discounts",
				},
			}),
			Value: m.discountFormatter.FormatDiscounts(m.boostRequest.RoleDiscounts),
		}
	}

	return &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: 0x00FF00,
			Title: m.localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:  "BoostRequest",
					One: "Boost Request",
				},
				PluralCount: 1,
			}),
			Description: content,
			Fields:      fields,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: advertiser.AvatarURL(""),
			},
		},
	}, nil
}
