package interactions

import "github.com/bwmarrin/discordgo"

func getGuildMember(discord *discordgo.Session, event *discordgo.InteractionCreate, guildID string) (*discordgo.Member, error) {
	if event.Member != nil {
		return event.Member, nil
	}
	member, err := discord.GuildMember(guildID, event.User.ID)
	return member, err
}
