package protocol

/*
$ watchman list-capabilities
{
    "version": "3.8.0",
    "capabilities": [
        "field-mode",
        "term-allof",
        "cmd-trigger"
    ]
}
*/

// A ListCapabilitiesRequest represents the Watchman list-capabilities command.
//
// See also: https://facebook.github.io/watchman/docs/cmd/list-capabilities.html
type ListCapabilitiesRequest struct{}

// Args returns values used to encode a request PDU.
func (req *ListCapabilitiesRequest) Args() []interface{} {
	return []interface{}{"list-capabilities"}
}

type listCapabilitiesResponse struct {
	response
	Capabilities []string
}

// A ListCapabilitiesResponse represents a response to the Watchman list-capabilities command.
type ListCapabilitiesResponse struct {
	listCapabilitiesResponse
}

// Version returns the Watchman server version.
func (res *ListCapabilitiesResponse) Version() string {
	return res.response.Version
}

// Warning returns a notice from the Watchman server that, if non-empty,
// should be shown to the user as an advisory so that the system can
// operate more effectively
func (res *ListCapabilitiesResponse) Warning() string {
	return res.response.Warning
}

// Capabilities returns the result of the Watchman list-capabilities command.
func (res *ListCapabilitiesResponse) Capabilities() []string {
	return res.listCapabilitiesResponse.Capabilities
}
