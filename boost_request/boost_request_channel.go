package boost_request

type BoostRequestChannel struct {
	ID                int
	GuildID           string
	FrontendChannelID string
	BackendChannelID  string
	UsesBuyerMessage  bool
	NotifiesBuyer     bool
}
