package quit

import "sync"

var (
	waiter sync.WaitGroup
	quitCh chan struct{}
)

func init() {
	quitCh = make(chan struct{})
}

func GetQuitChannel() chan struct{} {
	return quitCh
}

func GetWaiter() *sync.WaitGroup {
	return &waiter
}

// BroadcastShutdown broadcasts a message to all goroutines that are subscribing to the quit channel. This works because
// a channel close works like a message broadcast on all listeners subscribed to the channel. Ideally we can send a
// message normally, but in go message sends to channels always work in a fan-out fashion. That is, only one listener
// can get the message instead of all.
func BroadcastShutdown() {
	close(quitCh)
}
