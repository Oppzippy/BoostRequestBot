package commands

import (
	"fmt"
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var checkStealCreditsCommand = dgc.Command{
	Name:        "credits",
	Aliases:     []string{"credit", "stealcredits", "steals"},
	Description: "Checks the number of boost credits that you have available.",
	Usage:       "!boostrequest credits",
	Example:     "!boostrequest credits",
	IgnoreCase:  true,
	Handler:     checkStealCreditsHandler,
	Flags:       []string{"GUILD"},
}

func checkStealCreditsHandler(ctx *dgc.Ctx) {
	// TODO implement DM functionality
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	if ctx.Event.GuildID != "" {
		credits, err := repo.GetStealCreditsForUser(ctx.Event.GuildID, ctx.Event.Author.ID)
		if err != nil {
			log.Printf("Error fetching boost request steal credits: %v", err)
			respondText(ctx, genericError)
			return
		}
		respondText(ctx, fmt.Sprintf("You have %d boost request steal credits.", credits))
	}
}
