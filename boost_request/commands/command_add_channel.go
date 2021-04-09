package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var addChannelCommand = dgc.Command{
	Name:        "addchannel",
	Description: "Registers a channel for the bot to watch for boost requests.",
	Usage:       "!boostrequest addchannel <#frontend-channel> <#backend-channel>",
	Example:     "!boostrequest addchannel #boost-request-frontend #boost-request-backend",
	IgnoreCase:  true,
	Handler:     addChannelHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func addChannelHandler(ctx *dgc.Ctx) {
	if ctx.Arguments.Amount() != 2 {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}
	frontendID := ctx.Arguments.Get(0).AsChannelMentionID()
	backendID := ctx.Arguments.Get(1).AsChannelMentionID()

	if frontendID == "" || backendID == "" {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}

	frontend, err := ctx.Session.Channel(frontendID)
	if err != nil {
		log.Println("Error fetching channel", err)
		respondText(ctx, "An error has occurred.")
		return
	}
	backend, err := ctx.Session.Channel(backendID)
	if err != nil {
		log.Println("Error fetching channel", err)
		respondText(ctx, "An error has occurred.")
		return
	}

	if frontend.Type != discordgo.ChannelTypeGuildText || backend.Type != discordgo.ChannelTypeGuildText {
		respondText(ctx, "Frontend and backend channels must both be text channels.")
		return
	}

	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	err = repo.InsertBoostRequestChannel(&repository.BoostRequestChannel{
		GuildID:           ctx.Event.GuildID,
		FrontendChannelID: frontend.ID,
		BackendChannelID:  backend.ID,
		UsesBuyerMessage:  frontend.ID == backend.ID,
		SkipsBuyerDM:      frontend.ID == backend.ID,
	})

	if err != nil {
		log.Println("Error inserting boost request channel", err)
		respondText(ctx, "An error has occurred.")
		return
	}
	respondText(
		ctx,
		fmt.Sprintf(
			"Added boost request frontend %s with backend %s",
			frontend.Mention(),
			backend.Mention(),
		),
	)
}
