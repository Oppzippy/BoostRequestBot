package commands

import (
	"fmt"
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var setStealCreditsCommand = dgc.Command{
	Name:        "setcredits",
	Description: "Sets the number of boost request steal credits available to a user.",
	Usage:       "!boostrequest setcredits <@user> <credits>",
	Example:     "!boostrequest setcredits @JohnDoe 5",
	IgnoreCase:  true,
	Handler:     setStealCreditsHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func setStealCreditsHandler(ctx *dgc.Ctx) {
	if ctx.Arguments.Amount() != 2 {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}
	userID := ctx.Arguments.Get(0).AsUserMentionID()
	amount, err := ctx.Arguments.Get(1).AsInt()
	if userID == "" || err != nil {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}

	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	err = repo.UpdateStealCreditsForUser(ctx.Event.GuildID, userID, amount)
	if err != nil {
		log.Printf("Error updating steal credits: %v", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, fmt.Sprintf("Set steal credits to %d.", amount))
}
