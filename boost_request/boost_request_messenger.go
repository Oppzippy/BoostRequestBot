package boost_request

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type BoostRequestMessenger struct {
	Destroyed bool
	waitGroup *sync.WaitGroup
	quit      chan struct{}
}

var footer = &discordgo.MessageEmbedFooter{
	Text:    "Huokan Boosting Community",
	IconURL: "https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png",
}

func NewBoostRequestMessenger() *BoostRequestMessenger {
	brm := BoostRequestMessenger{
		Destroyed: false,
		waitGroup: new(sync.WaitGroup),
		quit:      make(chan struct{}),
	}
	return &brm
}

func (messenger *BoostRequestMessenger) SendBackendSignupMessage(discord *discordgo.Session, br *repository.BoostRequest) (*discordgo.Message, error) {
	message, err := discord.ChannelMessageSendEmbed(br.Channel.BackendChannelID, &discordgo.MessageEmbed{
		Color:       0x0000FF,
		Title:       "New Boost Request",
		Description: br.Message,
		Footer:      footer,
		Timestamp:   time.Now().Format(time.RFC3339),
	})

	if err != nil {
		return nil, err
	}

	discord.MessageReactionAdd(message.ChannelID, message.ID, ACCEPT_EMOJI)
	discord.MessageReactionAdd(message.ChannelID, message.ID, STEAL_EMOJI)

	return message, err
}

func (messenger *BoostRequestMessenger) SendBoostRequestCreatedDM(discord *discordgo.Session, br *repository.BoostRequest) (*discordgo.Message, error) {
	requester, err := discord.User(br.RequesterID)
	if err != nil {
		return nil, err
	}
	dmChannel, _ := discord.UserChannelCreate(requester.ID)
	message, err := discord.ChannelMessageSendComplex(dmChannel.ID, &discordgo.MessageSend{
		Content: "Please wait while we find an advertiser to complete your request.",
		Embed: &discordgo.MessageEmbed{
			Title: "Huokan Community Boost Request",
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
			messenger.sendTemporaryMessage(discord, br.Channel.FrontendChannelID, requester.Mention()+", I can't DM you. Please allow DMs from server members by right clicking the server and enabling \"Allow direct messages from server members.\" in Privacy Settings.")
		}
		return nil, err
	}
	return message, err
}

func (messenger *BoostRequestMessenger) SendBackendAdvertiserChosenMessage(discord *discordgo.Session, br *repository.BoostRequest) (*discordgo.Message, error) {
	advertiser, err := discord.User(br.AdvertiserID)

	if err != nil {
		return nil, err
	}

	message, err := discord.ChannelMessageSendEmbed(br.Channel.BackendChannelID, &discordgo.MessageEmbed{
		Color:       0xFF0000,
		Title:       "An advertiser has been selected.",
		Description: advertiser.Mention() + " will handle the following boost request.\n**Boost Request**\n" + br.Message,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: advertiser.AvatarURL(""),
		},
	})

	return message, err
}

func (messenger *BoostRequestMessenger) SendAdvertiserChosenDMToRequester(discord *discordgo.Session, br *repository.BoostRequest) (*discordgo.Message, error) {
	requester, err := discord.User(br.RequesterID)
	if err != nil {
		return nil, err
	}
	advertiser, err := discord.User(br.AdvertiserID)
	if err != nil {
		return nil, err
	}
	dmChannel, err := discord.UserChannelCreate(requester.ID)
	if err != nil {
		restErr, ok := err.(discordgo.RESTError)
		if ok && restErr.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			messenger.sendTemporaryMessage(discord, br.Channel.FrontendChannelID, requester.Mention()+", I can't DM you. Please allow DMs from server members by right clicking the server and enabling \"Allow direct messages from server members.\" in Privacy Settings and post your message again.")
		}
		return nil, err
	}

	message, err := discord.ChannelMessageSendEmbed(dmChannel.ID, &discordgo.MessageEmbed{
		Color:       0x00FF00,
		Title:       "Huokan Boosting Community Boost Request",
		Description: advertiser.Mention() + " (" + advertiser.String() + ") will reach out to you shortly. Anyone else that messages you regarding this boost request is not from Huokan and may attempt to scam you.",
		Footer:      footer,
		Timestamp:   time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: advertiser.AvatarURL(""),
		},
	})
	return message, err
}

func (messenger *BoostRequestMessenger) SendAdvertiserChosenDMToAdvertiser(discord *discordgo.Session, br *repository.BoostRequest) (*discordgo.Message, error) {
	requester, err := discord.User(br.RequesterID)
	if err != nil {
		return nil, err
	}
	advertiser, err := discord.User(br.AdvertiserID)
	if err != nil {
		return nil, err
	}
	dmChannel, err := discord.UserChannelCreate(advertiser.ID)
	if err != nil {
		restErr, ok := err.(discordgo.RESTError)
		if ok && restErr.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			messenger.sendTemporaryMessage(discord, br.Channel.BackendChannelID, advertiser.Mention()+", I can't DM you. Please allow DMs from server members by right clicking the server and enabling \"Allow direct messages from server members.\" in Privacy Settings.")
			_, err := discord.ChannelMessageSend(br.Channel.BackendChannelID, "Please DM "+requester.Mention()+" ("+requester.String()+").")
			if err != nil {
				log.Println("Failed to send backup message after failed DM", err)
			}
		}
		return nil, err
	}

	message, err := discord.ChannelMessageSendEmbed(dmChannel.ID, &discordgo.MessageEmbed{
		Color:       0xFF0000,
		Title:       "You have been selected to handle a boost request.",
		Description: "Please message " + requester.Mention() + " (" + requester.String() + ")\n**Boost Request**\n" + br.Message,
		Footer:      footer,
		Timestamp:   time.Now().Format(time.RFC3339),
	})

	return message, err
}

// Logs the creation of a boost request to a channel only moderators can view
func (messenger *BoostRequestMessenger) SendLogChannelMessage(
	discord *discordgo.Session, br *repository.BoostRequest, channelID string,
) (*discordgo.Message, error) {
	user, err := discord.User(br.RequesterID)
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
	message, err := discord.ChannelMessageSendEmbed(channelID, embed)
	return message, err
}

func (messenger *BoostRequestMessenger) Destroy(discord *discordgo.Session) {
	if !messenger.Destroyed {
		messenger.Destroyed = true
		close(messenger.quit)
		messenger.waitGroup.Wait()
	}
}

func (messenger *BoostRequestMessenger) sendTemporaryMessage(discord *discordgo.Session, channelID string, content string) {
	message, err := discord.ChannelMessageSend(channelID, content)
	if err == nil {
		messenger.waitGroup.Add(1)
		go func() {
			select {
			case <-time.After(30 * time.Second):
			case <-messenger.quit:
			}
			err := discord.ChannelMessageDelete(message.ChannelID, message.ID)
			if err != nil {
				log.Println("Error deleting temporary message", err)
			}
			messenger.waitGroup.Done()
		}()
	} else {
		log.Println("Error sending temporary message", err)
	}
}
