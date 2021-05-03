package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var listChannelsCommand = dgc.Command{
	Name:        "channels",
	Description: "Lists all boost request channels.",
	Usage:       "!boostrequest channels",
	Example:     "!boostrequest channels",
	IgnoreCase:  true,
	Handler:     listChannelsHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func listChannelsHandler(ctx *dgc.Ctx) {
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	channels, err := repo.GetBoostRequestChannels(ctx.Event.GuildID)
	if err != nil {
		log.Printf("Failed to fetch list of boost request channels: %v", err)
		respondText(ctx, genericError)
		return
	}
	sb := strings.Builder{}
	for i, brc := range channels {
		options := make([]string, 0, 2)
		if brc.SkipsBuyerDM {
			options = append(options, "doesn't dm buyer")
		}
		if brc.UsesBuyerMessage {
			options = append(options, "reacts directly to buyer's message")
		}
		if len(options) == 0 {
			options = append(options, "none")
		}

		sb.WriteString("**Channel ")
		sb.WriteString(fmt.Sprintf("%d", i+1))
		sb.WriteString("**\nFrontend Channel: <#")
		sb.WriteString(brc.FrontendChannelID)
		sb.WriteString(">\nBackend Channel: <#")
		sb.WriteString(brc.BackendChannelID)
		sb.WriteString(">\nOptions: ")
		sb.WriteString(strings.Join(options, ", "))
		sb.WriteString("\n")
	}

	if sb.Len() > 0 {
		respondText(ctx, sb.String())
	} else {
		respondText(ctx, "There are no boost request channels.")
	}
}
