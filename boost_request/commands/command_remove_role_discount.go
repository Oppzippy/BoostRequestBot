package commands

import (
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var removeRoleDiscountCommand = dgc.Command{
	Name:        "removerolediscount",
	Description: "Removes a role discount.",
	Usage:       "!boostrequest removerolediscount <@Role> <boostType>",
	Example:     "!boostrequest removerolediscount @Booster mythic+",
	IgnoreCase:  true,
	Handler:     removeRoleDiscountHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func removeRoleDiscountHandler(ctx *dgc.Ctx) {
	if ctx.Arguments.Amount() != 2 {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}
	roleID := ctx.Arguments.Get(0).AsRoleMentionID()
	boostType := ctx.Arguments.Get(1).Raw()

	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	rd, err := repo.GetRoleDiscountForBoostType(ctx.Event.GuildID, roleID, boostType)
	if err == repository.ErrNoResults {
		respondText(ctx, "This role does not have a discount for the specified boost request type.")
		return
	}
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
