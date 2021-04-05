package repository

type BoostRequestChannel struct {
	ID                int64
	GuildID           string
	FrontendChannelID string
	BackendChannelID  string
	UsesBuyerMessage  bool
	NotifiesBuyer     bool
}
