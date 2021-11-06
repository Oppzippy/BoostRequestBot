package commands

import (
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var removeWebhookCommand = dgc.Command{
	Name:        "removewebhook",
	Description: "Stop sending events to the webhook url.",
	Usage:       "!boostrequest removewebhook",
	Example:     "!boostrequest removewebhook",
	IgnoreCase:  true,
	Handler:     removeWebhookHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func removeWebhookHandler(ctx *dgc.Ctx) {
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	webhook, err := repo.GetWebhook(ctx.Event.GuildID)
	if err != nil {
		respondText(ctx, "There is no webhook to remove.")
		return
	}
	err = repo.DeleteWebhook(webhook)
	if err != nil {
		log.Printf("Error deleting webhook: %v", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Webhook removed.")
}
