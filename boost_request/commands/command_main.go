package commands

import "github.com/lus/dgc"

var MainCommand = dgc.Command{
	Name:        "boostrequest",
	Description: "Boost request bot administration commands.",
	Usage:       "!boostrequest <command>",
	Example:     "!boostrequest removechannels",
	Flags:       []string{"ADMIN", "GUILD"},
	IgnoreCase:  true,

	SubCommands: []*dgc.Command{
		&addChannelCommand,
		&removeChannelCommand,
		&removeChannelsCommand,
		&setPrivilegesCommand,
		&removePrivilegesCommand,
	},

	Handler: mainCommandHandler,
}

func mainCommandHandler(ctx *dgc.Ctx) {
}
