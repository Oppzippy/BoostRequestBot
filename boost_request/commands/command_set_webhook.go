package commands

import (
	"log"
	"net/url"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var setWebhookCommand = dgc.Command{
	Name:        "setwebhook",
	Description: "Sets a webhook url that will receive an http request when events occur.",
	Usage:       "!boostrequest setwebhook <webhookUrl>",
	Example:     "!boostrequest setwebhook https://localhost/boostRequestBotWebhook",
	IgnoreCase:  true,
	Handler:     setWebhookHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func setWebhookHandler(ctx *dgc.Ctx) {
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	webhookURL := ctx.Arguments.Get(0).Raw()
	if webhookURL == "" {
		respondText(ctx, "Usage: "+ctx.Command.Usage)
		return
	}
	u, err := url.ParseRequestURI(webhookURL)
	if err != nil || u.Host == "" || u.Scheme == "" {
		respondText(ctx, "Invalid url.")
		return
	}

	err = repo.InsertWebhook(repository.Webhook{
		GuildID: ctx.Event.GuildID,
		URL:     webhookURL,
	})
	if err != nil {
		log.Printf("Error inserting webhook: %v", err)
		respondText(ctx, genericError)
		return
	}
	respondText(ctx, "Webhook set.")
}
