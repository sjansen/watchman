package watchman

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"time"
)

type Connection struct {
	commands chan<- []interface{}
	results  <-chan result
	// metadata
	capabilities map[string]struct{}
	sockname     string
	version      string
}

func Connect(ctx context.Context) (*Connection, error) {
	sockname, err := sockname()
	if err != nil {
		return nil, err
	}

	socket, err := dial(sockname, 30*time.Second)
	if err != nil {
		return nil, err
	}

	c := &Connection{
		commands: writer(ctx, socket),
		results:  reader(ctx, socket),
		sockname: sockname,
	}
	err = c.init()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connection) HasCapability(capability string) bool {
	_, ok := c.capabilities[capability]
	return ok
}

func (c *Connection) SockName() string {
	return c.sockname
}

func (c *Connection) Version() string {
	return c.version
}

func (c *Connection) command(args ...interface{}) (map[string]interface{}, error) {
	c.commands <- args
	result := <-c.results
	return result.resp, result.err
}

func (c *Connection) init() (err error) {
	resp, err := c.command("list-capabilities")
	if err != nil {
		return
	}

	if version, ok := resp["version"].(string); ok {
		c.version = version
	}

	if capabilities, ok := resp["capabilities"].([]interface{}); ok {
		capset := map[string]struct{}{}
		for _, cap := range capabilities {
			capset[cap.(string)] = struct{}{}
		}
		c.capabilities = capset
	}

	return
}

func sockname() (string, error) {
	sockname := os.Getenv("WATCHMAN_SOCK")
	if sockname != "" {
		return sockname, nil
	}

	buffer := &bytes.Buffer{}
	cmd := exec.Command("watchman", "get-sockname")
	cmd.Stdout = buffer
	if err := cmd.Run(); err != nil {
		return "", err
	}

	var resp map[string]string
	if err := json.NewDecoder(buffer).Decode(&resp); err != nil {
		return "", err
	}

	sockname, ok := resp["sockname"]
	if !ok {
		return "", errors.New("unable to find watchman socket")
	}

	return sockname, nil
}
