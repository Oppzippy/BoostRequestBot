package commands

import (
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var removeRoleDiscountCommand = dgc.Command{
	Name:        "removerolediscount",
	Description: "Removes a role discount.",
	Usage:       "!boostrequest removerolediscount <@Role>",
	Example:     "!boostrequest removerolediscount @Booster",
	IgnoreCase:  true,
	Handler:     removeRoleDiscountHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func removeRoleDiscountHandler(ctx *dgc.Ctx) {
	if ctx.Arguments.Amount() != 1 {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}
	roleID := ctx.Arguments.Get(1).AsRoleMentionID()
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	rd, err := repo.GetRoleDiscountForRole(ctx.Event.GuildID, roleID)
	if err != nil {
		log.Printf("Error fetching role discount: %v", err)
		respondText(ctx, genericError)
		return
	}
	err = repo.DeleteRoleDiscount(rd)
	if err != nil {
		log.Printf("Error deleting role discount: %v", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Removed role discount.")
}
