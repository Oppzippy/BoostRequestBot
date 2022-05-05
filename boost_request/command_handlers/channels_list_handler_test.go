package command_handlers_test

import (
	"strings"
	"testing"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/mock/gomock"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/command_handlers"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository/mock_repository"
	"golang.org/x/text/language"
)

func TestChannelsListHandler_Handle(t *testing.T) {
	t.Parallel()
	type Test struct {
		name             string
		channels         []*repository.BoostRequestChannel
		expectedResponse string
	}
	tests := []Test{
		{
			name:             "no channels",
			channels:         []*repository.BoostRequestChannel{},
			expectedResponse: "There are no boost request channels.",
		},
		{
			name: "boost request channel with no options",
			channels: []*repository.BoostRequestChannel{
				{
					ID:                0,
					GuildID:           "",
					FrontendChannelID: "123",
					BackendChannelID:  "456",
					UsesBuyerMessage:  false,
					SkipsBuyerDM:      false,
				},
			},
			expectedResponse: "**Channel 1**\nFrontend Channel: <#123>\nBackend Channel: <#456>\nOptions: none",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			bundle := i18n.NewBundle(language.AmericanEnglish)
			repo := mock_repository.NewMockRepository(ctrl)
			handler := command_handlers.NewChannelsListHandler(bundle, repo)

			repo.EXPECT().GetBoostRequestChannels("GuildID").Return(test.channels, nil)

			response, err := handler.Handle(
				&discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						GuildID: "GuildID",
						Locale:  "",
					},
				},
				map[string]*discordgo.ApplicationCommandInteractionDataOption{},
			)
			if err != nil {
				t.Errorf("error handling command: %v", err)
			}

			if strings.TrimSpace(response.Data.Content) != test.expectedResponse {
				t.Errorf("expected %v, got %v", test.expectedResponse, response.Data.Content)
			}
		})
	}
}
