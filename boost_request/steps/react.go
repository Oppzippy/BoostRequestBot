package steps

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

const (
	AcceptEmoji = "👍"
	StealEmoji  = "⏭"
)

type reactStep struct {
	discord *discordgo.Session
	br      *repository.BoostRequest
}

func NewReactStep(discord *discordgo.Session, br *repository.BoostRequest) *reactStep {
	return &reactStep{
		discord: discord,
		br:      br,
	}
}

func (step *reactStep) Apply() (RevertFunction, error) {
	emojis := []string{}
	if step.br.Channel.UsesBuyerMessage {
		// We can't add buttons to someone else's message, so we fall back to using reactions for those boost requests
		emojis = []string{AcceptEmoji, StealEmoji}
	}
	revert, err := step.applyReactions(step.br.Channel.BackendChannelID, step.br.BackendMessageID, emojis)
	if err != nil {
		return revert, fmt.Errorf("reacting to boost request: %v", err)
	}
	return revert, nil
}

func (step *reactStep) applyReactions(channelID, messageID string, emojis []string) (RevertFunction, error) {
	var i int
	var emoji string
	var err error
	for i, emoji = range emojis {
		err = step.discord.MessageReactionAdd(channelID, messageID, emoji)
		if err != nil {
			// Make sure we don't revert the emoji that failed to apply
			i--
			break
		}
	}
	return func() error {
		if i >= 0 {
			return step.removeReactions(channelID, messageID, emojis[0:i])
		}
		return nil
	}, err
}

func (step *reactStep) removeReactions(channelID, messageID string, emojis []string) error {
	for _, emoji := range emojis {
		err := step.discord.MessageReactionRemove(channelID, messageID, emoji, "@me")
		if err != nil {
			return err
		}
	}
	return nil
}
