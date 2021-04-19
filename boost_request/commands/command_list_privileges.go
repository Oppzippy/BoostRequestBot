package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var listPrivilegesCommand = dgc.Command{
	Name:        "privileges",
	Description: "Lists all privileges.",
	Usage:       "!boostrequest privileges",
	Example:     "!boostrequest privileges",
	IgnoreCase:  true,
	Handler:     listPrivilegesHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func listPrivilegesHandler(ctx *dgc.Ctx) {
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	allPrivileges, err := repo.GetAdvertiserPrivilegesForGuild(ctx.Event.GuildID)
	if err != nil {
		log.Printf("Error listing all privileges: %v", err)
		respondText(ctx, genericError)
		return
	}

	sb := strings.Builder{}
	for _, p := range allPrivileges {
		// We won't ping the roles since AllowedMentions is empty
		sb.WriteString("<@&" + p.RoleID + ">")
		sb.WriteString(" Weight: ")
		sb.WriteString(strconv.FormatFloat(p.Weight, 'f', -1, 64))
		sb.WriteString(", Delay: ")
		sb.WriteString(fmt.Sprintf("%d", p.Delay))
		sb.WriteString("s\n")
	}
	if sb.Len() > 0 {
		_, err = ctx.Session.ChannelMessageSendComplex(ctx.Event.ChannelID, &discordgo.MessageSend{
			Reference: ctx.Event.Message.Reference(),
			Content:   sb.String(),
			// TODO add replied_user true when discordgo supports it
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		})
		if err != nil {
			log.Printf("Error sending message with list of privileges: %v", err)
		}
	} else {
		respondText(ctx, "No privileges are set.")
	}
}
