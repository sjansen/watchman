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

// A WatchListResponse represents a response to the Watchman watch-list command.
type WatchListResponse struct {
	response
	roots []string
}

// NewWatchListResponse converts a ResponsePDU to WatchListResponse
func NewWatchListResponse(pdu ResponsePDU) (res *WatchListResponse) {
	res = &WatchListResponse{}
	res.response.init(pdu)

	if x, ok := pdu["roots"]; ok {
		if roots, ok := x.([]interface{}); ok {
			res.roots = make([]string, len(roots))
			for i, root := range roots {
				res.roots[i] = root.(string)
			}
		}
	}
	return
}

// Roots returns the result of the Watchman watch-list command.
func (res *WatchListResponse) Roots() []string {
	return res.roots
}
