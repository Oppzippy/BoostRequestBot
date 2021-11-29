package active_request_test

import (
	"math"
	"testing"
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request/active_request"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestImmediateSignup(t *testing.T) {
	t.Parallel()
	start := time.Now()
	c := make(chan struct{})
	ar := active_request.NewActiveRequest(repository.BoostRequest{
		Channel: &repository.BoostRequestChannel{
			GuildID:           "guild",
			FrontendChannelID: "frontend",
			BackendChannelID:  "backend",
		},
		RequesterID: "requester",
		BackendMessages: []*repository.BoostRequestBackendMessage{
			{
				ChannelID: "backendChannelID",
				MessageID: "backendMessageID",
			},
		},
		Message:   "I would like one boost please!",
		CreatedAt: start,
	}, func(event *active_request.AdvertiserChosenEvent) {
		c <- struct{}{}
	})

	ar.AddSignup("advertiser", repository.AdvertiserPrivileges{
		GuildID: "guild",
		RoleID:  "advertiser",
		Weight:  1,
		Delay:   1,
	})

	select {
	case <-c:
		diff := time.Since(start)
		if diff < time.Duration(900)*time.Millisecond {
			t.Errorf("boost request resolved too fast: %v", diff)
			return
		}
		if diff > time.Duration(1500)*time.Millisecond {
			t.Errorf("boost request resolved too slow: %v", diff)
			return
		}
	case <-time.After(5 * time.Second):
		t.Error("timed out")
	}
}

func TestLateSignup(t *testing.T) {
	t.Parallel()
	c := make(chan struct{})
	ar := active_request.NewActiveRequest(repository.BoostRequest{
		Channel: &repository.BoostRequestChannel{
			GuildID:           "guild",
			FrontendChannelID: "frontend",
			BackendChannelID:  "backend",
		},
		RequesterID: "requester",
		BackendMessages: []*repository.BoostRequestBackendMessage{
			{
				ChannelID: "backendChannelID",
				MessageID: "backendMessageID",
			},
		},
		Message:   "I would like one boost please!",
		CreatedAt: time.Now(),
	}, func(event *active_request.AdvertiserChosenEvent) {
		c <- struct{}{}
	})
	<-time.After(2 * time.Second)

	ar.AddSignup("advertiser", repository.AdvertiserPrivileges{
		GuildID: "guild",
		RoleID:  "advertiser",
		Weight:  1,
		Delay:   1,
	})
	select {
	case <-time.After(250 * time.Millisecond):
		t.Error("the advertiser was not accepted immediately")
	case <-c:
	}
}

func TestSetAdvertiser(t *testing.T) {
	t.Parallel()
	c := make(chan struct{})
	ar := active_request.NewActiveRequest(repository.BoostRequest{
		Channel: &repository.BoostRequestChannel{
			GuildID:           "guild",
			FrontendChannelID: "frontend",
			BackendChannelID:  "backend",
		},
		RequesterID: "requester",
		BackendMessages: []*repository.BoostRequestBackendMessage{
			{
				ChannelID: "backendChannelID",
				MessageID: "backendMessageID",
			},
		},
		Message:   "I would like one boost please!",
		CreatedAt: time.Now(),
	}, func(event *active_request.AdvertiserChosenEvent) {
		c <- struct{}{}
	})

	ar.SetAdvertiser("advertiser")
	select {
	case <-time.After(250 * time.Millisecond):
		t.Error("timed out")
	case <-c:
	}
}

func TestRepeatedSetAdvertiser(t *testing.T) {
	t.Parallel()
	c := make(chan struct{})
	ar := active_request.NewActiveRequest(repository.BoostRequest{
		Channel: &repository.BoostRequestChannel{
			GuildID:           "guild",
			FrontendChannelID: "frontend",
			BackendChannelID:  "backend",
		},
		RequesterID: "requester",
		BackendMessages: []*repository.BoostRequestBackendMessage{
			{
				ChannelID: "backendChannelID",
				MessageID: "backendMessageID",
			},
		},
		Message:   "I would like one boost please!",
		CreatedAt: time.Now(),
	}, func(event *active_request.AdvertiserChosenEvent) {
		c <- struct{}{}
	})

	ar.SetAdvertiser("advertiser")
	ar.SetAdvertiser("advertiser2")
	ar.SetAdvertiser("advertiser3")

	var i int
	var done bool
	for !done {
		select {
		case <-time.After(250 * time.Millisecond):
			done = true
		case <-c:
			i++
		}
	}

	if i != 1 {
		t.Errorf("expected 1 advertiser to be chosen, but %d were chosen", i)
	}
}

func TestRepeatedSignupOfSameUser(t *testing.T) {
	t.Parallel()
	c := make(chan struct{})
	ar := active_request.NewActiveRequest(repository.BoostRequest{
		Channel: &repository.BoostRequestChannel{
			GuildID:           "guild",
			FrontendChannelID: "frontend",
			BackendChannelID:  "backend",
		},
		RequesterID: "requester",
		BackendMessages: []*repository.BoostRequestBackendMessage{
			{
				ChannelID: "backendChannelID",
				MessageID: "backendMessageID",
			},
		},
		Message:   "I would like one boost please!",
		CreatedAt: time.Now(),
	}, func(event *active_request.AdvertiserChosenEvent) {
		c <- struct{}{}
	})

	p := repository.AdvertiserPrivileges{
		GuildID: "guild",
		RoleID:  "advertiser",
		Weight:  1,
		Delay:   1,
	}

	ar.AddSignup("advertiser", p)
	ar.AddSignup("advertiser", p)
	ar.AddSignup("advertiser", p)
	ar.RemoveSignup("advertiser")

	var i int
	var done bool
	for !done {
		select {
		case <-time.After(1500 * time.Millisecond):
			done = true
		case <-c:
			i++
		}
	}

	if i != 0 {
		t.Errorf("expected 0 advertisers to be chosen, but %d were chosen", i)
	}
}

func TestRandomness(t *testing.T) {
	t.Parallel()
	numRuns := 10000
	winners := make(chan string)
	for i := 0; i < numRuns; i++ {
		runIteration(winners)
	}
	winCount := make(map[string]int)
	for i := 0; i < numRuns; i++ {
		winner := <-winners
		winCount[winner] = winCount[winner] + 1
	}

	advertiser1WinRate := float64(winCount["advertiser1"]) / float64(numRuns)
	advertiser2WinRate := float64(winCount["advertiser2"]) / float64(numRuns)
	advertiser3WinRate := float64(winCount["advertiser3"]) / float64(numRuns)

	if math.Abs(advertiser1WinRate-1.0/6.0) > 0.01 {
		t.Errorf("Expected %f win rate for advertiser1, got %f", 1.0/6.0, advertiser1WinRate)
	}
	if math.Abs(advertiser2WinRate-2.0/6.0) > 0.01 {
		t.Errorf("Expected %f win rate for advertiser2, got %f", 2.0/6.0, advertiser2WinRate)
	}
	if math.Abs(advertiser3WinRate-3.0/6.0) > 0.01 {
		t.Errorf("Expected %f win rate for advertiser3, got %f", 3.0/6.0, advertiser3WinRate)
	}
}

func runIteration(winners chan string) {
	ar := active_request.NewActiveRequest(repository.BoostRequest{
		Channel: &repository.BoostRequestChannel{
			GuildID:           "guild",
			FrontendChannelID: "frontend",
			BackendChannelID:  "backend",
		},
		RequesterID: "requester",
		BackendMessages: []*repository.BoostRequestBackendMessage{
			{
				ChannelID: "backendChannelID",
				MessageID: "backendMessageID",
			},
		},
		Message:   "I would like one boost please!",
		CreatedAt: time.Now(),
	}, func(event *active_request.AdvertiserChosenEvent) {
		winners <- event.UserID
	})

	ar.AddSignup("advertiser1", repository.AdvertiserPrivileges{
		GuildID: "guild",
		RoleID:  "advertiser1",
		Weight:  1,
		Delay:   1,
	})
	ar.AddSignup("advertiser2", repository.AdvertiserPrivileges{
		GuildID: "guild",
		RoleID:  "advertiser2",
		Weight:  2,
		Delay:   1,
	})
	ar.AddSignup("advertiser3", repository.AdvertiserPrivileges{
		GuildID: "guild",
		RoleID:  "advertiser3",
		Weight:  3,
		Delay:   1,
	})
}
