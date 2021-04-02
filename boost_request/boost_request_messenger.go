package boost_request

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type BoostRequestMessenger struct {
	waitGroup               sync.WaitGroup
	messagesPendingDeletion sync.Map
	destroyed               bool
}

var FOOTER = &discordgo.MessageEmbedFooter{
	Text:    "Huokan Boosting Community",
	IconURL: "https://cdn.discordapp.com/attachments/721652505796411404/749063535719481394/HuokanLogoCropped.png",
}

func (messenger *BoostRequestMessenger) SendBackendSignupMessage(discord *discordgo.Session, br *BoostRequest) (*discordgo.Message, error) {
	message, err := discord.ChannelMessageSendEmbed(br.Channel.BackendChannelID, &discordgo.MessageEmbed{
		Color:       0x0000FF,
		Title:       "New Boost Request",
		Description: br.Message,
		Footer:      FOOTER,
		Timestamp:   time.Now().Format(time.RFC3339),
	})

	return message, err
}

func (messenger *BoostRequestMessenger) SendBoostRequestCreatedDM(discord *discordgo.Session, br *BoostRequest) (*discordgo.Message, error) {
	requester, err := discord.User(br.RequesterID)
	if err != nil {
		return nil, err
	}
	dmChannel, err := discord.UserChannelCreate(requester.ID)
	if err != nil {
		restErr, ok := err.(discordgo.RESTError)
		if ok && restErr.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			messenger.sendTemporaryMessage(discord, br.Channel.FrontendChannelID, requester.Mention()+", I can't DM you. Please allow DMs from server members by right clicking the server and enabling \"Allow direct messages from server members.\" in Privacy Settings.")
		}
		return nil, err
	}
	message, err := discord.ChannelMessageSendComplex(dmChannel.ID, &discordgo.MessageSend{
		Content: "Please wait while we find an advertiser to complete your request.",
		Embed: &discordgo.MessageEmbed{
			Title: "Huokan Community Boost Request",
			Author: &discordgo.MessageEmbedAuthor{
				Name: requester.String(),
			},
			Description: br.Message,
			Footer:      FOOTER,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: requester.AvatarURL(""),
			},
		},
	})
	return message, err
}

func (messenger *BoostRequestMessenger) SendBackendAdvertiserChosenMessage(discord *discordgo.Session, br *BoostRequest) (*discordgo.Message, error) {
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

func (messenger *BoostRequestMessenger) SendAdvertiserChosenDMToRequester(discord *discordgo.Session, br *BoostRequest) (*discordgo.Message, error) {
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
		Footer:      FOOTER,
		Timestamp:   time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: advertiser.AvatarURL(""),
		},
	})
	return message, err
}

func (messenger *BoostRequestMessenger) SendAdvertiserChosenDMToAdvertiser(discord *discordgo.Session, br *BoostRequest) (*discordgo.Message, error) {
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
		Footer:      FOOTER,
		Timestamp:   time.Now().Format(time.RFC3339),
	})

	return message, err
}

func (messenger *BoostRequestMessenger) Destroy(discord *discordgo.Session) {
	messenger.destroyed = true
	messenger.waitGroup.Wait()
	messenger.messagesPendingDeletion.Range(func(_, value interface{}) bool {
		message, ok := value.(*discordgo.Message)
		if ok {
			discord.ChannelMessageDelete(message.ChannelID, message.ID)
		}
		return true
	})
}

func (messenger *BoostRequestMessenger) sendTemporaryMessage(discord *discordgo.Session, channelID string, content string) {
	message, err := discord.ChannelMessageSend(channelID, content)
	if err == nil {
		messenger.messagesPendingDeletion.Store(message.ID, message)
		go func() {
			time.Sleep(30 * time.Second)
			messenger.waitGroup.Add(1)
			defer messenger.waitGroup.Done()
			if !messenger.destroyed {
				err := discord.ChannelMessageDelete(message.ChannelID, message.ID)
				if err != nil {
					log.Println("Error deleting temporary message", err)
				}
				messenger.messagesPendingDeletion.Delete(message.ID)
			}
		}()
	} else {
		log.Println("Error sending temporary message", err)
	}
}
