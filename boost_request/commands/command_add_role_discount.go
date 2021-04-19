package commands

import (
	"log"
	"regexp"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/shopspring/decimal"
)

var setRoleDiscountCommand = dgc.Command{
	Name:        "setrolediscount",
	Description: "Sets a permanent discount when anyone with a particular role requests a boost.",
	Usage:       "!boostrequest setrolediscount <@Role> <discount%>",
	Example:     "!boostrequest setrolediscount @Booster 10%",
	IgnoreCase:  true,
	Handler:     setRoleDiscountHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

var percentRegex = regexp.MustCompile(`([0-9\.]+)%`)

func setRoleDiscountHandler(ctx *dgc.Ctx) {
	if ctx.Arguments.Amount() != 2 {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}
	roleID := ctx.Arguments.Get(0).AsRoleMentionID()
	discountStr := ctx.Arguments.Get(1).Raw()
	discountPercent, err := parsePercent(discountStr)
	if err != nil {
		log.Printf("Failed to parse percentage: %v", err)
		respondText(ctx, "Failed to parse discount percentage.")
		return
	}
	discount := discountPercent.Div(decimal.NewFromInt(100))

	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	err = repo.InsertRoleDiscount(&repository.RoleDiscount{
		GuildID:  ctx.Event.GuildID,
		RoleID:   roleID,
		Discount: discount,
	})
	if err != nil {
		log.Printf("Error inserting role discount: %v", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Added role discount.")
}

func parsePercent(percent string) (decimal.Decimal, error) {
	match := percentRegex.FindStringSubmatch(percent)
	if match != nil {
		return decimal.NewFromString(match[1])
	}
	return decimal.NewFromString(percent)
}
