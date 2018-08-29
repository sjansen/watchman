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

// A ListCapabilitiesResponse represents a response to the Watchman list-capabilities command.
type ListCapabilitiesResponse struct {
	response
	capabilities []string
}

// NewListCapabilitiesResponse converts a ResponsePDU to ListCapabilitiesResponse
func NewListCapabilitiesResponse(pdu ResponsePDU) (res *ListCapabilitiesResponse) {
	res = &ListCapabilitiesResponse{}
	res.response.init(pdu)

	if x, ok := pdu["capabilities"]; ok {
		if capabilities, ok := x.([]interface{}); ok {
			res.capabilities = make([]string, len(capabilities))
			for i, capability := range capabilities {
				res.capabilities[i] = capability.(string)
			}
		}
	}
	return
}

// Capabilities returns the result of the Watchman list-capabilities command.
func (res *ListCapabilitiesResponse) Capabilities() []string {
	return res.capabilities
}
