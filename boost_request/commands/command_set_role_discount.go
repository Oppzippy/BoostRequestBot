package commands

import (
	"fmt"
	"log"
	"regexp"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/shopspring/decimal"
)

var setRoleDiscountCommand = dgc.Command{
	Name:        "setrolediscount",
	Description: "Sets a permanent discount when anyone with a particular role requests a boost.",
	Usage:       "!boostrequest setrolediscount <@Role> <boostType> <discount%>",
	Example:     "!boostrequest setrolediscount @Booster mythic+ 10%",
	IgnoreCase:  true,
	Handler:     setRoleDiscountHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

var percentRegex = regexp.MustCompile(`([0-9\.]+)%`)

func setRoleDiscountHandler(ctx *dgc.Ctx) {
	if ctx.Arguments.Amount() != 3 {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}
	roleID := ctx.Arguments.Get(0).AsRoleMentionID()
	boostType := ctx.Arguments.Get(1).Raw()
	discountStr := ctx.Arguments.Get(2).Raw()

	discountPercent, err := parsePercent(discountStr)
	if err != nil {
		log.Printf("Failed to parse percentage: %v", err)
		respondText(ctx, "Failed to parse discount percentage.")
		return
	}
	discount := discountPercent.Div(decimal.NewFromInt(100))
	if discount.LessThanOrEqual(decimal.Zero) || discount.GreaterThan(decimal.NewFromInt(1)) {
		respondText(ctx, "Discount must be greater than 0% and less than or equal to 100%.")
		return
	}

	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	err = repo.InsertRoleDiscount(&repository.RoleDiscount{
		GuildID:   ctx.Event.GuildID,
		RoleID:    roleID,
		BoostType: boostType,
		Discount:  discount,
	})
	if err != nil {
		log.Printf("Error inserting role discount: %v", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, fmt.Sprintf("Added %v%% role discount.", discountPercent))
}

func parsePercent(percent string) (decimal.Decimal, error) {
	match := percentRegex.FindStringSubmatch(percent)
	if match != nil {
		return decimal.NewFromString(match[1])
	}
	return decimal.NewFromString(percent)
}
