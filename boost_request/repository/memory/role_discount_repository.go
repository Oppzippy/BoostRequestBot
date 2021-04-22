package mock

import "github.com/oppzippy/BoostRequestBot/boost_request/repository"

func (repo *memoryRepository) GetRoleDiscountForRole(guildID, roleID string) (*repository.RoleDiscount, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for _, rd := range repo.roleDiscounts {
		if rd.GuildID == guildID && rd.RoleID == roleID {
			return rd, nil
		}
	}
	return nil, repository.ErrNoResults
}

func (repo *memoryRepository) GetRoleDiscountsForGuild(guildID string) ([]*repository.RoleDiscount, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	discounts := make([]*repository.RoleDiscount, 0)
	for _, rd := range repo.roleDiscounts {
		if rd.GuildID == guildID {
			discounts = append(discounts, rd)
		}
	}
	return discounts, nil
}

func (repo *memoryRepository) InsertRoleDiscount(rd *repository.RoleDiscount) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	rd.ID = repo.lastID
	repo.lastID += 1
	repo.roleDiscounts = append(repo.roleDiscounts, rd)
	return nil
}

func (repo *memoryRepository) DeleteRoleDiscount(rd *repository.RoleDiscount) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for i, r := range repo.roleDiscounts {
		if r.ID == rd.ID {
			repo.roleDiscounts[i] = repo.roleDiscounts[len(repo.roleDiscounts)-1]
			repo.roleDiscounts = repo.roleDiscounts[:len(repo.roleDiscounts)-1]
			break
		}
	}
	return nil
}
