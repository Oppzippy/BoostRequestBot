package messages

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/message_utils"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type AdvertiserChosenDMToAdvertiser struct {
	localizer         *i18n.Localizer
	boostRequest      *repository.BoostRequest
	userProvider      userProvider
	discountFormatter *partials.DiscountFormatter
	embedPartial      *partials.BoostRequestEmbedPartial
}

func NewAdvertiserChosenDMToAdvertiser(
	localizer *i18n.Localizer, up userProvider, df *partials.DiscountFormatter, br *repository.BoostRequest,
) *AdvertiserChosenDMToAdvertiser {
	return &AdvertiserChosenDMToAdvertiser{
		localizer:         localizer,
		boostRequest:      br,
		userProvider:      up,
		discountFormatter: df,
		embedPartial:      partials.NewBoostRequestEmbedPartial(localizer, df, br),
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
	if requester != nil {
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

	if m.boostRequest.Price != 0 {
		description += "\n\n"
		description += m.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "AdvertiserConfirmReminder",
				Other: "Please collect {{.Gold}} from the buyer. Then, send a direct message to <@720340847928934531> (Huokan Bot) with the screenshot attached to the following command:\n`!brconfirm {{.BoostRequestID}}`\nOtherwise, the buyer will not receive any buyer points.",
			},
			TemplateData: map[string]string{
				"Gold":           message_utils.FormatCopper(m.localizer, m.boostRequest.Price-m.boostRequest.Discount),
				"BoostRequestID": m.boostRequest.ExternalID.String(),
			},
		})
	}

	embed, err := m.embedPartial.Embed(partials.BoostRequestEmbedConfiguration{
		Description:    description,
		Price:          true,
		AdvertiserCut:  true,
		Discount:       true,
		DiscountTotals: true,
		ID:             true,
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
	if requester != nil {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: requester.AvatarURL(""),
		}
	}

	return &discordgo.MessageSend{
		Embed: embed,
	}, nil
}
