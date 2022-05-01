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

func CommandPathToString(commandData *discordgo.ApplicationCommandInteractionData) string {
	sb := strings.Builder{}
	sb.WriteString(commandData.Name)

	for option := commandData.Options[0]; isOptionSubCommandOrGroup(option); option = option.Options[0] {
		sb.WriteRune('.')
		sb.WriteString(option.Name)
	}

	return sb.String()
}

func isOptionSubCommandOrGroup(option *discordgo.ApplicationCommandInteractionDataOption) bool {
	if option == nil {
		return false
	}
	return option.Type == discordgo.ApplicationCommandOptionSubCommand || option.Type == discordgo.ApplicationCommandOptionSubCommandGroup
}
