package context_key

type ContextKey int

const (
	// repository.Repository
	Repository ContextKey = iota
	// string
	GuildID
	// bool
	IsAuthorized
)
