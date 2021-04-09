package commands

import (
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var removeLogChannelCommand = dgc.Command{
	Name:        "removelogchannel",
	Description: "Stops boost request logging.",
	Usage:       "!boostrequest removelogchannel",
	Example:     "!boostrequest removelogchannel",
	IgnoreCase:  true,
	Handler:     removeLogChannelHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func removeLogChannelHandler(ctx *dgc.Ctx) {
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	err := repo.DeleteLogChannel(ctx.Event.GuildID)
	if err != nil {
		log.Println("Error deleting boost request log channel", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Boost request logging is disabled.")
}
