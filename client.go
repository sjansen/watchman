package watchman

import (
	"github.com/sjansen/watchman/protocol"
)

// Client provides a high-level interface to Watchman.
type Client struct {
	conn      *protocol.Connection
	stop      func(bool)
	requests  chan<- protocol.Request
	responses <-chan result
	updates   <-chan interface{}
}

// Connect connects to or starts the Watchman server and returns a
// new Client.
func Connect() (c *Client, err error) {
	conn, err := protocol.Connect()
	if err != nil {
		return
	}

	loop, stop := startEventLoop(conn)
	c = &Client{
		conn:      conn,
		stop:      stop,
		requests:  loop.requests,
		responses: loop.responses,
		updates:   loop.updates,
	}
	return
}

func (c *Client) send(req protocol.Request) (res protocol.ResponsePDU, err error) {
	c.requests <- req
	result := <-c.responses
	if result.err == nil {
		res = result.pdu
	} else {
		err = result.err
	}
	return
}

// AddWatch requests that the Watchman server monitor a directory for changes.
//
// Please note that Watchman may reuse an existing watch, or choose to start
// watching a parent of the requested directory.
//
// For details, see: https://facebook.github.io/watchman/docs/cmd/watch-project.html
func (c *Client) AddWatch(path string) (*Watch, error) {
	req := &protocol.WatchProjectRequest{Path: path}
	pdu, err := c.send(req)
	if err != nil {
		return nil, err
	}

	res := protocol.NewWatchProjectResponse(pdu)
	w := &Watch{
		client: c,
		root:   res.Watch(),
	}
	return w, nil
}

// Close closes the connection to the Watchman server.
func (c *Client) Close() error {
	c.stop(false)
	return nil
}

// HasCapability checks if the Watchman server supports a feature.
//
// For details, see: https://facebook.github.io/watchman/docs/capabilities.html
func (c *Client) HasCapability(capability string) bool {
	return c.conn.HasCapability(capability)
}

// ListWatches returns a list of directories that Watchman is monitoring.
func (c *Client) ListWatches() (roots []string, err error) {
	req := &protocol.WatchListRequest{}
	if pdu, err := c.send(req); err == nil {
		res := protocol.NewWatchListResponse(pdu)
		roots = res.Roots()
	}
	return
}

// Notifications returns a channel that emits unilateral messages
// from the Watchman server.
func (c *Client) Notifications() <-chan interface{} {
	return c.updates
}

// SockName returns the location of then UNIX domain socket used
// to communicate with the Watchman server.
func (c *Client) SockName() string {
	return c.conn.SockName()
}

// Version returns the version of the Watchman server.
func (c *Client) Version() string {
	return c.conn.Version()
}
