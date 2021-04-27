package middleware

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/lus/dgc"
)

// TODO consider caching user permissions and guild roles
type AdminOnlyMiddleware struct {
}

func NewAdminOnlyMiddleware() *AdminOnlyMiddleware {
	return &AdminOnlyMiddleware{}
}

func (mw *AdminOnlyMiddleware) Exec(next dgc.ExecutionHandler) dgc.ExecutionHandler {
	return func(ctx *dgc.Ctx) {
		var isAdmin bool
		permissions, err := mw.getPermissions(ctx.Session, ctx.Event.GuildID, ctx.Event.Author.ID)
		if err != nil {
			log.Printf("Error fetching permissions: %v", err)
			return
		}
		if permissions&discordgo.PermissionAdministrator != 0 {
			isAdmin = true
		}
		ctx.CustomObjects.Set("isAdmin", isAdmin)
		if isAdmin || !commandHasFlag(ctx.Command, "ADMIN") {
			next(ctx)
		}
	}
}

func (mw *AdminOnlyMiddleware) getPermissions(discord *discordgo.Session, guildID, userID string) (int64, error) {
	guild, err := discord.State.Guild(guildID)
	if err != nil {
		return 0, fmt.Errorf("fetching guild from state cache: %w", err)
	}
	if userID == guild.OwnerID {
		return discordgo.PermissionAll, nil
	}

	member, err := discord.GuildMember(guildID, userID)
	if err != nil {
		return 0, fmt.Errorf("fetching guild member from discord: %w", err)
	}

	allRolesByID := mw.indexRolesByID(guild.Roles)
	permissions := int64(0)
	for _, roleID := range member.Roles {
		role := allRolesByID[roleID]
		permissions |= role.Permissions
	}

	return permissions, nil
}

func (mw *AdminOnlyMiddleware) indexRolesByID(roles []*discordgo.Role) map[string]*discordgo.Role {
	rolesByID := make(map[string]*discordgo.Role)
	for _, role := range roles {
		rolesByID[role.ID] = role
	}
	return rolesByID
}
