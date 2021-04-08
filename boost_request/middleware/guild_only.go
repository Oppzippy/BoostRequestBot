package middleware

import "github.com/lus/dgc"

func GuildOnlyMiddleware(next dgc.ExecutionHandler) dgc.ExecutionHandler {
	return func(ctx *dgc.Ctx) {
		if commandHasFlag(ctx.Command, "GUILD") && ctx.Event.GuildID == "" {
			return
		}
		next(ctx)
	}
}
