package boost_request

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/boost_request_manager"
	"github.com/oppzippy/BoostRequestBot/boost_request/command_handlers"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/oppzippy/BoostRequestBot/framework/slash_commands"
)

func registerSlashCommandHandlers(
	scr *slash_commands.SlashCommandRegistry,
	bundle *i18n.Bundle,
	repo repository.Repository,
	discord *discordgo.Session,
	brm *boost_request_manager.BoostRequestManager,
) {
	// User commands
	scr.RegisterCommand([]string{"boostrequest", "autosignup", "start"}, command_handlers.NewAutoSignupEnableHandler(bundle, repo, brm).Handle)
	scr.RegisterCommand([]string{"boostrequest", "autosignup", "stop"}, command_handlers.NewAutoSignupDisableHandler(bundle, repo, brm).Handle)

	// Admin commands
	scr.RegisterCommand([]string{"channels", "add"}, command_handlers.NewChannelsAddHandler(bundle, repo, discord).Handle)
	scr.RegisterCommand([]string{"channels", "list"}, command_handlers.NewChannelsListHandler(bundle, repo).Handle)
	scr.RegisterCommand([]string{"channels", "remove"}, command_handlers.NewChannelsRemoveHandler(bundle, repo).Handle)

	scr.RegisterCommand([]string{"credits", "add"}, command_handlers.NewCreditsAddHandler(bundle, repo).Handle)
	scr.RegisterCommand([]string{"credits", "list"}, command_handlers.NewCreditsListHandler(bundle, repo, discord).Handle)
	scr.RegisterCommand([]string{"credits", "set"}, command_handlers.NewCreditsSetHandler(bundle, repo).Handle)

	scr.RegisterCommand([]string{"logchannel", "set"}, command_handlers.NewLogChannelSetHandler(bundle, repo, discord).Handle)
	scr.RegisterCommand([]string{"logchannel", "remove"}, command_handlers.NewLogChannelRemoveHandler(bundle, repo).Handle)

	scr.RegisterCommand([]string{"privileges", "list"}, command_handlers.NewPrivilegesListHandler(bundle, repo).Handle)
	scr.RegisterCommand([]string{"privileges", "remove"}, command_handlers.NewPrivilegesRemoveHandler(bundle, repo).Handle)
	scr.RegisterCommand([]string{"privileges", "set"}, command_handlers.NewPrivilegesSetHandler(bundle, repo).Handle)

	scr.RegisterCommand([]string{"rollchannel", "remove"}, command_handlers.NewRollChannelRemoveHandler(bundle, repo).Handle)
	scr.RegisterCommand([]string{"rollchannel", "set"}, command_handlers.NewRollChannelSetHandler(bundle, repo, discord).Handle)

	scr.RegisterCommand([]string{"webhook", "list"}, command_handlers.NewWebhookListHandler(bundle, repo).Handle)
	scr.RegisterCommand([]string{"webhook", "remove"}, command_handlers.NewWebhookRemoveHandler(bundle, repo).Handle)
	scr.RegisterCommand([]string{"webhook", "set"}, command_handlers.NewWebhookSetHandler(bundle, repo).Handle)
}
