package commands

import (
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var setLogChannelCommand = dgc.Command{
	Name:        "setlogchannel",
	Description: "Sets a channel for all boost requests to be logged to.",
	Usage:       "!boostrequest logchannel <#logchannel>",
	Example:     "!boostrequest logchannel #boost-request-log",
	IgnoreCase:  true,
	Handler:     setLogChannelHandler,
}

func setLogChannelHandler(ctx *dgc.Ctx) {
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
	err := repo.InsertLogChannel(ctx.Event.GuildID, channelID)
	if err != nil {
		log.Println("Error setting boost request log channel", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Boost request log channel set to <#"+channelID+">.")
}
