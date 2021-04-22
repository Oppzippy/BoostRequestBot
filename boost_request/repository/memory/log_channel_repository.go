package mock

import "github.com/oppzippy/BoostRequestBot/boost_request/repository"

func (repo *memoryRepository) GetLogChannel(guildID string) (channelID string, err error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	id, ok := repo.logChannels[guildID]
	if ok {
		return id, nil
	} else {
		return "", repository.ErrNoResults
	}
}

func (repo *memoryRepository) InsertLogChannel(guildID, channelID string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	repo.logChannels[guildID] = channelID
	return nil
}

func (repo *memoryRepository) DeleteLogChannel(guildID string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	delete(repo.logChannels, guildID)
	return nil
}
