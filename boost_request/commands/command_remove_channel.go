package commands

import (
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var removeChannelCommand = dgc.Command{
	Name:        "removechannel",
	Description: "Stops watching (unregisters) a boost request channel.",
	Usage:       "!boostrequest removechannel <#frontend-channel>",
	Example:     "!boostrequest removechannel #boost-request-frontend",
	IgnoreCase:  true,
	Handler:     removeChannelHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func removeChannelHandler(ctx *dgc.Ctx) {
	if ctx.Arguments.Amount() != 1 {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}
	channelID := ctx.Arguments.Get(0).AsChannelMentionID()

	if channelID == "" {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}

	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	brc, err := repo.GetBoostRequestChannelByFrontendChannelID(ctx.Event.GuildID, channelID)
	if err != nil {
		if err == repository.ErrBoostRequestChannelNotFound {
			respondText(ctx, "<#"+channelID+"> is not a boost request frontend.")
			return
		}
		log.Printf("Error fetching boost request channel: %v", err)
		respondText(ctx, genericError)
		return
	}
	err = repo.DeleteBoostRequestChannel(brc)
	if err != nil {
		log.Printf("Error deleting boost request channel: %v", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Removed boost request channel <#"+channelID+">")
}
