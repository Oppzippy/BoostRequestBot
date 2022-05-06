package command_handlers_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"

	"github.com/oppzippy/BoostRequestBot/boost_request/command_handlers"
)

type testCase struct {
	name             string
	interaction      *discordgo.Interaction
	options          map[string]*discordgo.ApplicationCommandInteractionDataOption
	expectedResponse string
	setup            func(t *testing.T) command_handlers.CommandHandler
}

func testCommandHandler(t *testing.T, tests []*testCase) {
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			handler := test.setup(t)

			response, err := handler.Handle(&discordgo.InteractionCreate{
				Interaction: test.interaction,
			}, test.options)

			if err != nil {
				t.Errorf("error handing command: %v", err)
			}

			if response.Data.Content != test.expectedResponse {
				t.Errorf("expected %v, got %v", test.expectedResponse, response.Data.Content)
			}
		})
	}
}
