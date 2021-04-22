package mock

import (
	"sync"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type memoryRepository struct {
	privileges    []*repository.AdvertiserPrivileges
	apiKeys       []*repository.APIKey
	channels      []*repository.BoostRequestChannel
	boostRequests []*repository.BoostRequest
	logChannels   map[string]string
	roleDiscounts []*repository.RoleDiscount
	stealCredits  map[string]int
	lastID        int64
	mutex         sync.Mutex
}

func NewRepository() *memoryRepository {
	return &memoryRepository{
		privileges:    make([]*repository.AdvertiserPrivileges, 0),
		apiKeys:       make([]*repository.APIKey, 0),
		channels:      make([]*repository.BoostRequestChannel, 0),
		boostRequests: make([]*repository.BoostRequest, 0),
		logChannels:   make(map[string]string),
		roleDiscounts: make([]*repository.RoleDiscount, 0),
		stealCredits:  make(map[string]int),
	}
}
