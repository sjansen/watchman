package watchman

import (
	"io"

	"github.com/sjansen/watchman/protocol"
)

func producer(conn *protocol.Connection) <-chan protocol.ResponsePDU {
	ch := make(chan protocol.ResponsePDU)
	go func() {
		defer close(ch)

		for {
			pdu, err := conn.Recv()
			if err != nil {
				return
			}
			ch <- pdu
		}
	}()
	return ch
}

// Client provides a high-level interface to the Watchman service.
type Client struct {
	conn *protocol.Connection
	recv <-chan protocol.ResponsePDU
}

// Connect connects to or starts the Watchman server and returns a new Client.
func Connect() (c *Client, err error) {
	conn, err := protocol.Connect()
	if err != nil {
		return
	}

	c = &Client{
		conn: conn,
		recv: producer(conn),
	}
	return
}

func (c *Client) request(req protocol.Request) (res protocol.ResponsePDU, err error) {
	if err = c.conn.Send(req); err != nil {
		return
	}
	for pdu := range c.recv {
		if !pdu.IsUnilateral() {
			return pdu, nil
		}
	}
	// TODO replace EOF?
	return nil, io.EOF
}

// Close closes the connection to the Watchman server.
func (c *Client) Close() error {
	return c.conn.Close()
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
	if pdu, err := c.request(req); err == nil {
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
	if pdu, err := c.request(req); err == nil {
		res := protocol.NewWatchProjectResponse(pdu)
		w = &Watch{
			client: c,
			root:   res.Watch(),
		}
	}
	return
}
