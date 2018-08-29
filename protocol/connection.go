package protocol

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"time"
)

// Connection provides a low-level interface to the Watchman service.
type Connection struct {
	reader *bufio.Reader
	socket io.Writer
	// metadata
	capabilities map[string]struct{}
	sockname     string
	version      string
}

// Connect connects to or starts the Watchman server and returns a new Connection.
func Connect() (*Connection, error) {
	sockname, err := sockname()
	if err != nil {
		return nil, err
	}

	socket, err := dial(sockname, 30*time.Second)
	if err != nil {
		return nil, err
	}

	c := &Connection{
		reader:   bufio.NewReader(socket),
		socket:   socket,
		sockname: sockname,
	}
	err = c.init()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// HasCapability checks if the Watchman server supports a specific feature.
func (c *Connection) HasCapability(capability string) bool {
	_, ok := c.capabilities[capability]
	return ok
}

// SockName returns the UNIX domain socket used to communicate with the Watchman server.
func (c *Connection) SockName() string {
	return c.sockname
}

// Version returns the version of the Watchman server.
func (c *Connection) Version() string {
	return c.version
}

func (c *Connection) init() (err error) {
	if err = c.Send(&ListCapabilitiesRequest{}); err != nil {
		return
	}

	pdu, err := c.Recv()
	if err != nil {
		return
	}

	res := NewListCapabilitiesResponse(pdu)
	capset := map[string]struct{}{}
	for _, cap := range res.Capabilities() {
		capset[cap] = struct{}{}
	}
	c.capabilities = capset
	c.version = res.Version()

	return
}

// Recv reads and decodes a response PDU from the Watchman server.
func (c *Connection) Recv() (pdu ResponsePDU, err error) {
	line, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(line, &pdu); err != nil {
		return nil, err
	} else if msg, ok := pdu["error"]; ok {
		err = &WatchmanError{msg: msg.(string)}
		return nil, err
	}

	return
}

// Send encodes and sends a request PDU to the Watchman server.
func (c *Connection) Send(req Request) (err error) {
	args := req.Args()
	b, err := json.Marshal(args)
	if err != nil {
		return
	}

	_, err = c.socket.Write(b)
	if err != nil {
		return
	}

	_, err = c.socket.Write([]byte("\n"))
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
