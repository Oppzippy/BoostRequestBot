package boost_request_test

import (
	"testing"
	"time"

	"github.com/oppzippy/BoostRequestBot/boost_request"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
)

func TestImmediateSignup(t *testing.T) {
	start := time.Now()
	c := make(chan struct{})
	ar := boost_request.NewActiveRequest(repository.BoostRequest{
		Channel: repository.BoostRequestChannel{
			GuildID:           "guild",
			FrontendChannelID: "frontend",
			BackendChannelID:  "backend",
		},
		RequesterID:      "requester",
		BackendMessageID: "backendMessage",
		Message:          "I would like one boost please!",
		CreatedAt:        start,
	}, func(br repository.BoostRequest, userID string) {
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
	c := make(chan struct{})
	ar := boost_request.NewActiveRequest(repository.BoostRequest{
		Channel: repository.BoostRequestChannel{
			GuildID:           "guild",
			FrontendChannelID: "frontend",
			BackendChannelID:  "backend",
		},
		RequesterID:      "requester",
		BackendMessageID: "backendMessage",
		Message:          "I would like one boost please!",
		CreatedAt:        time.Now(),
	}, func(br repository.BoostRequest, userID string) {
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
	c := make(chan struct{})
	ar := boost_request.NewActiveRequest(repository.BoostRequest{
		Channel: repository.BoostRequestChannel{
			GuildID:           "guild",
			FrontendChannelID: "frontend",
			BackendChannelID:  "backend",
		},
		RequesterID:      "requester",
		BackendMessageID: "backendMessage",
		Message:          "I would like one boost please!",
		CreatedAt:        time.Now(),
	}, func(br repository.BoostRequest, userID string) {
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
	c := make(chan struct{})
	ar := boost_request.NewActiveRequest(repository.BoostRequest{
		Channel: repository.BoostRequestChannel{
			GuildID:           "guild",
			FrontendChannelID: "frontend",
			BackendChannelID:  "backend",
		},
		RequesterID:      "requester",
		BackendMessageID: "backendMessage",
		Message:          "I would like one boost please!",
		CreatedAt:        time.Now(),
	}, func(br repository.BoostRequest, userID string) {
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
