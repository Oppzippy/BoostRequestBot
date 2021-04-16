package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
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
	if ctx.Event.GuildID != "" {
		checkStealCreditsHandlerGuild(ctx)
	}
}

func checkStealCreditsHandlerGuild(ctx *dgc.Ctx) {
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)

	if ctx.Arguments.Amount() == 0 {
		credits, err := repo.GetStealCreditsForUser(ctx.Event.GuildID, ctx.Event.Author.ID)
		if err != nil {
			log.Printf("Error fetching boost request steal credits: %v", err)
			respondText(ctx, genericError)
			return
		}
		respondText(ctx, fmt.Sprintf("You have %d boost request steal credits.", credits))
	} else {
		isAdmin, ok := ctx.CustomObjects.Get("isAdmin")
		if ok && isAdmin.(bool) {
			checkStealCreditsForUsers(ctx)
		}
	}
}

func checkStealCreditsForUsers(ctx *dgc.Ctx) {
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)

	userIDs := make([]string, 0, 1)
	for i := 0; i < ctx.Arguments.Amount(); i++ {
		userIDs = append(userIDs, ctx.Arguments.Get(i).AsUserMentionID())
	}

	sb := strings.Builder{}
	for _, userID := range userIDs {
		credits, err := repo.GetStealCreditsForUser(ctx.Event.GuildID, userID)
		if err != nil {
			log.Printf("Error fetching boost request steal credits for user in admin check credits command: %v", err)
			respondText(ctx, genericError)
			return
		}
		// Users won't get pinged since AllowedMentions is empty
		sb.WriteString("<@")
		sb.WriteString(userID)
		sb.WriteString(">: ")
		sb.WriteString(fmt.Sprintf("%d", credits))
		sb.WriteString(" credits\n")
	}

	ctx.Session.ChannelMessageSendComplex(ctx.Event.ChannelID, &discordgo.MessageSend{
		Reference: ctx.Event.Message.Reference(),
		Content:   sb.String(),
		// TODO add replied_user true when discordgo supports it
		AllowedMentions: &discordgo.MessageAllowedMentions{},
	})
}
