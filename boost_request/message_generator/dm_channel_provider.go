package message_generator

type dmChannelProvider interface {
	DMChannel(userID string) (channelID string, err error)
}
