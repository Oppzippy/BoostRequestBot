//go:generate mockgen -source repository.go -destination mock_repository/repository.go
package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNoResults = errors.New("not found")
var ErrInvalidOperation = errors.New("invalid math operation")
var ErrDuplicate = errors.New("the object already exists")

// Repository All methods are thread safe.
// All method that interact with a database may return database errors in addition to the specified ones.
type Repository interface {
	ApiKeyRepository
	BoostRequestChannelRepository
	BoostRequestRepository
	AdvertiserPrivilegesRepository
	LogChannelRepository
	StealCreditRepository
	RollChannelRepository
	WebhookRepository
	AutoSignupSessionRepository
	DelayedMessageRepository
}

type ApiKeyRepository interface {
	// GetAPIKey Returns either the API key or ErrNoResults.
	GetAPIKey(key string) (*APIKey, error)
	// NewAPIKey Creates a new api key for the specified guild and returns it.
	NewAPIKey(guildID string) (*APIKey, error)
}

type BoostRequestRepository interface {
	// GetUnresolvedBoostRequests Returns all unresolved boost requests for any guild.
	GetUnresolvedBoostRequests() ([]*BoostRequest, error)
	// GetBoostRequestByBackendMessageID Returns the boost request associated with a backend message or ErrNoResults.
	GetBoostRequestByBackendMessageID(backendChannelID, backendMessageID string) (*BoostRequest, error)
	// GetBoostRequestById Returns a boost request in the specified guild with the specified id
	GetBoostRequestById(guildID string, boostRequestID uuid.UUID) (*BoostRequest, error)
	// InsertBoostRequest Stores a boost request. Boost requests are unique by backend message id, so this will error if a boost request already exists.
	// br will have its ID field updated to match the newly inserted row'd id.
	InsertBoostRequest(br *BoostRequest) error
	// ResolveBoostRequest Updates fields necessary to flag a boost request as resolved.
	ResolveBoostRequest(br *BoostRequest) error
	// DeleteBoostRequest Deletes a boost request
	DeleteBoostRequest(br *BoostRequest) error
	InsertBoostRequestDelayedMessage(br *BoostRequest, delayedMessage *DelayedMessage) error
	GetBoostRequestDelayedMessageIDs(br *BoostRequest) ([]int64, error)
}

type BoostRequestChannelRepository interface {
	// GetBoostRequestChannelByFrontendChannelID Returns the boost request channel with the specified frontend channel. Frontend channels are unique.
	// Returns ErrNoResults if the channel does not exist.
	GetBoostRequestChannelByFrontendChannelID(guildID string, frontendChannelID string) (*BoostRequestChannel, error)
	// GetBoostRequestChannels Returns all boost request channels in a guild.
	GetBoostRequestChannels(guildID string) ([]*BoostRequestChannel, error)
	// InsertBoostRequestChannel Creates a new boost request channel. If the channel already exists, it will be updated.
	// brc will have its ID field updated to match the newly inserted row's id.
	InsertBoostRequestChannel(brc *BoostRequestChannel) error
	// DeleteBoostRequestChannel Deletes a specific boost request channel.
	DeleteBoostRequestChannel(brc *BoostRequestChannel) error
	// DeleteBoostRequestChannelsInGuild Deletes all boost request channels in a guild.
	DeleteBoostRequestChannelsInGuild(guildID string) error
}

type AdvertiserPrivilegesRepository interface {
	// GetAdvertiserPrivilegesForGuild Returns advertiser privileges for all roles in a guild.
	GetAdvertiserPrivilegesForGuild(guildID string) ([]*AdvertiserPrivileges, error)
	// GetAdvertiserPrivilegesForRole Returns advertiser privileges for a particular role in a guild or ErrNoResults.
	GetAdvertiserPrivilegesForRole(guildID, roleID string) (*AdvertiserPrivileges, error)
	// InsertAdvertiserPrivileges Creates advertiser privileges or updates privileges if the role already has privileges set.
	InsertAdvertiserPrivileges(privileges *AdvertiserPrivileges) error
	// DeleteAdvertiserPrivileges Deletes specific advertiser privileges.
	DeleteAdvertiserPrivileges(privileges *AdvertiserPrivileges) error
}

type LogChannelRepository interface {
	// GetLogChannel Returns a guild's log channel or ErrNoResults.
	GetLogChannel(guildID string) (channelID string, err error)
	// InsertLogChannel Creates a log channel for a guild or updates it if it already exists.
	InsertLogChannel(guildID, channelID string) error
	// DeleteLogChannel Deletes a guild's log channel.
	DeleteLogChannel(guildID string) error
}

type StealCreditRepository interface {
	// GetStealCreditsForUser Returns the number of boost request steal credits avilable to a user.
	GetStealCreditsForUser(guildID, userID string) (int, error)
	// GetGlobalStealCreditsForUser Returns the number of boost request steal credits available to a user across all guilds.
	GetGlobalStealCreditsForUser(userID string) (map[string]int, error)
	// AdjustStealCreditsForUser Performs a math operation on the number of boost request steal credits in a user's possession.
	// If a valid operation is not specified, ErrInvalidOperation will be returned.
	AdjustStealCreditsForUser(guildID, userID string, operation Operation, amount int) error
	// UpdateStealCreditsForUser Short version of AdjustStealCreditsForUser with OperationSet
	UpdateStealCreditsForUser(guildID, userID string, amount int) error
}

type RollChannelRepository interface {
	// GetRollChannel Returns the channel ID that boost request RNG rolls should be posted to
	// or ErrNoResults if rolls should not be posted
	GetRollChannel(guildID string) (channelID string, err error)
	// InsertRollChannel Sets the channel that boost request RNG rolls should be posted to
	InsertRollChannel(guildID, channelID string) error
	// DeleteRollChannel Stops posting boost request RNG rolls
	DeleteRollChannel(guildID string) error
}

type WebhookRepository interface {
	GetWebhook(guildId string) (Webhook, error)
	InsertWebhook(webhook Webhook) error
	DeleteWebhook(Webhook Webhook) error
	InsertQueuedWebhook(webhook Webhook, body string) error
	GetQueuedWebhooks() ([]*QueuedWebhookRequest, error)
	InsertWebhookAttempt(attempt WebhookAttempt) error
}

type AutoSignupSessionRepository interface {
	IsAutoSignupEnabled(guildID, advertiserID string) (bool, error)
	EnableAutoSignup(guildID, advertiserID string, expiresAt time.Time) (*AutoSignupSession, error)
	CancelAutoSignup(guildID, advertiserID string) error
	GetEnabledAutoSignups() ([]*AutoSignupSession, error)
	GetEnabledAutoSignupsInGuild(guildID string) ([]*AutoSignupSession, error)
	InsertAutoSignupDelayedMessages(autoSignup *AutoSignupSession, delayedMessages []*DelayedMessage) error
	GetAutoSignupDelayedMessageIDs(guildID, advertiserID string) ([]int64, error)
}

type DelayedMessageRepository interface {
	GetDelayedMessages() ([]*DelayedMessage, error)
	InsertDelayedMessage(delayedMessage *DelayedMessage) error
	DeleteDelayedMessage(id int64) error
	FlagDelayedMessageAsSent(message *DelayedMessage) error
}
