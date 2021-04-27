package commands

import (
	"fmt"
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var removeChannelsCommand = dgc.Command{
	Name:        "removechannels",
	Description: "Stops watching (unregisters) all boost request channels.",
	Usage:       "!boostrequest removechannels",
	Example:     "!boostrequest removechannels",
	IgnoreCase:  true,
	Handler:     removeChannelsHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func removeChannelsHandler(ctx *dgc.Ctx) {
	if ctx.Arguments.Amount() != 0 {
		respondText(ctx, fmt.Sprintf("Usage: %v. Did you mean !boostrequest removechannel?", ctx.Command.Usage))
		return
	}
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	err := repo.DeleteBoostRequestChannelsInGuild(ctx.Event.GuildID)
	if err != nil {
		log.Printf("Error deleting all channels in guild: %v", err)
		respondText(ctx, genericError)
	} else {
		respondText(ctx, "Removed all boost request channels.")
	}
}
