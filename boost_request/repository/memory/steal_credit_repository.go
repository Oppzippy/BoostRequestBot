package memory

import (
	"fmt"
	"strings"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *memoryRepository) GetStealCreditsForUser(guildID, userID string) (int, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	credits := repo.stealCredits[fmt.Sprintf("%s:%s", guildID, userID)]
	return credits, nil
}

func (repo *memoryRepository) GetGlobalStealCreditsForUser(userID string) (map[string]int, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	creditsByGuild := make(map[string]int)
	for key, credits := range repo.stealCredits {
		parts := strings.Split(key, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected one ':' in key, found %d", len(parts)-1)
		}
		if parts[1] == userID {
			creditsByGuild[parts[0]] = credits
		}
	}
	return creditsByGuild, nil
}

func (repo *memoryRepository) AdjustStealCreditsForUser(guildID, userID string, operation repository.Operation, amount int) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	key := fmt.Sprintf("%s:%s", guildID, userID)

	switch operation {
	case repository.OperationSet:
		repo.stealCredits[key] = amount
	case repository.OperationAdd:
		repo.stealCredits[key] += amount
	case repository.OperationSubtract:
		repo.stealCredits[key] -= amount
	case repository.OperationMultiply:
		repo.stealCredits[key] *= amount
	case repository.OperationDivide:
		repo.stealCredits[key] /= amount
	default:
		return repository.ErrInvalidOperation
	}
	return nil
}

func (repo *memoryRepository) UpdateStealCreditsForUser(guildID, userID string, amount int) error {
	return repo.AdjustStealCreditsForUser(guildID, userID, repository.OperationSet, amount)
}
