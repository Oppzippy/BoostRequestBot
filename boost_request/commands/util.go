package commands

import (
	"log"

	"github.com/lus/dgc"
)

const genericError = "An error has occurred."

func respondText(ctx *dgc.Ctx, message string) {
	err := ctx.RespondText(message)
	if err != nil {
		log.Printf("An error occurred while responding to !%s: %s", ctx.Command.Name, message)
	}
}
