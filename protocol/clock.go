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

type clockResponse struct {
	response
	Clock string
}

// A ClockResponse represents a response to the Watchman clock command.
type ClockResponse struct {
	clockResponse
}

// Version returns the Watchman server version.
func (res *ClockResponse) Version() string {
	return res.response.Version
}

// Warning returns a notice from the Watchman server that, if non-empty,
// should be shown to the user as an advisory so that the system can
// operate more effectively
func (res *ClockResponse) Warning() string {
	return res.response.Warning
}

// Clock returns the result of the Watchman clock command.
func (res *ClockResponse) Clock() string {
	return res.clockResponse.Clock
}
