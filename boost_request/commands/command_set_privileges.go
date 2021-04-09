package commands

import (
	"log"
	"strconv"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var setPrivilegesCommand = dgc.Command{
	Name:        "setprivileges",
	Description: "Sets the weight and wait time for an advertiser role.",
	Usage:       "!boostrequest setprivileges <@role> <weight> <delayInSeconds>",
	Example:     "!boostrequest setprivileges @Advertiser 1.0 60",
	IgnoreCase:  true,
	Handler:     setPrivilegesHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func setPrivilegesHandler(ctx *dgc.Ctx) {
	if ctx.Arguments.Amount() != 3 {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}
	roleID := ctx.Arguments.Get(0).AsRoleMentionID()
	weightRaw := ctx.Arguments.Get(1).Raw()
	weight, weightErr := strconv.ParseFloat(weightRaw, 64)
	delaySeconds, delayErr := ctx.Arguments.Get(2).AsInt()
	if roleID == "" || weightErr != nil || delayErr != nil {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}

	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	err := repo.InsertAdvertiserPrivileges(&repository.AdvertiserPrivileges{
		GuildID: ctx.Event.GuildID,
		RoleID:  roleID,
		Weight:  weight,
		Delay:   delaySeconds,
	})
	if err != nil {
		log.Println("Error setting privileges", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Privileges set.")
}
