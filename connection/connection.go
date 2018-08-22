package connection

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

// Request is the interface used to encode a request PDU.
type Request interface {
	Args() []interface{}
}

// Response is the interface common to all response PDUs.
type Response interface {
	Version() string
	Warning() string
}

type response struct {
	Version string
	Warning string
}

// Unilateral is the interface that wraps the PDU method.
//
// PDU provides access to response data decoded to primitive Go values.
// This is most useful when Connection.Recv() is unable to return a
// more specific Unilateral implementation.
type Unilateral interface {
	PDU() map[string]interface{}
}

type unilateral struct {
	pdu map[string]json.RawMessage
}

func (u *unilateral) PDU() map[string]interface{} {
	pdu := map[string]interface{}{}
	for key, raw := range u.pdu {
		var val interface{}
		if err := json.Unmarshal(raw, &val); err == nil {
			pdu[key] = val
		}
	}
	return pdu
}

// Connection provides a low-level interface to the Watchman service.
type Connection struct {
	reader *bufio.Reader
	socket io.Writer
	// metadata
	capabilities map[string]struct{}
	sockname     string
	version      string
}

// New connects to or starts the Watchman server and returns a new Connection.
func New() (*Connection, error) {
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

	res := &ListCapabilitiesResponse{}
	if _, err = c.Recv(res); err != nil {
		return
	}

	capset := map[string]struct{}{}
	for _, cap := range res.Capabilities() {
		capset[cap] = struct{}{}
	}
	c.capabilities = capset
	c.version = res.Version()

	return
}

// Recv reads and decodes a response PDU from the Watchman server.
// A non-nil Unilateral return value indicates the PDU was sent as
// a result of an older request, instead of the immediately previous
// request. Most notably, as a result of the subscribe command.
//
// When the PDU was sent as a response to to immediately previous
// request, the provided Response will be used to decode the PDU.
//
// When the PDU was sent unilaterally, Recv will attempt to return the
// most specific Unilateral implementation available. Type assertion
// should be used to determine that implementation. Currently, only
// Subscription responses are recognized, but additional Unilateral
// implementations may be added in the future.
func (c *Connection) Recv(res Response) (Unilateral, error) {
	line, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	var pdu map[string]json.RawMessage
	if err = json.Unmarshal(line, &pdu); err != nil {
		return nil, err
	} else if msg, ok := pdu["error"]; ok {
		err = &WatchmanError{string(msg)}
		return nil, err
	} else if _, ok := pdu["unilateral"]; ok {
		if _, ok := pdu["subscription"]; ok {
			sub := &Subscription{
				unilateral: unilateral{pdu: pdu},
			}
			err = json.Unmarshal(line, sub)
			return sub, err
		}
		return &unilateral{pdu: pdu}, nil
	}

	err = json.Unmarshal(line, interface{}(res))
	return nil, err
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
