package command_handlers_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository/mock_repository"
	"golang.org/x/text/language"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/command_handlers"
)

func TestStealCreditsSetHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []*testCase{
		{
			name:             "set credits",
			expectedResponse: "Set <@123>'s steal credits to 2.",
			interaction:      &discordgo.Interaction{GuildID: "GuildID"},
			options: map[string]*discordgo.ApplicationCommandInteractionDataOption{
				"user": {
					Name:  "user",
					Type:  discordgo.ApplicationCommandOptionUser,
					Value: "123",
				},
				"credits": {
					Name:  "credits",
					Type:  discordgo.ApplicationCommandOptionInteger,
					Value: float64(2),
				},
			},
			setup: func(t *testing.T) command_handlers.CommandHandler {
				ctrl := gomock.NewController(t)
				bundle := i18n.NewBundle(language.AmericanEnglish)
				repo := mock_repository.NewMockRepository(ctrl)

				repo.EXPECT().UpdateStealCreditsForUser("GuildID", "123", 2).Return(nil)

				return command_handlers.NewCreditsSetHandler(bundle, repo)
			},
		},
	}

	testCommandHandler(t, tests)
}
