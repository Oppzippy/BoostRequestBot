package command_handlers_test

import (
	"testing"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"

	"github.com/golang/mock/gomock"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/command_handlers"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository/mock_repository"
	"golang.org/x/text/language"

	"github.com/bwmarrin/discordgo"
)

func TestCreditsListHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []*testCase{
		{
			name:             "nonexistent user",
			expectedResponse: "<@123> has 0 steal credits.",
			interaction:      &discordgo.Interaction{GuildID: "GuildID"},
			options: map[string]*discordgo.ApplicationCommandInteractionDataOption{
				"user": {
					Name:  "user",
					Type:  discordgo.ApplicationCommandOptionUser,
					Value: "123",
				},
			},
			setup: func(t *testing.T) command_handlers.CommandHandler {
				ctrl := gomock.NewController(t)
				bundle := i18n.NewBundle(language.AmericanEnglish)
				repo := mock_repository.NewMockRepository(ctrl)

				repo.EXPECT().GetStealCreditsForUser("GuildID", "123").Return(0, repository.ErrNoResults)

				return command_handlers.NewCreditsListHandler(bundle, repo)
			},
		},
		{
			name:             "user exists",
			expectedResponse: "<@123> has 2 steal credits.",
			interaction:      &discordgo.Interaction{GuildID: "GuildID"},
			options: map[string]*discordgo.ApplicationCommandInteractionDataOption{
				"user": {
					Name:  "user",
					Type:  discordgo.ApplicationCommandOptionUser,
					Value: "123",
				},
			},
			setup: func(t *testing.T) command_handlers.CommandHandler {
				ctrl := gomock.NewController(t)
				bundle := i18n.NewBundle(language.AmericanEnglish)
				repo := mock_repository.NewMockRepository(ctrl)

				repo.EXPECT().GetStealCreditsForUser("GuildID", "123").Return(2, nil)

				return command_handlers.NewCreditsListHandler(bundle, repo)
			},
		},
	}
	testCommandHandler(t, tests)
}
