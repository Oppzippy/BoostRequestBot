package testing_util

import "github.com/bwmarrin/discordgo"

func CommandPathToInteractionData(path []string, options []*discordgo.ApplicationCommandInteractionDataOption) discordgo.ApplicationCommandInteractionData {
	return discordgo.ApplicationCommandInteractionData{
		Name:    path[0],
		Options: commandPathSubcommandOptions(path[1:], options),
	}
}

func commandPathSubcommandOptions(path []string, options []*discordgo.ApplicationCommandInteractionDataOption) []*discordgo.ApplicationCommandInteractionDataOption {
	if len(path) == 0 {
		return options
	}
	return []*discordgo.ApplicationCommandInteractionDataOption{
		{
			Name:    path[0],
			Type:    discordgo.ApplicationCommandOptionSubCommand,
			Options: commandPathSubcommandOptions(path, options),
		},
	}
}
