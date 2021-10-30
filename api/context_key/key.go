package context_key

type ContextKey int

const (
	// repository.Repository
	Repository ContextKey = iota
	// boost_request.BoostRequestManager
	BooostRequestManager
	// string
	GuildID
	// bool
	IsAuthorized
)
