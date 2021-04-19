package commands

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/shopspring/decimal"
)

var listRoleDiscountsCommand = dgc.Command{
	Name:        "rolediscounts",
	Description: "Lists all role discounts.",
	Usage:       "!boostrequest rolediscounts",
	Example:     "!boostrequest rolediscounts",
	IgnoreCase:  true,
	Handler:     listRoleDiscountsHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func listRoleDiscountsHandler(ctx *dgc.Ctx) {
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	discounts, err := repo.GetRoleDiscountsForGuild(ctx.Event.GuildID)
	if err != nil {
		log.Printf("Error fetching role discounts: %v", err)
		respondText(ctx, genericError)
		return
	}
	sb := strings.Builder{}
	for _, rd := range discounts {
		discountPercent := rd.Discount.Mul(decimal.NewFromInt(100))
		sb.WriteString("<@&")
		sb.WriteString(rd.RoleID)
		sb.WriteString(">: ")
		sb.WriteString(discountPercent.String())
		sb.WriteString("%\n")
	}
	if sb.Len() > 0 {
		ctx.Session.ChannelMessageSendComplex(ctx.Event.ChannelID, &discordgo.MessageSend{
			Reference: ctx.Event.Message.Reference(),
			Content:   sb.String(),
			// TODO add replied_user true when discordgo supports it
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		})
	} else {
		respondText(ctx, "There are no role discounts.")
	}
}
