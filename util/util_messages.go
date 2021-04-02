package util

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func NotifyDMsDisabled(discord *discordgo.Session, channelID string, userID string) {
	go func() {
		message, err := discord.ChannelMessageSend(channelID, "<@"+userID+">, I can't DM you. Please allow DMs from server members by right clicking the server and enabling \"Allow direct messages from server members.\" in Privacy Settings and try again.")
		if err != nil {
			time.Sleep(30 * time.Second)
			discord.ChannelMessageDelete(message.ChannelID, message.ID)
		}
	}()
}
