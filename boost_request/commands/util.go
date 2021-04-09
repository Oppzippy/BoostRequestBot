package commands

import (
	"log"

	"github.com/lus/dgc"
)

const genericError = "An error has occurred."

func respondText(ctx *dgc.Ctx, message string) {
	_, err := ctx.Session.ChannelMessageSendReply(ctx.Event.ChannelID, message, ctx.Event.Message.Reference())
	if err != nil {
		log.Printf("An error occurred while responding to !%s: %v", ctx.Command.Name, err)
	}
}
