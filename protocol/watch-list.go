package protocol

/*
$ watchman watch-list
{
  "version": "1.9",
  "roots": [
    "/home/wez/watchman"
  ]
}
*/

// A WatchListRequest represents the Watchman watch-list command.
//
// See also: https://facebook.github.io/watchman/docs/cmd/watch-list.html
type WatchListRequest struct{}

// Args returns values used to encode a request PDU.
func (req *WatchListRequest) Args() []interface{} {
	return []interface{}{"watch-list"}
}

type watchListResponse struct {
	response
	Roots []string
}

// A WatchListResponse represents a response to the Watchman watch-list command.
type WatchListResponse struct {
	watchListResponse
}

// Version returns the Watchman server version.
func (res *WatchListResponse) Version() string {
	return res.response.Version
}

// Warning returns a notice from the Watchman server that, if non-empty,
// should be shown to the user as an advisory so that the system can
// operate more effectively
func (res *WatchListResponse) Warning() string {
	return res.response.Warning
}

// Roots returns the result of the Watchman watch-list command.
func (res *WatchListResponse) Roots() []string {
	return res.watchListResponse.Roots
}
