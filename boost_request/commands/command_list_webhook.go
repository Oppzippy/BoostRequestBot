package commands

import (
	"log"

	"github.com/lus/dgc"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

var listWebhookCommand = dgc.Command{
	Name:        "webhook",
	Description: "Lists the webhook",
	Usage:       "!boostrequest webhook",
	Example:     "!boostrequest webhook",
	IgnoreCase:  true,
	Handler:     listWebhookHandler,
	Flags:       []string{"ADMIN", "GUILD"},
}

func listWebhookHandler(ctx *dgc.Ctx) {
	repo := ctx.CustomObjects.MustGet("repo").(repository.Repository)
	webhook, err := repo.GetWebhook(ctx.Event.GuildID)
	if err == repository.ErrNoResults {
		respondText(ctx, "A webhook is not set.")
		return
	}
	if err != nil {
		log.Printf("Error fetching webhook: %v", err)
		respondText(ctx, genericError)
		return
	}
	channel, err := ctx.Session.UserChannelCreate(ctx.Event.Message.Author.ID)
	if err != nil {
		log.Printf("Error creating dm channel: %v", err)
		respondText(ctx, "Failed to send DM.")
		return
	}

	_, err = ctx.Session.ChannelMessageSend(channel.ID, "Webhook: "+webhook.URL)
	if err != nil {
		log.Printf("Error sending dm: %v", err)
		respondText(ctx, "Failed to send DM.")
		return
	}

	respondText(ctx, "Sent webhook url in a DM.")
}
