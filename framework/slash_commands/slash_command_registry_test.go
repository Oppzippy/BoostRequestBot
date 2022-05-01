package slash_commands_test

import (
	"io"
	"log"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"github.com/oppzippy/BoostRequestBot/framework/slash_commands"
	"github.com/oppzippy/BoostRequestBot/framework/slash_commands/mock_slash_commands"
	"github.com/oppzippy/BoostRequestBot/util/testing_util"
)

func TestSlashCommandRegistryCommandMatching(t *testing.T) {
	t.Parallel()

	type test struct {
		name            string
		commandPath     []string
		executedCommand []string
		shouldSucceed   bool
		options         []*discordgo.ApplicationCommandInteractionDataOption
		expectedOptions map[string]interface{}
	}

	tests := []test{
		{
			name:            "basic command",
			commandPath:     []string{"test"},
			executedCommand: []string{"test"},
			shouldSucceed:   true,
		},
		{
			name:            "1 level deep subcommands",
			commandPath:     []string{"test subcommand1"},
			executedCommand: []string{"test subcommand1"},
			shouldSucceed:   true,
		},
		{
			name:            "2 level deep subcommands",
			commandPath:     []string{"test subcommand1 subcommand2"},
			executedCommand: []string{"test subcommand1 subcommand2"},
			shouldSucceed:   true,
		},
		{
			name:            "wrong command name",
			commandPath:     []string{"test"},
			executedCommand: []string{"test2"},
			shouldSucceed:   false,
		},
		{
			name:            "subcommand not deep enough",
			commandPath:     []string{"test subcommand1 subcommand2"},
			executedCommand: []string{"test subcommand1"},
			shouldSucceed:   false,
		},
		{
			name:            "subcommand too deep",
			commandPath:     []string{"test subcommand1"},
			executedCommand: []string{"test subcommand1 subcommand2"},
			shouldSucceed:   false,
		},
		{
			name:            "command with options",
			commandPath:     []string{"test"},
			executedCommand: []string{"test"},
			options: []*discordgo.ApplicationCommandInteractionDataOption{
				{
					Name:  "option",
					Type:  discordgo.ApplicationCommandOptionBoolean,
					Value: true,
				},
			},
			expectedOptions: map[string]interface{}{
				"option": true,
			},
			shouldSucceed: true,
		},
		{
			name:            "subcommand with options",
			commandPath:     []string{"test subcommand"},
			executedCommand: []string{"test subcommand"},
			options: []*discordgo.ApplicationCommandInteractionDataOption{
				{
					Name:  "option",
					Type:  discordgo.ApplicationCommandOptionBoolean,
					Value: true,
				},
			},
			expectedOptions: map[string]interface{}{
				"option": true,
			},
			shouldSucceed: true,
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			registry := slash_commands.NewSlashCommandRegistry()
			registry.Logger = log.New(io.Discard, "", 0)

			expectedResponse := "success!"
			registry.RegisterCommand(testCase.commandPath, func(interaction *discordgo.InteractionCreate, options map[string]*discordgo.ApplicationCommandInteractionDataOption) (*discordgo.InteractionResponse, error) {
				if testCase.expectedOptions != nil {
					for _, option := range options {
						if testCase.expectedOptions[option.Name] != option.Value {
							t.Errorf("missing option %v", option.Name)
						}
					}
				}
				return &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: expectedResponse,
					},
				}, nil
			})

			responder := mock_slash_commands.NewMockInteractionResponder(ctrl)
			if testCase.shouldSucceed {
				responder.EXPECT().InteractionRespond(gomock.Any(), &responseMatcher{expectedContent: expectedResponse}).Return(nil)
			}
			registry.OnInteraction(responder, &discordgo.InteractionCreate{
				Interaction: &discordgo.Interaction{
					Type: discordgo.InteractionApplicationCommand,
					Data: testing_util.CommandPathToInteractionData(testCase.executedCommand, testCase.options),
				},
			})
		})
	}

}

type responseMatcher struct {
	expectedContent string
}

func (m *responseMatcher) Matches(x interface{}) bool {
	response, ok := x.(*discordgo.InteractionResponse)
	if !ok {
		return false
	}
	return response.Data.Content == m.expectedContent
}

func (m *responseMatcher) String() string {
	return "compares response content to a constant"
}
