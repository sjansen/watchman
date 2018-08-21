package connection

/*
$ watchman clock /path/to/dir
{"clock":"c:1531594843:978:9:172","version":"4.9.0"}

["clock", "/path/to/dir", {"sync_timeout": 100}]
{"error":"sync_timeout expired","version":"4.9.0"}
*/

type ClockRequest struct {
	Path        string
	SyncTimeout int
}

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

type ClockResponse struct {
	clockResponse
}

func (res *ClockResponse) Version() string {
	return res.response.Version
}

func (res *ClockResponse) Warning() string {
	return res.response.Warning
}

func (res *ClockResponse) Clock() string {
	return res.clockResponse.Clock
}
