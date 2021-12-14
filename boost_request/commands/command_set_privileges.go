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
	Usage:       "!boostrequest setprivileges <@role> <weight> <delayInSeconds> [autoSignupDurationInMinutes]",
	Example:     "!boostrequest setprivileges @Advertiser 1.0 60",
	IgnoreCase:  true,
	Handler:     setPrivilegesHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func setPrivilegesHandler(ctx *dgc.Ctx) {
	if ctx.Arguments.Amount() < 3 {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}
	roleID := ctx.Arguments.Get(0).AsRoleMentionID()
	weightRaw := ctx.Arguments.Get(1).Raw()
	weight, weightErr := strconv.ParseFloat(weightRaw, 64)
	delaySeconds, delayErr := ctx.Arguments.Get(2).AsInt()
	var autoSignupDurationMinutes int
	var autoSignupDurationErr error
	if ctx.Arguments.Amount() >= 4 {
		autoSignupDurationMinutes, autoSignupDurationErr = ctx.Arguments.Get(3).AsInt()
	}
	if roleID == "" || weightErr != nil || delayErr != nil || autoSignupDurationErr != nil {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}

	if weight <= 0 {
		respondText(ctx, "Weight must be greater than 0.")
		return
	}
	if delaySeconds < 0 {
		respondText(ctx, "Delay must be greater than or equal to 0.")
		return
	}
	if autoSignupDurationMinutes < 0 {
		respondText(ctx, "Auto signup duration must be greater than or equal to 0.")
		return
	}

	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	err := repo.InsertAdvertiserPrivileges(&repository.AdvertiserPrivileges{
		GuildID:            ctx.Event.GuildID,
		RoleID:             roleID,
		Weight:             weight,
		Delay:              delaySeconds,
		AutoSignupDuration: autoSignupDurationMinutes * 60,
	})
	if err != nil {
		log.Printf("Error setting privileges: %v", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Privileges set.")
}
