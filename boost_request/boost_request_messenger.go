package boost_request

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/message_generator"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/oppzippy/BoostRequestBot/roll"
	"github.com/shopspring/decimal"
)

type BoostRequestMessenger struct {
	Destroyed bool
	discord   *discordgo.Session
	bundle    *i18n.Bundle
	waitGroup *sync.WaitGroup
	quit      chan struct{}
}

var footer = &discordgo.MessageEmbedFooter{
	Text:    "Huokan Boosting Community",
	IconURL: "https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png",
}

func NewBoostRequestMessenger(discord *discordgo.Session, bundle *i18n.Bundle) *BoostRequestMessenger {
	brm := BoostRequestMessenger{
		Destroyed: false,
		discord:   discord,
		bundle:    bundle,
		waitGroup: new(sync.WaitGroup),
		quit:      make(chan struct{}),
	}
	return &brm
}

func (messenger *BoostRequestMessenger) SendBackendSignupMessage(br *repository.BoostRequest) (*discordgo.Message, error) {
	gen := message_generator.NewGenerator(messenger.localizer("en"), messenger.discord)
	bsm := gen.BackendSignupMessage(br)

	message, err := messenger.send(bsm)

	return message, err
}

func (messenger *BoostRequestMessenger) SendBoostRequestCreatedDM(br *repository.BoostRequest) (*discordgo.Message, error) {
	requester, err := messenger.discord.User(br.RequesterID)
	if err != nil {
		return nil, err
	}
	dmChannel, _ := messenger.discord.UserChannelCreate(requester.ID)
	message, err := messenger.discord.ChannelMessageSendComplex(dmChannel.ID, &discordgo.MessageSend{
		Content: "Please wait while we find an advertiser to complete your request.",
		Embed: &discordgo.MessageEmbed{
			Title: "Huokan Boosting Community Boost Request",
			Author: &discordgo.MessageEmbedAuthor{
				Name: requester.String(),
			},
			Description: br.Message,
			Footer:      footer,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: requester.AvatarURL(""),
			},
		},
	})
	if err != nil {
		restErr, ok := err.(*discordgo.RESTError)
		if ok && restErr.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			messenger.sendTemporaryMessage(
				br.Channel.FrontendChannelID,
				requester.Mention()+", I can't DM you. Please allow DMs from server members by right clicking the server "+
					"and enabling \"Allow direct messages from server members.\" in Privacy Settings.",
			)
		}
		return nil, err
	}
	return message, err
}

func (messenger *BoostRequestMessenger) SendBackendAdvertiserChosenMessage(
	br *repository.BoostRequest,
) (*discordgo.Message, error) {
	advertiser, err := messenger.discord.User(br.AdvertiserID)

	if err != nil {
		return nil, err
	}

	message, err := messenger.discord.ChannelMessageSendEmbed(br.Channel.BackendChannelID, &discordgo.MessageEmbed{
		Color:       0xFF0000,
		Title:       "An advertiser has been selected.",
		Description: advertiser.Mention() + " will handle the following boost request.",
		Fields:      messenger.formatBoostRequest(br),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: advertiser.AvatarURL(""),
		},
	})

	return message, err
}

func (messenger *BoostRequestMessenger) SendAdvertiserChosenDMToRequester(
	br *repository.BoostRequest,
) (*discordgo.Message, error) {
	requester, err := messenger.discord.User(br.RequesterID)
	if err != nil {
		return nil, err
	}
	advertiser, err := messenger.discord.User(br.AdvertiserID)
	if err != nil {
		return nil, err
	}
	dmChannel, err := messenger.discord.UserChannelCreate(requester.ID)
	if err != nil {
		restErr, ok := err.(discordgo.RESTError)
		if ok && restErr.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			messenger.sendTemporaryMessage(
				br.Channel.FrontendChannelID,
				requester.Mention()+", I can't DM you. Please allow DMs from server members by right clicking the "+
					"server and enabling \"Allow direct messages from server members.\" in Privacy Settings and post "+
					"your message again.",
			)
		}
		return nil, err
	}

	sb := strings.Builder{}
	sb.WriteString(advertiser.Mention())
	sb.WriteString(" ")
	sb.WriteString(advertiser.String())
	sb.WriteString(" will reach out to you shortly.")
	sb.WriteString(" Anyone else that messages you regarding this boost request is not from Huokan and may attempt to scam you.")
	var fields []*discordgo.MessageEmbedField

	if br.RoleDiscounts != nil && len(br.RoleDiscounts) != 0 {
		fields = make([]*discordgo.MessageEmbedField, 1)
		fields[0] = &discordgo.MessageEmbedField{
			Name:  "You are eligible for discounts",
			Value: messenger.formatDiscounts(br),
		}
	}

	message, err := messenger.discord.ChannelMessageSendEmbed(dmChannel.ID, &discordgo.MessageEmbed{
		Color:       0x00FF00,
		Title:       "Huokan Boosting Community Boost Request",
		Description: sb.String(),
		Fields:      fields,
		Footer:      footer,
		Timestamp:   time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: advertiser.AvatarURL(""),
		},
	})
	return message, err
}

func (messenger *BoostRequestMessenger) SendAdvertiserChosenDMToAdvertiser(
	br *repository.BoostRequest,
) (*discordgo.Message, error) {
	if br.EmbedFields != nil {
		m, err := messenger.sendAdvertiserChosenDMToAdvertiserWithBotRequester(br)
		return m, err
	} else {
		m, err := messenger.sendAdvertiserChosenDMToAdvertiserWithHumanRequester(br)
		return m, err
	}
}

func (messenger *BoostRequestMessenger) SendRoll(
	channelID string, br *repository.BoostRequest, rollResults *roll.WeightedRollResults,
) (*discordgo.Message, error) {
	if rollResults == nil {
		return nil, errors.New("rollResults must not be nil")
	}

	sb := strings.Builder{}
	var weightAccumulator float64
	for iter := rollResults.Iterator(); iter.HasNext(); {
		advertiserID, weight, isChosenItem := iter.Next()
		weightAccumulator += weight

		sb.WriteString(fmt.Sprintf(
			"<@%s>: %s to %s",
			advertiserID,
			messenger.formatFloat(weightAccumulator-weight),
			messenger.formatFloat(weightAccumulator),
		))
		if isChosenItem {
			sb.WriteString(fmt.Sprintf(
				"   **<-- %s**",
				messenger.formatFloat(rollResults.Roll()),
			))
		}
		sb.WriteString("\n")
	}

	message, err := messenger.discord.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: br.Message,
		Embed: &discordgo.MessageEmbed{
			Title:       "Roll Results",
			Description: sb.String(),
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer:      footer,
		},
		AllowedMentions: &discordgo.MessageAllowedMentions{},
	})

	return message, err
}

func (messenger *BoostRequestMessenger) sendAdvertiserChosenDMToAdvertiserWithHumanRequester(
	br *repository.BoostRequest,
) (*discordgo.Message, error) {
	requester, err := messenger.discord.User(br.RequesterID)
	if err != nil {
		return nil, err
	}
	advertiser, err := messenger.discord.User(br.AdvertiserID)
	if err != nil {
		return nil, err
	}
	dmChannel, err := messenger.discord.UserChannelCreate(advertiser.ID)
	if err != nil {
		restErr, ok := err.(discordgo.RESTError)
		if ok && restErr.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			messenger.sendTemporaryMessage(
				br.Channel.BackendChannelID,
				advertiser.Mention()+", I can't DM you. Please allow DMs from server members by right clicking the server and enabling \"Allow direct messages from server members.\" in Privacy Settings.",
			)
			_, err := messenger.discord.ChannelMessageSend(br.Channel.BackendChannelID, "Please DM "+requester.Mention()+" ("+requester.String()+").")
			if err != nil {
				log.Printf("Failed to send backup message after failed DM: %v", err)
			}
		}
		return nil, err
	}

	sb := strings.Builder{}
	sb.WriteString("Please message ")
	sb.WriteString(requester.Mention())
	sb.WriteString(" ")
	sb.WriteString(requester.String())
	sb.WriteString(".")

	fields := messenger.formatBoostRequest(br)

	if br.RoleDiscounts != nil && len(br.RoleDiscounts) != 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "The requester is eligible for discounts",
			Value: messenger.formatDiscounts(br),
		})
	}

	message, err := messenger.discord.ChannelMessageSendEmbed(dmChannel.ID, &discordgo.MessageEmbed{
		Color:       0xFF0000,
		Title:       "You have been selected to handle a boost request.",
		Description: sb.String(),
		Fields:      fields,
		Footer:      footer,
		Timestamp:   time.Now().Format(time.RFC3339),
	})

	return message, err
}

func (messenger *BoostRequestMessenger) sendAdvertiserChosenDMToAdvertiserWithBotRequester(
	br *repository.BoostRequest,
) (*discordgo.Message, error) {
	advertiser, err := messenger.discord.User(br.AdvertiserID)
	if err != nil {
		return nil, err
	}
	dmChannel, err := messenger.discord.UserChannelCreate(advertiser.ID)
	if err != nil {
		restErr, ok := err.(discordgo.RESTError)
		if ok && restErr.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			messenger.sendTemporaryMessage(br.Channel.BackendChannelID, advertiser.Mention()+", I can't DM you. Please allow DMs from server members by right clicking the server and enabling \"Allow direct messages from server members.\" in Privacy Settings.")
		}
		return nil, err
	}

	message, err := messenger.discord.ChannelMessageSendEmbed(dmChannel.ID, &discordgo.MessageEmbed{
		Color:       0xFF0000,
		Title:       "You have been selected to handle a boost request.",
		Description: "Please message the user listed below.",
		Fields:      messenger.formatBoostRequest(br),
		Footer:      footer,
		Timestamp:   time.Now().Format(time.RFC3339),
	})

	return message, err
}

// Logs the creation of a boost request to a channel only moderators can view
func (messenger *BoostRequestMessenger) SendLogChannelMessage(
	br *repository.BoostRequest, channelID string,
) (*discordgo.Message, error) {
	if br.EmbedFields != nil {
		// TODO return an error
		return nil, nil
	}
	user, err := messenger.discord.User(br.RequesterID)
	if err != nil {
		return nil, err
	}

	embed := &discordgo.MessageEmbed{
		Color:       0x0000FF,
		Title:       "New Boost Request",
		Description: br.Message,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Requested By",
				Value: user.Mention() + " " + user.String(),
			},
		},
		Footer:    footer,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	message, err := messenger.discord.ChannelMessageSendEmbed(channelID, embed)
	return message, err
}

func (messenger *BoostRequestMessenger) SendCreditsUpdateDM(userID string, credits int) (*discordgo.Message, error) {
	dmChannel, err := messenger.discord.UserChannelCreate(userID)
	if err != nil {
		return nil, err
	}
	var plural string
	if credits != 1 {
		plural = "s"
	}
	message, err := messenger.discord.ChannelMessageSend(
		dmChannel.ID,
		fmt.Sprintf("You now have %d boost request steal credit%s.", credits, plural),
	)
	return message, err
}

func (messenger *BoostRequestMessenger) Destroy() {
	if !messenger.Destroyed {
		messenger.Destroyed = true
		close(messenger.quit)
		messenger.waitGroup.Wait()
	}
}

func (messenger *BoostRequestMessenger) sendTemporaryMessage(channelID string, content string) {
	message, err := messenger.discord.ChannelMessageSend(channelID, content)
	if err == nil {
		messenger.waitGroup.Add(1)
		go func() {
			select {
			case <-time.After(30 * time.Second):
			case <-messenger.quit:
			}
			err := messenger.discord.ChannelMessageDelete(message.ChannelID, message.ID)
			if err != nil {
				log.Printf("Error deleting temporary message: %v", err)
			}
			messenger.waitGroup.Done()
		}()
	} else {
		log.Printf("Error sending temporary message: %v", err)
	}
}

func (messenger *BoostRequestMessenger) formatBoostRequest(br *repository.BoostRequest) []*discordgo.MessageEmbedField {
	var fields []*discordgo.MessageEmbedField
	if br.EmbedFields != nil {
		fields = repository.ToDiscordEmbedFields(br.EmbedFields)
	} else {
		fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Boost Request",
				Value: br.Message,
			},
		}
	}
	return fields
}

func (messenger *BoostRequestMessenger) formatDiscounts(br *repository.BoostRequest) string {
	sb := strings.Builder{}
	if br.RoleDiscounts != nil && len(br.RoleDiscounts) != 0 {
		for _, roleDiscount := range br.RoleDiscounts {
			sb.WriteString(roleDiscount.Discount.Mul(decimal.NewFromInt(100)).String())
			sb.WriteString("% discount on ")
			sb.WriteString(roleDiscount.BoostType)

			roleName := messenger.getRoleName(roleDiscount.GuildID, roleDiscount.RoleID)
			if roleName != "" {
				sb.WriteString(" (")
				sb.WriteString(roleName)
				sb.WriteString(")")
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func (messenger *BoostRequestMessenger) getRoleName(guildID, roleID string) string {
	guild, err := messenger.discord.State.Guild(guildID)

	if err == nil {
		roles := guild.Roles
		for _, role := range roles {
			if role.ID == roleID {
				return role.Name
			}
		}
	}
	return ""
}

func (messenger *BoostRequestMessenger) formatFloat(f float64) string {
	return strings.TrimRight(
		strings.TrimRight(fmt.Sprintf("%.2f", f), "0"),
		".",
	)
}

func (messenger *BoostRequestMessenger) localizer(langs ...string) *i18n.Localizer {
	return i18n.NewLocalizer(messenger.bundle, langs...)
}

type sendable interface {
	ChannelID() (string, error)
	Message() (*discordgo.MessageSend, error)
}

func (messenger *BoostRequestMessenger) send(sendableMessage sendable) (*discordgo.Message, error) {
	channelID, err := sendableMessage.ChannelID()
	if err != nil {
		// TODO handle
		return nil, err
	}
	message, err := sendableMessage.Message()
	if err != nil {
		// TODO handle
		return nil, err
	}
	if message.Embed != nil {
		message.Embed.Footer = &discordgo.MessageEmbedFooter{
			Text:    "Huokan Boosting Community",
			IconURL: "https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png",
		}
		message.Embed.Timestamp = time.Now().Format(time.RFC3339)
	}
	m, err := messenger.discord.ChannelMessageSendComplex(channelID, message)
	return m, err
}
