package commands

import (
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
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	err := repo.DeleteBoostRequestChannelsInGuild(ctx.Event.GuildID)
	if err != nil {
		log.Println("Error deleting all channels in guild", err)
		respondText(ctx, genericError)
	} else {
		respondText(ctx, "Removed all boost request channels.")
	}
}
