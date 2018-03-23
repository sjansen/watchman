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

type connection struct {
	commands chan<- string
	events   <-chan string
}

type object map[string]interface{}

type eventloop struct {
	commands    chan<- string
	results     <-chan object
	unilaterals <-chan object
}

func loop(c *connection) (l *eventloop) {
	commands := make(chan string)
	results := make(chan object)
	unilaterals := make(chan object)
	l = &eventloop{
		commands:    commands,
		results:     results,
		unilaterals: unilaterals,
	}

	expectCommand := func() (ok bool) {
		for {
			select {
			case command, ok := <-commands:
				if ok {
					c.commands <- command
				}
				return ok
			case data, ok := <-c.events:
				if ok {
					var event object
					if err := json.Unmarshal([]byte(data), &event); err != nil {
						ok = false
						event = object{"error": err.Error()}
					}
					unilaterals <- event
				}
				return ok
			}
		}
	}

	expectResult := func() (ok bool) {
		for {
			data, ok := <-c.events
			if ok {
				var event object
				if err := json.Unmarshal([]byte(data), &event); err != nil {
					ok = false
					event = object{"error": err.Error()}
				}
				if _, u8l := event["log"]; u8l {
					unilaterals <- event
				} else if _, u8l := event["subscription"]; u8l {
					unilaterals <- event
				} else {
					results <- event
				}
			}
			return ok
		}
	}

	go func() {
		defer close(commands)
		defer close(results)
		defer close(unilaterals)
		for {
			if ok := expectCommand(); !ok {
				return
			}
			if ok := expectResult(); !ok {
				return
			}
		}
	}()

	return
}

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
