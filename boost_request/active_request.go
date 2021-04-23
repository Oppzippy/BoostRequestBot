package boost_request

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

type activeRequest struct {
	AdvertiserChosenCallback func(br repository.BoostRequest, userID string)
	boostRequest             repository.BoostRequest
	signupsByDelay           map[int][]userWithPrivileges
	quit                     chan struct{}
	mutex                    sync.Mutex
	inactive                 bool
}

type userWithPrivileges struct {
	userID     string
	privileges repository.AdvertiserPrivileges
}

func NewActiveRequest(br repository.BoostRequest, cb func(br repository.BoostRequest, userID string)) *activeRequest {
	return &activeRequest{
		AdvertiserChosenCallback: cb,
		boostRequest:             br,
		quit:                     make(chan struct{}),
		signupsByDelay:           make(map[int][]userWithPrivileges),
	}
}

func (r *activeRequest) AddSignup(userID string, privileges repository.AdvertiserPrivileges) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.inactive {
		return
	}

	endTime := r.boostRequest.CreatedAt.Add(time.Duration(privileges.Delay) * time.Second)
	if now := time.Now(); now.After(endTime) || now.Equal(endTime) {
		r.setAdvertiserWithoutLocking(userID)
		return
	}

	uwp := userWithPrivileges{userID, privileges}

	if signups := r.signupsByDelay[privileges.Delay]; signups != nil {
		r.signupsByDelay[privileges.Delay] = append(signups, uwp)
	} else {
		r.signupsByDelay[privileges.Delay] = []userWithPrivileges{uwp}
		go r.waitForDelay(privileges.Delay, endTime)
	}
}

func (r *activeRequest) SetAdvertiser(userID string) (ok bool) {
	r.mutex.Lock()
	ok = r.setAdvertiserWithoutLocking(userID)
	r.mutex.Unlock()
	return ok
}

func (r *activeRequest) setAdvertiserWithoutLocking(userID string) (ok bool) {
	ok = !r.inactive
	if ok {
		close(r.quit)
		r.inactive = true
		go r.AdvertiserChosenCallback(r.boostRequest, userID)
	}
	return ok
}

// mutex should be locked before calling this method
func (r *activeRequest) chooseAdvertiser(delay int) (string, bool) {
	users := r.signupsByDelay[delay]
	var totalWeight float64
	for _, user := range users {
		totalWeight += user.privileges.Weight
	}
	var chosenWeight float64 = rand.Float64() * totalWeight

	var currentWeight float64
	for _, user := range users {
		currentWeight += user.privileges.Weight
		if chosenWeight < currentWeight {
			return user.userID, true
		}
	}
	return "", false
}

func (r *activeRequest) waitForDelay(delay int, endTime time.Time) {
	select {
	case <-r.quit:
		return
	case <-time.After(time.Until(endTime)):
	}
	r.mutex.Lock()
	if !r.inactive {
		if advertiserID, ok := r.chooseAdvertiser(delay); ok {
			r.setAdvertiserWithoutLocking(advertiserID)
		}
	}
	r.mutex.Unlock()
}
