package commands

import (
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var removePrivilegesCommand = dgc.Command{
	Name:        "removeprivileges",
	Description: "Removes a role's boost request privileges.",
	Usage:       "!boostrequest removeprivileges <@role>",
	Example:     "!boostrequest removeprivileges @Advertiser",
	IgnoreCase:  true,
	Handler:     removePrivilegesHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func removePrivilegesHandler(ctx *dgc.Ctx) {
	if ctx.Arguments.Amount() != 1 {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}
	roleID := ctx.Arguments.Get(0).AsRoleMentionID()

	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	privileges, err := repo.GetAdvertiserPrivilegesForRole(ctx.Event.GuildID, roleID)
	if err != nil {
		log.Println("Error fetching advertiser privileges for role", err)
		respondText(ctx, genericError)
		return
	}
	if privileges == nil {
		respondText(ctx, "This role has no privileges.")
		return
	}
	err = repo.DeleteAdvertiserPrivileges(privileges)
	if err != nil {
		log.Println("Error deleting advertiser privileges", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Removed privileges.")
}
