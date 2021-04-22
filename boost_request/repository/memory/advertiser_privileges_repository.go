package mock

import (
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *memoryRepository) GetAdvertiserPrivilegesForGuild(guildID string) ([]*repository.AdvertiserPrivileges, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	priveleges := make([]*repository.AdvertiserPrivileges, 0)
	for _, p := range repo.privileges {
		if p.GuildID == guildID {
			priveleges = append(priveleges, p)
		}
	}
	return priveleges, nil
}

func (repo *memoryRepository) GetAdvertiserPrivilegesForRole(guildID, roleID string) (*repository.AdvertiserPrivileges, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for _, p := range repo.privileges {
		if p.GuildID == guildID && p.RoleID == roleID {
			return p, nil
		}
	}
	return nil, repository.ErrNoResults
}

func (repo *memoryRepository) InsertAdvertiserPrivileges(privileges *repository.AdvertiserPrivileges) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	privileges.ID = repo.lastID
	repo.lastID += 1
	repo.privileges = append(repo.privileges, privileges)
	return nil
}

func (repo *memoryRepository) DeleteAdvertiserPrivileges(privileges *repository.AdvertiserPrivileges) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for i, p := range repo.privileges {
		if p.ID == privileges.ID {
			repo.privileges[i] = repo.privileges[len(repo.privileges)-1]
			repo.privileges = repo.privileges[:len(repo.privileges)-1]
			break
		}
	}
	return nil
}
