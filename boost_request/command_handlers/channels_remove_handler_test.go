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

func TestChannelsRemoveHandler_Handle(t *testing.T) {
	t.Parallel()
	tests := []*testCase{
		{
			name:             "delete a channel",
			expectedResponse: "Removed boost request channel <#123>.",
			interaction:      &discordgo.Interaction{GuildID: "GuildID"},
			options: map[string]*discordgo.ApplicationCommandInteractionDataOption{
				"frontend-channel": {
					Name:  "frontend-channel",
					Type:  discordgo.ApplicationCommandOptionChannel,
					Value: "123",
				},
			},
			setup: func(t *testing.T) command_handlers.CommandHandler {
				ctrl := gomock.NewController(t)
				repo := mock_repository.NewMockRepository(ctrl)
				bundle := i18n.NewBundle(language.AmericanEnglish)

				channel := &repository.BoostRequestChannel{
					ID:                2,
					GuildID:           "GuildID",
					FrontendChannelID: "123",
					BackendChannelID:  "456",
					UsesBuyerMessage:  false,
					SkipsBuyerDM:      false,
				}

				repo.EXPECT().GetBoostRequestChannelByFrontendChannelID("GuildID", "123").Return(channel, nil)
				repo.EXPECT().DeleteBoostRequestChannel(&boostRequestChannelMatcher{channel}).Return(nil)

				return command_handlers.NewChannelsRemoveHandler(bundle, repo)
			},
		},
		{
			name:             "delete a nonexistent channel",
			expectedResponse: "<#1> is not a boost request frontend.",
			interaction:      &discordgo.Interaction{GuildID: "GuildID"},
			options: map[string]*discordgo.ApplicationCommandInteractionDataOption{
				"frontend-channel": {
					Name:  "frontend-channel",
					Type:  discordgo.ApplicationCommandOptionChannel,
					Value: "1",
				},
			},
			setup: func(t *testing.T) command_handlers.CommandHandler {
				ctrl := gomock.NewController(t)
				repo := mock_repository.NewMockRepository(ctrl)
				bundle := i18n.NewBundle(language.AmericanEnglish)

				repo.EXPECT().GetBoostRequestChannelByFrontendChannelID("GuildID", "1").Return(nil, repository.ErrNoResults)

				return command_handlers.NewChannelsRemoveHandler(bundle, repo)
			},
		},
	}
	testCommandHandler(t, tests)
}

type boostRequestChannelMatcher struct {
	channelToMatch *repository.BoostRequestChannel
}

func (m *boostRequestChannelMatcher) Matches(channel interface{}) bool {
	brc, ok := channel.(*repository.BoostRequestChannel)
	if !ok {
		return false
	}
	return brc.ID == m.channelToMatch.ID
}

func (m *boostRequestChannelMatcher) String() string {
	return "Matches a boost request channel"
}
