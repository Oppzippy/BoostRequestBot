package repository

import (
	"errors"
)

var ErrNoResults = errors.New("not found")
var ErrInvalidOperation = errors.New("invalid math operation")
var ErrDuplicate = errors.New("the object already exists")

// All methods are thread safe.
// All method that interact with a database may return database errors in addition to the specified ones.
type Repository interface {
	ApiKeyRepository
	BoostRequestChannelRepository
	BoostRequestRepository
	AdvertiserPrivilegesRepository
	LogChannelRepository
	RoleDiscountRepository
	StealCreditRepository
}

type ApiKeyRepository interface {
	// Returns either the API key or ErrNoResults.
	GetAPIKey(key string) (*APIKey, error)
}

type BoostRequestRepository interface {
	// Returns all unresolved boost requests for any guild.
	GetUnresolvedBoostRequests() ([]*BoostRequest, error)
	// Returns the boost request associated with a backend message or ErrNoResults.
	GetBoostRequestByBackendMessageID(backendChannelID, backendMessageID string) (*BoostRequest, error)
	// Stores a boost request. Boost requests are unique by backend message id, so this will error if a boost request already exists.
	// br will have its ID field updated to match the newly inserted row'd id.
	InsertBoostRequest(br *BoostRequest) error
	// Updates fields necessary to flag a boost request as resolved.
	ResolveBoostRequest(br *BoostRequest) error
}

type BoostRequestChannelRepository interface {
	// Returns the boost request channel with the specified frontend channel. Frontend channels are unique.
	// Returns ErrNoResults if the channel does not exist.
	GetBoostRequestChannelByFrontendChannelID(guildID string, frontendChannelID string) (*BoostRequestChannel, error)
	// Returns all boost request channels in a guild.
	GetBoostRequestChannels(guildID string) ([]*BoostRequestChannel, error)
	// Creates a new boost request channel. If the channel already exists, it will be updated.
	// brc will have its ID field updated to match the newly inserted row's id.
	InsertBoostRequestChannel(brc *BoostRequestChannel) error
	// Deletes a specific boost request channel.
	DeleteBoostRequestChannel(brc *BoostRequestChannel) error
	// Deletes all boost request channels in a guild.
	DeleteBoostRequestChannelsInGuild(guildID string) error
}

type AdvertiserPrivilegesRepository interface {
	// Returns advertiser privileges for all roles in a guild.
	GetAdvertiserPrivilegesForGuild(guildID string) ([]*AdvertiserPrivileges, error)
	// Returns advertiser privileges for a particular role in a guild or ErrNoResults.
	GetAdvertiserPrivilegesForRole(guildID, roleID string) (*AdvertiserPrivileges, error)
	// Creates advertiser privileges or updates privileges if the role already has privileges set.
	InsertAdvertiserPrivileges(privileges *AdvertiserPrivileges) error
	// Deletes specific advertiser privileges.
	DeleteAdvertiserPrivileges(privileges *AdvertiserPrivileges) error
}

type LogChannelRepository interface {
	// Returns a guild's log channel or ErrNoResults.
	GetLogChannel(guildID string) (channelID string, err error)
	// Creates a log channel for a guild or updates it if it already exists.
	InsertLogChannel(guildID, channelID string) error
	// Deletes a guild's log channel.
	DeleteLogChannel(guildID string) error
}

type RoleDiscountRepository interface {
	// Returns the best discount of each boost type available to the provided roles.
	GetBestDiscountsForRoles(guildID string, roleID []string) ([]*RoleDiscount, error)
	// Returns a role's discounts on all types of boosts.
	GetRoleDiscountsForRole(guildID, roleID string) ([]*RoleDiscount, error)
	// Returns a role's discount for a specific boost type or ErrNoResults.
	GetRoleDiscountForBoostType(guildID, roleID, boostType string) (*RoleDiscount, error)
	// Returns all role discounts in a guild.
	GetRoleDiscountsForGuild(guildID string) ([]*RoleDiscount, error)
	// Creates a role discount in a guild. If the discount already exists, it will be updated.
	InsertRoleDiscount(rd *RoleDiscount) error
	// Deletes a role discount in a guild.
	DeleteRoleDiscount(rd *RoleDiscount) error
}

type StealCreditRepository interface {
	// Returns the number of boost request steal credits avilable to a user.
	GetStealCreditsForUser(guildID, userID string) (int, error)
	// Returns the number of boost request steal credits available to a user across all guilds.
	GetGlobalStealCreditsForUser(userID string) (map[string]int, error)
	// Performs a math operation on the number of boost request steal credits in a user's possession.
	// If a valid operation is not specified, ErrInvalidOperation will be returned.
	AdjustStealCreditsForUser(guildID, userID string, operation Operation, amount int) error
	// Short version of AdjustStealCreditsForUser with OperationSet
	UpdateStealCreditsForUser(guildID, userID string, amount int) error
}
