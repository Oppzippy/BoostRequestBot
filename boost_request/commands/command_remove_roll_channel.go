package commands

import (
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var removeRollChannelCommand = dgc.Command{
	Name:        "removerollchannel",
	Description: "Disables posting RNG rolls when an advertiser has been chosen for a boost request.",
	Usage:       "!boostrequest removerollchannel",
	Example:     "!boostrequest removerollchannel",
	IgnoreCase:  true,
	Handler:     removeRollChannelHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func removeRollChannelHandler(ctx *dgc.Ctx) {
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	err := repo.DeleteRollChannel(ctx.Event.GuildID)
	if err != nil {
		log.Printf("Error deleting boost request roll channel: %v", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Boost request roll posts have been disabled.")
}
