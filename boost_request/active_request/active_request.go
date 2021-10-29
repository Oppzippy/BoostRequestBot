package active_request

import (
	"sync"
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/oppzippy/BoostRequestBot/roll"
)

type AdvertiserChosenEvent struct {
	BoostRequest repository.BoostRequest
	UserID       string
	RollResults  *roll.WeightedRollResults
}

type ActiveRequest struct {
	AdvertiserChosenCallback func(*AdvertiserChosenEvent)
	boostRequest             repository.BoostRequest
	signupsByDelay           map[int][]userWithPrivileges
	userDelays               map[string]int
	quit                     chan struct{}
	mutex                    sync.Mutex
	inactive                 bool
}

type userWithPrivileges struct {
	userID     string
	privileges repository.AdvertiserPrivileges
}

func NewActiveRequest(br repository.BoostRequest, cb func(*AdvertiserChosenEvent)) *ActiveRequest {
	return &ActiveRequest{
		AdvertiserChosenCallback: cb,
		boostRequest:             br,
		quit:                     make(chan struct{}),
		signupsByDelay:           make(map[int][]userWithPrivileges),
		userDelays:               make(map[string]int),
	}
}

func (r *ActiveRequest) AddSignup(userID string, privileges repository.AdvertiserPrivileges) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.inactive {
		return
	}
	_, isSignedUp := r.userDelays[userID]
	if isSignedUp {
		r.removeSignupWithoutLocking(userID)
	}

	endTime := r.boostRequest.CreatedAt.Add(time.Duration(privileges.Delay) * time.Second)
	if now := time.Now(); now.After(endTime) || now.Equal(endTime) {
		r.setAdvertiserWithoutLocking(&AdvertiserChosenEvent{
			BoostRequest: r.boostRequest,
			UserID:       userID,
		})
		return
	}

	r.userDelays[userID] = privileges.Delay
	uwp := userWithPrivileges{userID, privileges}

	if signups := r.signupsByDelay[privileges.Delay]; signups != nil {
		r.signupsByDelay[privileges.Delay] = append(signups, uwp)
	} else {
		r.signupsByDelay[privileges.Delay] = []userWithPrivileges{uwp}
		go r.waitForDelay(privileges.Delay, endTime)
	}
}

func (r *ActiveRequest) RemoveSignup(userID string) {
	r.mutex.Lock()
	r.removeSignupWithoutLocking(userID)
	r.mutex.Unlock()
}

func (r *ActiveRequest) removeSignupWithoutLocking(userID string) {
	if r.inactive {
		return
	}

	delay, isSignedUp := r.userDelays[userID]
	if !isSignedUp {
		return
	}
	delete(r.userDelays, userID)
	signups, ok := r.signupsByDelay[delay]
	if !ok {
		return
	}

	for i, user := range signups {
		if user.userID == userID {
			length := len(signups)
			signups[i] = signups[length-1]
			r.signupsByDelay[delay] = signups[:length-1]
			break
		}
	}
}

func (r *ActiveRequest) SetAdvertiser(userID string) (ok bool) {
	r.mutex.Lock()
	ok = r.setAdvertiserWithoutLocking(&AdvertiserChosenEvent{
		BoostRequest: r.boostRequest,
		UserID:       userID,
	})
	r.mutex.Unlock()
	return ok
}

func (r *ActiveRequest) setAdvertiserWithoutLocking(event *AdvertiserChosenEvent) (ok bool) {
	ok = !r.inactive
	if ok {
		close(r.quit)
		r.inactive = true
		go r.AdvertiserChosenCallback(event)
	}
	return ok
}

// mutex should be locked before calling this method
func (r *ActiveRequest) chooseAdvertiser(delay int) (rollInfo *roll.WeightedRollResults, ok bool) {
	users := r.signupsByDelay[delay]
	if len(users) == 0 {
		return nil, false
	}

	weightedRoll := roll.NewWeightedRoll(len(users))
	for _, user := range users {
		weightedRoll.AddItem(user.userID, user.privileges.Weight)
	}
	results, ok := weightedRoll.Roll()
	return results, ok
}

func (r *ActiveRequest) waitForDelay(delay int, endTime time.Time) {
	select {
	case <-r.quit:
		return
	case <-time.After(time.Until(endTime)):
	}
	r.mutex.Lock()
	if !r.inactive {
		if rollResults, ok := r.chooseAdvertiser(delay); ok {
			advertiserID := rollResults.ChosenItem()
			r.setAdvertiserWithoutLocking(&AdvertiserChosenEvent{
				BoostRequest: r.boostRequest,
				UserID:       advertiserID,
				RollResults:  rollResults,
			})
		}
	}
	r.mutex.Unlock()
}

func (r *ActiveRequest) Destroy() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if !r.inactive {
		close(r.quit)
		r.inactive = true
	}
}
