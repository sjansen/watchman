package watchman

import (
	"time"

	"github.com/sjansen/watchman/protocol"
)

// A Watch represents a directory, or watched root, that Watchman is watching for changes.
type Watch struct {
	client *Client
	root   string
}

// Clock returns the current clock value for a watched root.
//
// For details, see: https://facebook.github.io/watchman/docs/cmd/clock.html
func (w *Watch) Clock(syncTimeout time.Duration) (clock string, err error) {
	timeout := syncTimeout.Nanoseconds() / int64(time.Millisecond)
	req := &protocol.ClockRequest{
		Path:        w.root,
		SyncTimeout: int(timeout),
	}
	pdu, err := w.client.send(req)
	if err == nil {
		res := protocol.NewClockResponse(pdu)
		clock = res.Clock()
	}
	return
}

// SubscribeOptions configures an optional set of parameters for Subscribe.
type SubscribeOptions struct {
	EmptyOnFreshInstance bool
}

// Subscribe requests notification when changes occur under a watched root.
func (w *Watch) Subscribe(name, root string, opts ...*SubscribeOptions) (s *Subscription, err error) {
	req := &protocol.SubscribeRequest{
		Name: name,
		Root: root,
	}
	if len(opts) > 0 && opts[0] != nil {
		req.EmptyOnFreshInstance = opts[0].EmptyOnFreshInstance
	}
	_, err = w.client.send(req)
	if err == nil {
		s = &Subscription{
			client: w.client,
			name:   name,
			root:   root,
		}
	}
	return
}
