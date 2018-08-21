package watchman

import (
	"github.com/sjansen/watchman/connection"
)

type Client struct {
	Conn *connection.Connection
}

func (c *Client) HasCapability(capability string) bool {
	return c.Conn.HasCapability(capability)
}

func (c *Client) SockName() string {
	return c.Conn.SockName()
}

func (c *Client) Version() string {
	return c.Conn.Version()
}

func (c *Client) WatchList() (roots []string, err error) {
	req := &connection.WatchListRequest{}
	if err = c.Conn.Send(req); err != nil {
		return
	}

	res := &connection.WatchListResponse{}
	for {
		if unilateral, err := c.Conn.Recv(res); err != nil {
			return nil, err
		} else if unilateral == nil {
			break
		}
	}

	roots = res.Roots()
	return
}

func (c *Client) WatchProject(path string) (w *Watch, err error) {
	req := &connection.WatchProjectRequest{Path: path}
	if err = c.Conn.Send(req); err != nil {
		return
	}

	res := &connection.WatchProjectResponse{}
	for {
		if unilateral, err := c.Conn.Recv(res); err != nil {
			return nil, err
		} else if unilateral == nil {
			break
		}
	}

	w = &Watch{
		conn: c.Conn,
		root: res.Watch(),
	}
	return
}
