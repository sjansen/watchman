package protocol

/*
$ watchman clock /path/to/dir
{"clock":"c:1531594843:978:9:172","version":"4.9.0"}

["clock", "/path/to/dir", {"sync_timeout": 100}]
{"error":"sync_timeout expired","version":"4.9.0"}
*/

// A ClockRequest represents the Watchman clock command.
//
// See also: https://facebook.github.io/watchman/docs/cmd/clock.html
type ClockRequest struct {
	Path        string
	SyncTimeout int
}

// Args returns values used to encode a request PDU.
func (req *ClockRequest) Args() []interface{} {
	if req.SyncTimeout < 1 {
		return []interface{}{"clock", req.Path}
	}
	m := map[string]int{"sync_timeout": req.SyncTimeout}
	return []interface{}{"clock", req.Path, m}
}

// A ClockResponse represents a response to the Watchman clock command.
type ClockResponse struct {
	response
	clock string
}

// NewClockResponse converts a ResponsePDU to ClockResponse
func NewClockResponse(pdu ResponsePDU) (res *ClockResponse) {
	res = &ClockResponse{}
	res.response.init(pdu)

	if x, ok := pdu["clock"]; ok {
		if clock, ok := x.(string); ok {
			res.clock = clock
		}
	}
	return
}

// Clock returns the result of the Watchman clock command.
func (res *ClockResponse) Clock() string {
	return res.clock
}
