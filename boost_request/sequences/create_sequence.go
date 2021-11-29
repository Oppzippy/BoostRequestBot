package sequences

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/oppzippy/BoostRequestBot/boost_request/active_request"
	"github.com/oppzippy/BoostRequestBot/boost_request/messenger"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/oppzippy/BoostRequestBot/boost_request/steps"
)

type CreateSequenceArgs struct {
	Repository               repository.Repository
	BoostRequest             *repository.BoostRequest
	Discord                  *discordgo.Session
	Messenger                *messenger.BoostRequestMessenger
	ActiveRequests           *sync.Map
	BackendMessageChannelIDs []string
	SetWinnerCallback        func(*active_request.AdvertiserChosenEvent)
}

func RunCreateHumanRequesterSequence(args CreateSequenceArgs) error {
	steps := []steps.RevertableStep{
		steps.NewSendCreatedDMStep(args.Discord, *args.Messenger, args.BoostRequest),
		steps.NewSendMessageStep(args.Discord, args.Messenger, args.BoostRequest, args.BackendMessageChannelIDs),
		steps.NewInsertBoostRequestStep(args.Repository, args.BoostRequest),
		steps.NewStoreActiveRequestStep(args.ActiveRequests, args.BoostRequest, args.SetWinnerCallback),
		steps.NewReactStep(args.Discord, args.BoostRequest),
		steps.NewPostToLogChannelStep(args.Repository, args.BoostRequest, args.Messenger),
	}

	_, err := runSequence(steps)
	return err
}

func RunCreateBotRequesterSequence(args CreateSequenceArgs) error {
	steps := []steps.RevertableStep{
		steps.NewSendMessageStep(args.Discord, args.Messenger, args.BoostRequest, args.BackendMessageChannelIDs),
		steps.NewInsertBoostRequestStep(args.Repository, args.BoostRequest),
		steps.NewStoreActiveRequestStep(args.ActiveRequests, args.BoostRequest, args.SetWinnerCallback),
		steps.NewReactStep(args.Discord, args.BoostRequest),
		steps.NewPostToLogChannelStep(args.Repository, args.BoostRequest, args.Messenger),
	}

	_, err := runSequence(steps)
	return err
}
