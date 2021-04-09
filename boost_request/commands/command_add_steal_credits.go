package commands

import (
	"fmt"
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var addStealCreditsCommand = dgc.Command{
	Name:        "addcredits",
	Description: "Adds boost request steal credits to a user. Add a negative amount to subtract credits.",
	Usage:       "!boostrequest addcredits <@user> <credits>",
	Example:     "!boostrequest addcredits @JohnDoe 5",
	IgnoreCase:  true,
	Handler:     addStealCreditsHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func addStealCreditsHandler(ctx *dgc.Ctx) {
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
	err = repo.AdjustStealCreditsForUser(ctx.Event.GuildID, userID, repository.OperationAdd, amount)
	if err != nil {
		log.Printf("Error updating steal credits: %v", err)
		respondText(ctx, genericError)
		return
	}
	newCredits, err := repo.GetStealCreditsForUser(ctx.Event.GuildID, userID)
	if err != nil {
		log.Printf("Error fetching steal credits: %v", err)
		respondText(ctx, fmt.Sprintf("Added %d steal credits.", amount))
	} else {
		respondText(ctx, fmt.Sprintf("Added %d steal credits. New total is %d.", amount, newCredits))
	}
}
