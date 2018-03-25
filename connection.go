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

type command []interface{}

type Connection struct {
	commands chan<- string
	results  <-chan object
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

	server := serverFromSocket(ctx, socket)
	l, stop := loop(server)
	defer stop(false)

	c := &Connection{
		commands: l.commands,
		results:  l.results,
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

func (c *Connection) command(args ...interface{}) (object, error) {
	command, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	c.commands <- string(command)
	event := <-c.results
	if msg, ok := event["error"]; ok {
		return event, errors.New(msg.(string))
	}
	return event, nil
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
