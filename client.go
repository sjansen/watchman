package watchman

import (
	"github.com/sjansen/watchman/protocol"
)

// Client provides a high-level interface to the Watchman service.
type Client struct {
	conn        *protocol.Connection
	stop        func(bool)
	requests    chan<- protocol.Request
	responses   <-chan result
	unilaterals <-chan protocol.ResponsePDU
}

// Connect connects to or starts the Watchman server and returns a new Client.
func Connect() (c *Client, err error) {
	conn, err := protocol.Connect()
	if err != nil {
		return
	}

	loop, stop := startEventLoop(conn)
	go func() { // TODO stop throwing away unilaterals
		for range loop.unilaterals {
			continue
		}
	}()

	c = &Client{
		conn:        conn,
		stop:        stop,
		requests:    loop.requests,
		responses:   loop.responses,
		unilaterals: loop.unilaterals,
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

// Close closes the connection to the Watchman server.
func (c *Client) Close() error {
	c.stop(false)
	return nil
}

// HasCapability checks if the Watchman server supports a specific feature.
func (c *Client) HasCapability(capability string) bool {
	return c.conn.HasCapability(capability)
}

// SockName returns the UNIX domain socket used to communicate with the Watchman server.
func (c *Client) SockName() string {
	return c.conn.SockName()
}

// Version returns the version of the Watchman server.
func (c *Client) Version() string {
	return c.conn.Version()
}

// WatchList returns a list of the dirs the Watchman server is watching.
func (c *Client) WatchList() (roots []string, err error) {
	req := &protocol.WatchListRequest{}
	if pdu, err := c.send(req); err == nil {
		res := protocol.NewWatchListResponse(pdu)
		roots = res.Roots()
	}
	return
}

// WatchProject requests that the Watchman server monitor a dir or one of its parents for changes.
//
// For details, see: https://facebook.github.io/watchman/docs/cmd/watch-project.html
func (c *Client) WatchProject(path string) (w *Watch, err error) {
	req := &protocol.WatchProjectRequest{Path: path}
	if pdu, err := c.send(req); err == nil {
		res := protocol.NewWatchProjectResponse(pdu)
		w = &Watch{
			client: c,
			root:   res.Watch(),
		}
	}
	return
}
