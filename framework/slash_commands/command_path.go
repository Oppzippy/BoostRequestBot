package slash_commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MatchesCommandPath(commandData discordgo.ApplicationCommandInteractionData, path []string) bool {
	if len(path) == 0 {
		return true
	}
	commandName, subCommandPath := path[0], path[1:]
	if commandData.Name != commandName {
		return false
	}

	option := commandData.Options[0]
	for i, name := range subCommandPath {
		isCommand := option.Type == discordgo.ApplicationCommandOptionSubCommandGroup ||
			option.Type == discordgo.ApplicationCommandOptionSubCommand
		if !isCommand || option.Name != name {
			return false
		}
		if len(option.Options) >= 1 {
			option = option.Options[0]
		} else if i != len(subCommandPath)-1 {
			return false
		}
	}

	return true
}

func getCommandPathString(commandData *discordgo.ApplicationCommandInteractionData) string {
	return commandPathToString(getCommandPath(commandData))
}

func commandPathToString(path []string) string {
	return strings.Join(path, ".")
}

func getCommandPath(commandData *discordgo.ApplicationCommandInteractionData) []string {
	path := []string{commandData.Name}
	if len(commandData.Options) == 0 {
		return path
	}

	for option := commandData.Options[0]; isOptionSubCommandOrGroup(option); {
		path = append(path, option.Name)

		if len(option.Options) > 0 {
			option = option.Options[0]
		} else {
			option = nil
		}
	}
	return path
}

func isOptionSubCommandOrGroup(option *discordgo.ApplicationCommandInteractionDataOption) bool {
	if option == nil {
		return false
	}
	return option.Type == discordgo.ApplicationCommandOptionSubCommand || option.Type == discordgo.ApplicationCommandOptionSubCommandGroup
}
