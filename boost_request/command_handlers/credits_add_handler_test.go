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
	tests := []*testCase{
		{
			name:             "add credits",
			expectedResponse: "Added 2 steal credits. New total is 4.",
			interaction:      &discordgo.Interaction{GuildID: "GuildID"},
			options: map[string]*discordgo.ApplicationCommandInteractionDataOption{
				"user": {
					Name:    "user",
					Type:    discordgo.ApplicationCommandOptionUser,
					Value:   "UserID",
					Options: nil,
					Focused: false,
				},
				"credits": {
					Name:    "credits",
					Type:    discordgo.ApplicationCommandOptionInteger,
					Value:   float64(2),
					Options: nil,
					Focused: false,
				},
			},
			setup: func(t *testing.T) command_handlers.CommandHandler {
				ctrl := gomock.NewController(t)
				bundle := i18n.NewBundle(language.AmericanEnglish)
				repo := mock_repository.NewMockRepository(ctrl)

				repo.EXPECT().AdjustStealCreditsForUser("GuildID", "UserID", repository.OperationAdd, 2)
				repo.EXPECT().GetStealCreditsForUser("GuildID", "UserID").Return(4, nil)

				return command_handlers.NewCreditsAddHandler(bundle, repo)
			},
		},
		{
			name:             "subtract credits",
			expectedResponse: "Added -2 steal credits. New total is 3.",
			interaction:      &discordgo.Interaction{GuildID: "GuildID"},
			options: map[string]*discordgo.ApplicationCommandInteractionDataOption{
				"user": {
					Name:  "user",
					Type:  discordgo.ApplicationCommandOptionUser,
					Value: "UserID",
				},
				"credits": {
					Name:  "credits",
					Type:  discordgo.ApplicationCommandOptionInteger,
					Value: float64(-2),
				},
			},
			setup: func(t *testing.T) command_handlers.CommandHandler {
				ctrl := gomock.NewController(t)
				bundle := i18n.NewBundle(language.AmericanEnglish)
				repo := mock_repository.NewMockRepository(ctrl)

				repo.EXPECT().AdjustStealCreditsForUser("GuildID", "UserID", repository.OperationAdd, -2)
				repo.EXPECT().GetStealCreditsForUser("GuildID", "UserID").Return(3, nil)

				return command_handlers.NewCreditsAddHandler(bundle, repo)
			},
		},
	}
	testCommandHandler(t, tests)
}
