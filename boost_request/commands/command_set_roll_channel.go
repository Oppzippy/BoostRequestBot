package commands

import (
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var setRollChannelCommand = dgc.Command{
	Name:        "setrollchannel",
	Description: "Sets a channel for boost request RNG rolls to be posted to.",
	Usage:       "!boostrequest setrollchannel <#rollchannel>",
	Example:     "!boostrequest setrollchannel #rolls ",
	IgnoreCase:  true,
	Handler:     setRollChannelHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func setRollChannelHandler(ctx *dgc.Ctx) {
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
	err := repo.InsertRollChannel(ctx.Event.GuildID, channelID)
	if err != nil {
		log.Printf("Error setting boost request roll channel: %v", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Boost request roll channel set to <#"+channelID+">.")
}
