package command_handlers_test

import (
	"testing"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"

	"github.com/bwmarrin/discordgo"

	"github.com/golang/mock/gomock"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/command_handlers"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository/mock_repository"
	"golang.org/x/text/language"
)

func TestCreditsAddHandler_Handle(t *testing.T) {
	t.Parallel()
	type Test struct {
		name             string
		userID           string
		creditsToAdd     int
		newCredits       int
		expectedResponse string
	}

	tests := []*Test{
		{
			name:             "add credits",
			userID:           "1",
			creditsToAdd:     2,
			newCredits:       4,
			expectedResponse: "Added 2 steal credits. New total is 4.",
		},
		{
			name:             "subtract credits",
			userID:           "1",
			creditsToAdd:     -2,
			newCredits:       3,
			expectedResponse: "Added -2 steal credits. New total is 3.",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			bundle := i18n.NewBundle(language.AmericanEnglish)
			repo := mock_repository.NewMockRepository(ctrl)
			handler := command_handlers.NewCreditsAddHandler(bundle, repo)

			repo.EXPECT().AdjustStealCreditsForUser("GuildID", test.userID, repository.OperationAdd, test.creditsToAdd)
			repo.EXPECT().GetStealCreditsForUser("GuildID", test.userID).Return(test.newCredits, nil)

			response, err := handler.Handle(&discordgo.InteractionCreate{
				Interaction: &discordgo.Interaction{
					GuildID: "GuildID",
				},
			}, map[string]*discordgo.ApplicationCommandInteractionDataOption{
				"user": {
					Name:    "user",
					Type:    discordgo.ApplicationCommandOptionUser,
					Value:   test.userID,
					Options: nil,
					Focused: false,
				},
				"credits": {
					Name:    "credits",
					Type:    discordgo.ApplicationCommandOptionInteger,
					Value:   float64(test.creditsToAdd),
					Options: nil,
					Focused: false,
				},
			})

			if err != nil {
				t.Errorf("error handing command: %v", err)
			}

			if response.Data.Content != test.expectedResponse {
				t.Errorf("expected %v, got %v", test.expectedResponse, response.Data.Content)
			}
		})
	}
}
