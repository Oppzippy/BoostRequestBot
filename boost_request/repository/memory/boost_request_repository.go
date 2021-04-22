package mock

import "github.com/oppzippy/BoostRequestBot/boost_request/repository"

func (repo *memoryRepository) GetUnresolvedBoostRequests() ([]*repository.BoostRequest, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	requests := make([]*repository.BoostRequest, 0)
	for _, br := range repo.boostRequests {
		if !br.IsResolved {
			requests = append(requests, br)
		}
	}
	return requests, nil
}

func (repo *memoryRepository) GetBoostRequestByBackendMessageID(backendChannelID, backendMessageID string) (*repository.BoostRequest, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for _, br := range repo.boostRequests {
		if br.Channel.BackendChannelID == backendChannelID && br.BackendMessageID == backendMessageID {
			return br, nil
		}
	}
	return nil, repository.ErrNoResults
}

func (repo *memoryRepository) InsertBoostRequest(br *repository.BoostRequest) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	br.ID = repo.lastID
	repo.lastID += 1
	repo.boostRequests = append(repo.boostRequests, br)
	return nil
}

func (repo *memoryRepository) ResolveBoostRequest(br *repository.BoostRequest) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for i, b := range repo.boostRequests {
		if b.ID == br.ID {
			repo.boostRequests[i] = br
			break
		}
	}
	return nil
}
