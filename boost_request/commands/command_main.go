package commands

import "github.com/lus/dgc"

var MainCommand = dgc.Command{
	Name:        "boostrequest",
	Description: "Boost request bot administration commands.",
	Usage:       "!boostrequest <command>",
	Example:     "!boostrequest removechannels",
	IgnoreCase:  true,

	SubCommands: []*dgc.Command{
		&addChannelCommand,
		&removeChannelCommand,
		&removeChannelsCommand,
		&setPrivilegesCommand,
		&removePrivilegesCommand,
		&setLogChannelCommand,
		&removeLogChannelCommand,
		&addStealCreditsCommand,
		&setStealCreditsCommand,
		&checkStealCreditsCommand,
	},

	Handler: mainCommandHandler,
}

func mainCommandHandler(ctx *dgc.Ctx) {
}
