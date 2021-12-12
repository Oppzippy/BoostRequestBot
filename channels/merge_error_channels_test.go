package channels_test

import (
	"errors"
	"testing"

	"github.com/oppzippy/BoostRequestBot/channels"
)

func TestMergeClosedErrorChannels(t *testing.T) {
	t.Parallel()
	c1, c2 := make(chan error), make(chan error)
	close(c1)
	close(c2)

	c := channels.MergeErrorChannels(c1, c2)
	_, ok := <-c
	if ok {
		t.Error("received message from channel that should be closed")
	}
}

func TestMultipleMessagesToErrorChannels(t *testing.T) {
	t.Parallel()
	c1, c2, c3 := make(chan error), make(chan error), make(chan error)
	c := channels.MergeErrorChannels(c1, c2, c3)

	c1 <- errors.New("test")
	c2 <- errors.New("test2")
	c3 <- errors.New("test3")

	var i int
	for range c {
		i++
	}
	if i != 3 {
		t.Errorf("expected 3 messages, got %d", i)
	}
}

func TestMultipleMessagesToSameErrorChannel(t *testing.T) {
	t.Parallel()
	c1 := make(chan error)
	c2 := make(chan error)
	c := channels.MergeErrorChannels(c1, c2)
	close(c2)

	go func() {
		c1 <- errors.New("test1")
		c1 <- errors.New("test2")
		c1 <- errors.New("test3")
	}()

	var i int
	for range c {
		i++
	}
	if i != 1 {
		t.Errorf("expected 1 message, got %d", i)
	}
}
