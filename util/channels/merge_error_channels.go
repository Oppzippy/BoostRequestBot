package channels

import "sync"

// If channels send 0 messages, they must be closed. Only the first message is captured per channel.
// Once all channels have sent a message or are closed, the returned channel closes.
func MergeErrorChannels(channels ...<-chan error) <-chan error {
	mergedChannel := make(chan error, len(channels))

	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	for _, c := range channels {
		currentChannel := c
		go func() {
			err, ok := <-currentChannel
			if ok {
				mergedChannel <- err
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(mergedChannel)
	}()

	return mergedChannel
}
