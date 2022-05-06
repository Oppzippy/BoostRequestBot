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

func TestChannelsListHandler_Handle(t *testing.T) {
	t.Parallel()
	interaction := &discordgo.Interaction{
		GuildID: "GuildID",
		Locale:  "",
	}
	tests := []*testCase{
		{
			name:             "no channels",
			expectedResponse: "There are no boost request channels.",
			interaction:      interaction,
			options:          nil,
			setup: func(t *testing.T) command_handlers.CommandHandler {
				ctrl := gomock.NewController(t)
				repo := mock_repository.NewMockRepository(ctrl)
				bundle := i18n.NewBundle(language.AmericanEnglish)

				repo.EXPECT().GetBoostRequestChannels("GuildID").Return([]*repository.BoostRequestChannel{}, nil)
				return command_handlers.NewChannelsListHandler(bundle, repo)
			},
		},
		{
			name:             "boost request channel with no options",
			expectedResponse: "**Channel 1**\nFrontend Channel: <#123>\nBackend Channel: <#456>\nOptions: none",
			interaction:      interaction,
			options:          nil,
			setup: func(t *testing.T) command_handlers.CommandHandler {
				ctrl := gomock.NewController(t)
				repo := mock_repository.NewMockRepository(ctrl)
				bundle := i18n.NewBundle(language.AmericanEnglish)

				repo.EXPECT().GetBoostRequestChannels("GuildID").Return([]*repository.BoostRequestChannel{
					{
						ID:                0,
						GuildID:           "GuildID",
						FrontendChannelID: "123",
						BackendChannelID:  "456",
						UsesBuyerMessage:  false,
						SkipsBuyerDM:      false,
					},
				}, nil)
				return command_handlers.NewChannelsListHandler(bundle, repo)
			},
		},
	}

	testCommandHandler(t, tests)
}
