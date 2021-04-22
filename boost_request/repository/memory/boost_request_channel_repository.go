package mock

import "github.com/oppzippy/BoostRequestBot/boost_request/repository"

func (repo *memoryRepository) GetBoostRequestChannelByFrontendChannelID(guildID string, frontendChannelID string) (*repository.BoostRequestChannel, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for _, brc := range repo.channels {
		if brc.GuildID == guildID && brc.FrontendChannelID == frontendChannelID {
			return brc, nil
		}
	}
	return nil, repository.ErrNoResults
}

func (repo *memoryRepository) InsertBoostRequestChannel(brc *repository.BoostRequestChannel) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	brc.ID = repo.lastID
	repo.lastID += 1
	repo.channels = append(repo.channels, brc)
	return nil
}

func (repo *memoryRepository) DeleteBoostRequestChannel(brc *repository.BoostRequestChannel) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for i, c := range repo.channels {
		if c.ID == brc.ID {
			repo.channels[i] = repo.channels[len(repo.channels)-1]
			repo.channels = repo.channels[:len(repo.channels)-1]
			break
		}
	}
	return nil
}

func (repo *memoryRepository) DeleteBoostRequestChannelsInGuild(guildID string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for i, c := range repo.channels {
		if c.GuildID == guildID {
			repo.channels[i] = repo.channels[len(repo.channels)-1]
			repo.channels = repo.channels[:len(repo.channels)-1]
		}
	}
	return nil
}

func (repo *memoryRepository) GetBoostRequestChannels(guildID string) ([]*repository.BoostRequestChannel, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	channels := make([]*repository.BoostRequestChannel, 0)
	for _, brc := range repo.channels {
		if brc.GuildID == guildID {
			channels = append(channels, brc)
		}
	}
	return channels, nil
}
