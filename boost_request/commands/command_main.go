package commands

import "github.com/lus/dgc"

var MainCommand = dgc.Command{
	Name:        "boostrequest",
	Description: "Boost request bot administration commands.",
	Usage:       "!boostrequest <command>",
	Example:     "!boostrequest removechannels",
	IgnoreCase:  true,

	SubCommands: []*dgc.Command{
		// Boost request channels
		&addChannelCommand,
		&listChannelsCommand,
		&removeChannelCommand,
		&removeChannelsCommand,
		// Log channel
		&setLogChannelCommand,
		&removeLogChannelCommand,
		// Advertiser privileges
		&setPrivilegesCommand,
		&listPrivilegesCommand,
		&removePrivilegesCommand,
		// Steal credits
		&addStealCreditsCommand,
		&setStealCreditsCommand,
		&checkStealCreditsCommand,
		// Role discounts
		&setRoleDiscountCommand,
		&listRoleDiscountsCommand,
		&removeRoleDiscountCommand,
		// Boost request RNG Rolls
		&setRollChannelCommand,
		&removeRollChannelCommand,
	},

	Handler: mainCommandHandler,
}

func mainCommandHandler(ctx *dgc.Ctx) {
}
