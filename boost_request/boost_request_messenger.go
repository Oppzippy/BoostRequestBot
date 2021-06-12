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
	m := message_generator.NewBackendSignupMessage(
		messenger.localizer("en"),
		message_generator.NewDiscountFormatter(
			messenger.localizer("en"),
			message_generator.NewDiscordRoleNameProvider(messenger.discord),
		),
		br,
	)

	message, err := messenger.send(&MessageDestination{
		DestinationID:   br.Channel.BackendChannelID,
		DestinationType: DestinationChannel,
	}, m)

	return message, err
}

func (messenger *BoostRequestMessenger) SendBoostRequestCreatedDM(br *repository.BoostRequest) (*discordgo.Message, error) {
	localizer := messenger.localizer("en")

	m := message_generator.NewBoostRequestCreatedDM(localizer, messenger.discord, br)

	message, err := messenger.send(&MessageDestination{
		DestinationID:     br.RequesterID,
		DestinationType:   DestinationUser,
		FallbackChannelID: br.Channel.FrontendChannelID,
	}, m)

	return message, err
}

func (messenger *BoostRequestMessenger) SendBackendAdvertiserChosenMessage(
	br *repository.BoostRequest,
) (*discordgo.Message, error) {
	m := message_generator.NewBackendAdvertiserChosenMessage(messenger.localizer("en"), messenger.discord, br)

	message, err := messenger.send(&MessageDestination{
		DestinationID:   br.Channel.BackendChannelID,
		DestinationType: DestinationChannel,
	}, m)

	return message, err
}

func (messenger *BoostRequestMessenger) SendAdvertiserChosenDMToRequester(
	br *repository.BoostRequest,
) (*discordgo.Message, error) {
	localizer := messenger.localizer("en")
	m := message_generator.NewAdvertiserChosenDMToRequester(
		localizer,
		messenger.discord,
		message_generator.NewDiscountFormatter(
			localizer,
			message_generator.NewDiscordRoleNameProvider(messenger.discord),
		),
		br,
	)
	message, err := messenger.send(&MessageDestination{
		DestinationID:     br.RequesterID,
		DestinationType:   DestinationUser,
		FallbackChannelID: br.Channel.FrontendChannelID,
	}, m)

	return message, err
}

func (messenger *BoostRequestMessenger) SendAdvertiserChosenDMToAdvertiser(
	br *repository.BoostRequest,
) (*discordgo.Message, error) {
	localizer := messenger.localizer("en")
	m := message_generator.NewAdvertiserChosenDMToAdvertiser(
		localizer,
		messenger.discord,
		message_generator.NewDiscountFormatter(
			localizer,
			message_generator.NewDiscordRoleNameProvider(messenger.discord),
		),
		br,
	)
	message, err := messenger.send(&MessageDestination{
		DestinationID:     br.AdvertiserID,
		DestinationType:   DestinationUser,
		FallbackChannelID: br.Channel.BackendChannelID,
	}, m)
	return message, err
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
	Message() (*discordgo.MessageSend, error)
}

func (messenger *BoostRequestMessenger) send(dest *MessageDestination, sendableMessage sendable) (*discordgo.Message, error) {
	channelID, err := dest.ResolveChannelID(messenger.discord)
	if err != nil {
		return nil, fmt.Errorf("resolving channel id: %v", err)
	}
	message, err := sendableMessage.Message()
	if err != nil {
		return nil, fmt.Errorf("generating message: %v", err)
	}
	if message.Embed != nil {
		message.Embed.Footer = &discordgo.MessageEmbedFooter{
			Text:    "Huokan Boosting Community",
			IconURL: "https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png",
		}
		message.Embed.Timestamp = time.Now().Format(time.RFC3339)
	}
	m, err := messenger.discord.ChannelMessageSendComplex(channelID, message)

	if err != nil && dest.DestinationType == DestinationUser {
		restErr, ok := err.(discordgo.RESTError)
		if ok && restErr.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			if dest.FallbackChannelID != "" {
				messenger.sendDMBlockedMessage(dest.FallbackChannelID, dest.DestinationID)
			}
		}
	}

	return m, err
}

func (messenger *BoostRequestMessenger) sendDMBlockedMessage(channelID, userID string) {
	m := message_generator.NewDMBlockedMessage(messenger.localizer("en"), userID)
	message, err := m.Message()
	if err != nil {
		log.Printf("error generating dm blocked message: %v", err)
		return
	}
	messenger.sendTemporaryMessage(channelID, message)
}

func (messenger *BoostRequestMessenger) sendTemporaryMessage(channelID string, content *discordgo.MessageSend) {
	message, err := messenger.discord.ChannelMessageSendComplex(channelID, content)
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
