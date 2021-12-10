package messenger

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages"
	"github.com/oppzippy/BoostRequestBot/boost_request/messages/partials"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/oppzippy/BoostRequestBot/channels"
	"github.com/oppzippy/BoostRequestBot/roll"
)

type BoostRequestMessenger struct {
	Destroyed                bool
	discord                  *discordgo.Session
	bundle                   *i18n.Bundle
	messageBroker            *messageBroker
	rnp                      *messages.DiscordRoleNameProvider
	delayedMessageRepository repository.DelayedMessageRepository
}

func NewBoostRequestMessenger(
	discord *discordgo.Session,
	bundle *i18n.Bundle,
	delayedMessageRepository repository.DelayedMessageRepository,
) *BoostRequestMessenger {
	brm := BoostRequestMessenger{
		Destroyed:                false,
		discord:                  discord,
		bundle:                   bundle,
		messageBroker:            newMessageBroker(discord),
		rnp:                      messages.NewDiscordRoleNameProvider(discord),
		delayedMessageRepository: delayedMessageRepository,
	}
	brm.loadDelayedMessages()

	return &brm
}

func (messenger *BoostRequestMessenger) loadDelayedMessages() error {
	delayedMessages, err := messenger.delayedMessageRepository.GetDelayedMessages()
	if err != nil {
		return err
	}
	for _, delayedMessage := range delayedMessages {
		_, errChannel := messenger.addExistingDelayedMessage(delayedMessage)
		go func() {
			err, ok := <-errChannel
			if ok && err != nil {
				log.Printf("Error sending delayed message: %v", err)
			}
		}()
	}
	return nil
}

func (messenger *BoostRequestMessenger) Destroy() {
	messenger.messageBroker.Destroy()
}

func (messenger *BoostRequestMessenger) SendBackendSignupMessage(
	br *repository.BoostRequest,
	channelID string,
	buttonConfiguration BackendSignupMessageButtonConfiguration,
) (*discordgo.Message, error) {
	localizer := messenger.localizer("en")
	m := messages.NewBackendSignupMessage(
		localizer,
		partials.NewDiscountFormatter(localizer, messenger.rnp),
		br,
		buttonConfiguration,
	)

	message, err := messenger.send(&MessageDestination{
		DestinationID:   channelID,
		DestinationType: DestinationChannel,
	}, m)

	return message, err
}

func (messenger *BoostRequestMessenger) SendBoostRequestCreatedDM(br *repository.BoostRequest) (*discordgo.Message, error) {
	localizer := messenger.localizer("en")

	m := messages.NewBoostRequestCreatedDM(localizer,
		messenger.discord,
		partials.NewDiscountFormatter(localizer, messenger.rnp),
		br,
	)

	var fallbackChannelID string
	if br.Channel != nil {
		fallbackChannelID = br.Channel.FrontendChannelID
	}

	message, err := messenger.send(&MessageDestination{
		DestinationID:     br.RequesterID,
		DestinationType:   DestinationUser,
		FallbackChannelID: fallbackChannelID,
	}, m)

	return message, err
}

func (messenger *BoostRequestMessenger) SendBackendAdvertiserChosenMessage(
	br *repository.BoostRequest,
) ([]*discordgo.Message, error) {
	localizer := messenger.localizer("en")
	m := messages.NewBackendAdvertiserChosenMessage(
		localizer,
		messenger.discord,
		partials.NewDiscountFormatter(localizer, messenger.rnp),
		br,
	)

	if br.Channel != nil && br.Channel.UsesBuyerMessage {
		message, err := messenger.send(&MessageDestination{
			DestinationID:   br.Channel.BackendChannelID,
			DestinationType: DestinationChannel,
		}, m)

		return []*discordgo.Message{message}, err
	}
	content, err := m.Message()
	if err != nil {
		return nil, err
	}

	messages := make([]*discordgo.Message, 0, len(br.BackendMessages))
	errorMessages := make([]string, 0)
	for _, backendMessage := range br.BackendMessages {
		message, err := messenger.discord.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:         backendMessage.MessageID,
			Channel:    backendMessage.ChannelID,
			Embed:      content.Embed,
			Components: []discordgo.MessageComponent{},
		})
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
		messages = append(messages, message)
	}
	if len(messages) == 0 {
		return nil, errors.New(strings.Join(errorMessages, "\n"))
	}
	return messages, nil
}

func (messenger *BoostRequestMessenger) SendAdvertiserChosenDMToRequester(
	br *repository.BoostRequest,
) (*discordgo.Message, error) {
	localizer := messenger.localizer("en")
	m := messages.NewAdvertiserChosenDMToRequester(
		localizer,
		messenger.discord,
		partials.NewDiscountFormatter(localizer, messenger.rnp),
		br,
	)

	var fallbackChannelID string
	if br.Channel != nil {
		fallbackChannelID = br.Channel.FrontendChannelID
	}
	message, err := messenger.send(&MessageDestination{
		DestinationID:     br.RequesterID,
		DestinationType:   DestinationUser,
		FallbackChannelID: fallbackChannelID,
	}, m)

	return message, err
}

func (messenger *BoostRequestMessenger) SendAdvertiserChosenDMToAdvertiser(
	br *repository.BoostRequest,
) (*discordgo.Message, error) {
	localizer := messenger.localizer("en")
	m := messages.NewAdvertiserChosenDMToAdvertiser(
		localizer,
		messenger.discord,
		partials.NewDiscountFormatter(localizer, messenger.rnp),
		br,
	)

	var fallbackChannelID string
	if br.Channel != nil {
		fallbackChannelID = br.Channel.BackendChannelID
	}
	message, err := messenger.send(&MessageDestination{
		DestinationID:     br.AdvertiserID,
		DestinationType:   DestinationUser,
		FallbackChannelID: fallbackChannelID,
	}, m)
	return message, err
}

func (messenger *BoostRequestMessenger) SendRoll(
	channelID string, br *repository.BoostRequest, rollResults *roll.WeightedRollResults,
) (*discordgo.Message, error) {
	m := messages.NewBoostRequestRollMessage(messenger.localizer("en"), br, rollResults)
	message, err := messenger.send(&MessageDestination{
		DestinationID:   channelID,
		DestinationType: DestinationChannel,
	}, m)

	return message, err
}

// Logs the creation of a boost request to a channel only moderators can view
func (messenger *BoostRequestMessenger) SendLogChannelMessage(
	br *repository.BoostRequest, channelID string,
) (*discordgo.Message, error) {
	m := messages.NewLogChannelMessage(messenger.localizer("en"), messenger.discord, br)
	message, err := messenger.send(&MessageDestination{
		DestinationID:   channelID,
		DestinationType: DestinationChannel,
	}, m)
	return message, err
}

func (messenger *BoostRequestMessenger) SendCreditsUpdateDM(userID string, credits int) (*discordgo.Message, error) {
	m := messages.NewCreditsUpdatedDM(messenger.localizer("en"), credits)
	message, err := messenger.send(&MessageDestination{
		DestinationID:   userID,
		DestinationType: DestinationUser,
	}, m)
	return message, err
}

func (messenger *BoostRequestMessenger) SendAutoSignUpMessages(
	userID string,
	expiresAt time.Time,
) ([]*repository.DelayedMessage, <-chan error) {
	expiringSoonDelayedMessage, expiringSoonErrChannel := messenger.sendAutoSignUpExpiringSoonMessage(userID, expiresAt, 5*time.Minute)
	expiredDelayedMessage, expiredErrChannel := messenger.sendAutoSignUpExpiredMessage(userID, expiresAt)

	errChannel := channels.MergeErrorChannels(expiredErrChannel, expiringSoonErrChannel)

	delayedMessages := make([]*repository.DelayedMessage, 0, 2)
	if expiringSoonDelayedMessage != nil {
		delayedMessages = append(delayedMessages, expiringSoonDelayedMessage)
	}
	if expiredDelayedMessage != nil {
		delayedMessages = append(delayedMessages, expiredDelayedMessage)
	}
	return delayedMessages, errChannel
}

func (messenger *BoostRequestMessenger) sendAutoSignUpExpiringSoonMessage(
	userID string,
	expiresAt time.Time,
	warningTime time.Duration,
) (*repository.DelayedMessage, <-chan error) {
	m := messages.NewAutoSignUpExpiringSoonMessage(messenger.localizer("en"), warningTime)
	delay := time.Until(expiresAt.Add(warningTime * -1))
	if delay > 0 {
		messenger.sendDelayed(&MessageDestination{
			DestinationID:   userID,
			DestinationType: DestinationUser,
		}, m, delay)
	}
	errChannel := make(chan error)
	close(errChannel)
	return nil, errChannel
}

func (messenger *BoostRequestMessenger) sendAutoSignUpExpiredMessage(
	userID string,
	expiresAt time.Time,
) (*repository.DelayedMessage, <-chan error) {
	m := messages.NewAutoSignUpExpiredMessage(messenger.localizer("en"))
	delay := time.Until(expiresAt)
	delayedMessage, _, errChannel := messenger.sendDelayed(&MessageDestination{
		DestinationID:   userID,
		DestinationType: DestinationUser,
	}, m, delay)
	return delayedMessage, errChannel
}

func (messenger *BoostRequestMessenger) send(dest *MessageDestination, mg messageGenerator) (*discordgo.Message, error) {
	m, err := messenger.messageBroker.Send(dest, mg)
	if dest.DestinationType == DestinationUser && err == errDMBlocked {
		_, dmBlockedErr := messenger.sendDMBlockedMessage(dest.FallbackChannelID, dest.DestinationID)
		if dmBlockedErr != nil {
			return nil, fmt.Errorf("dm was blocked: %v, error sending dm blocked message: %v", err, dmBlockedErr)
		}
		return nil, fmt.Errorf("dm was blocked but the user was informed of the issue: %v", err)
	}
	return m, err
}

func (messenger *BoostRequestMessenger) sendDelayed(
	dest *MessageDestination,
	mg messageGenerator,
	delay time.Duration,
) (*repository.DelayedMessage, <-chan *discordgo.Message, <-chan error) {
	delayedMessageDTO, err := messenger.storeDelayedMessageInRepository(dest, mg, delay)
	if err != nil {
		messageChannel, errChannel := make(chan *discordgo.Message, 1), make(chan error, 1)
		errChannel <- err
		close(messageChannel)
		close(errChannel)
		return nil, messageChannel, errChannel
	}

	messageChannel, errChannel := messenger.addExistingDelayedMessage(delayedMessageDTO)
	return delayedMessageDTO, messageChannel, errChannel
}

func (messenger *BoostRequestMessenger) storeDelayedMessageInRepository(
	dest *MessageDestination,
	mg messageGenerator,
	delay time.Duration,
) (*repository.DelayedMessage, error) {
	destType := repository.DestinationTypeChannel
	if dest.DestinationType == DestinationUser {
		destType = repository.DestinationTypeUser
	}
	message, err := mg.Message()
	if err != nil {
		return nil, err
	}

	delayedMessageDTO := repository.DelayedMessage{
		DestinationID:     dest.DestinationID,
		DestinationType:   destType,
		FallbackChannelID: dest.FallbackChannelID,
		Message:           message,
		SendAt:            time.Now().Add(delay),
	}
	err = messenger.delayedMessageRepository.InsertDelayedMessage(&delayedMessageDTO)
	if err != nil {
		return nil, err
	}

	return &delayedMessageDTO, nil
}

func (messenger *BoostRequestMessenger) addExistingDelayedMessage(
	delayedMessageDTO *repository.DelayedMessage,
) (<-chan *discordgo.Message, <-chan error) {
	messageChannel, errChannel := make(chan *discordgo.Message, 1), make(chan error, 1)

	go func() {
		defer close(messageChannel)
		defer close(errChannel)

		destType := DestinationChannel
		if delayedMessageDTO.DestinationType == repository.DestinationTypeUser {
			destType = DestinationUser
		}

		sendDelayedMessage, sendDelayedErr := messenger.messageBroker.SendDelayed(
			&MessageDestination{
				DestinationID:     delayedMessageDTO.DestinationID,
				DestinationType:   destType,
				FallbackChannelID: delayedMessageDTO.FallbackChannelID,
			},
			messages.NewStaticMessage(delayedMessageDTO.Message),
			time.Until(delayedMessageDTO.SendAt),
		)
		select {
		case m := <-sendDelayedMessage:
			messageChannel <- m
		case err := <-sendDelayedErr:
			errChannel <- err
		}
		messenger.delayedMessageRepository.FlagDelayedMessageAsSent(delayedMessageDTO)
	}()

	return messageChannel, errChannel
}

func (messenger *BoostRequestMessenger) sendDMBlockedMessage(channelID, userID string) (*discordgo.Message, error) {
	sentMessage, errChannel := messenger.messageBroker.SendTemporaryMessage(
		&MessageDestination{
			DestinationID:   channelID,
			DestinationType: DestinationChannel,
		},
		messages.NewDMBlockedMessage(messenger.localizer("en"), userID),
	)
	select {
	case err := <-errChannel:
		return nil, err
	default:
		go func() {
			err, ok := <-errChannel
			if ok && err != nil {
				log.Printf("send DM blocked message: %v", err)
			}
		}()
		return sentMessage, nil
	}
}

func (messenger *BoostRequestMessenger) localizer(langs ...string) *i18n.Localizer {
	return i18n.NewLocalizer(messenger.bundle, langs...)
}
