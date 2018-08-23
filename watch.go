package watchman

import (
	"github.com/sjansen/watchman/protocol"
)

// A Watch represents a directory, or watched root, that Watchman is watching for changes.
type Watch struct {
	conn *protocol.Connection
	root string
}

// Clock returns the current clock value for a watched root.
//
// For details, see: https://facebook.github.io/watchman/docs/cmd/clock.html
func (w *Watch) Clock(syncTimeout int) (clock string, err error) {
	req := &protocol.ClockRequest{
		Path:        w.root,
		SyncTimeout: syncTimeout,
	}
	if err = w.conn.Send(req); err != nil {
		return
	}

	res := &protocol.ClockResponse{}
	for {
		if unilateral, err := w.conn.Recv(res); err != nil {
			return "", err
		} else if unilateral == nil {
			break
		}
	}

	clock = res.Clock()
	return
}
