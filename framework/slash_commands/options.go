package slash_commands

import "github.com/bwmarrin/discordgo"

func parseOptions(commandData *discordgo.ApplicationCommandInteractionData) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := getInnermostCommandOptions(commandData)
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}
	return optionMap
}

func getInnermostCommandOptions(commandData *discordgo.ApplicationCommandInteractionData) []*discordgo.ApplicationCommandInteractionDataOption {
	option := commandData.Options[0]
	if !isOptionSubCommandOrGroup(option) {
		return commandData.Options
	}

	for isOptionSubCommandOrGroup(option.Options[0]) {
		option = option.Options[0]
	}
	return option.Options
}
