package watchman

import (
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
func (w *Watch) Clock(syncTimeout int) (clock string, err error) {
	req := &protocol.ClockRequest{
		Path:        w.root,
		SyncTimeout: syncTimeout,
	}
	res := &protocol.ClockResponse{}
	if err = w.client.handle(req, res); err == nil {
		clock = res.Clock()
	}
	return
}
