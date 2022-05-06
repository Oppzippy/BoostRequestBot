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

const channelIDToRemove = "1"

func TestChannelsRemoveHandler_Handle(t *testing.T) {
	t.Parallel()
	type Test struct {
		name             string
		channel          *repository.BoostRequestChannel
		expectedResponse string
	}
	tests := []Test{
		{
			name: "delete a channel",
			channel: &repository.BoostRequestChannel{
				GuildID:           "GuildID",
				FrontendChannelID: "1",
			},
			expectedResponse: "Removed boost request channel <#1>.",
		},
		{
			name:             "delete a nonexistent channel",
			channel:          nil,
			expectedResponse: "<#1> is not a boost request frontend.",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			bundle := i18n.NewBundle(language.AmericanEnglish)
			repo := mock_repository.NewMockRepository(ctrl)
			handler := command_handlers.NewChannelsRemoveHandler(bundle, repo)

			if test.channel != nil {
				repo.EXPECT().GetBoostRequestChannelByFrontendChannelID("GuildID", channelIDToRemove).Return(test.channel, nil)
				repo.EXPECT().DeleteBoostRequestChannel(&boostRequestChannelMatcher{test.channel}).Return(nil).AnyTimes()
			} else {
				repo.EXPECT().GetBoostRequestChannelByFrontendChannelID("GuildID", channelIDToRemove).Return(nil, repository.ErrNoResults)
			}

			response, err := handler.Handle(
				&discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						GuildID: "GuildID",
						Locale:  "",
					},
				},
				map[string]*discordgo.ApplicationCommandInteractionDataOption{
					"frontend-channel": {
						Name:    "frontend-channel",
						Type:    discordgo.ApplicationCommandOptionChannel,
						Value:   channelIDToRemove,
						Options: nil,
						Focused: false,
					},
				},
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
