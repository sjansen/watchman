package watchman

import (
	"github.com/sjansen/watchman/protocol"
)

// Client provides a high-level interface to the Watchman service.
type Client struct {
	conn *protocol.Connection
}

// Connect connects to or starts the Watchman server and returns a new Client.
func Connect() (c *Client, err error) {
	conn, err := protocol.Connect()
	if err != nil {
		return
	}
	c = &Client{conn: conn}
	return
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
	if err = c.conn.Send(req); err != nil {
		return
	}

	res := &protocol.WatchListResponse{}
	for {
		if unilateral, err := c.conn.Recv(res); err != nil {
			return nil, err
		} else if unilateral == nil {
			break
		}
	}

	roots = res.Roots()
	return
}

// WatchProject requests that the Watchman server monitor a dir or one of its parents for changes.
//
// For details, see: https://facebook.github.io/watchman/docs/cmd/watch-project.html
func (c *Client) WatchProject(path string) (w *Watch, err error) {
	req := &protocol.WatchProjectRequest{Path: path}
	if err = c.conn.Send(req); err != nil {
		return
	}

	res := &protocol.WatchProjectResponse{}
	for {
		if unilateral, err := c.conn.Recv(res); err != nil {
			return nil, err
		} else if unilateral == nil {
			break
		}
	}

	w = &Watch{
		conn: c.conn,
		root: res.Watch(),
	}
	return
}
