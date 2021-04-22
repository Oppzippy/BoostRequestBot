package mock

import (
	"errors"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func (repo *memoryRepository) GetAPIKey(key string) (*repository.APIKey, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// TODO
	return nil, errors.New("memory GetApiKey not implemented")
}
