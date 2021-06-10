package message_generator

import "github.com/bwmarrin/discordgo"

type DiscordRoleNameProvider struct {
	discord *discordgo.Session
}

func NewDiscordRoleNameProvider(discord *discordgo.Session) *DiscordRoleNameProvider {
	return &DiscordRoleNameProvider{
		discord: discord,
	}
}

func (rnp *DiscordRoleNameProvider) RoleName(guildID, roleID string) string {
	guild, err := rnp.discord.State.Guild(guildID)

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
