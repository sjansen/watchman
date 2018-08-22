package watchman

import (
	"github.com/sjansen/watchman/connection"
)

// A Watch represents a directory, or watched root, that Watchman is watching for changes.
type Watch struct {
	conn *connection.Connection
	root string
}

// Clock returns the current clock value for a watched root.
//
// For details, see: https://facebook.github.io/watchman/docs/cmd/clock.html
func (w *Watch) Clock(syncTimeout int) (clock string, err error) {
	req := &connection.ClockRequest{
		Path:        w.root,
		SyncTimeout: syncTimeout,
	}
	if err = w.conn.Send(req); err != nil {
		return
	}

	res := &connection.ClockResponse{}
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
