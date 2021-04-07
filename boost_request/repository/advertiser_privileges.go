package repository

type AdvertiserPrivileges struct {
	ID      int64
	GuildID string
	RoleID  string
	Weight  float64
	Delay   int
}
