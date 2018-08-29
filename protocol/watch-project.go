package protocol

/*
$ watchman watch-project ~/www/some/child/dir
{
  "version": "3.0.1",
  "watch": "/Users/wez/www",
  "relative_path": "some/child/dir"
}
*/

// A WatchProjectRequest represents the Watchman watch-project command.
//
// See also: https://facebook.github.io/watchman/docs/cmd/watch-project.html
type WatchProjectRequest struct {
	Path string
}

// Args returns values used to encode a request PDU.
func (req *WatchProjectRequest) Args() []interface{} {
	return []interface{}{"watch-project", req.Path}
}

// A WatchProjectResponse represents a response to the Watchman watch-project command.
type WatchProjectResponse struct {
	response
	relativePath string
	watch        string
}

// NewWatchProjectResponse converts a ResponsePDU to WatchProjectResponse
func NewWatchProjectResponse(pdu ResponsePDU) (res *WatchProjectResponse) {
	res = &WatchProjectResponse{}
	res.response.init(pdu)

	if x, ok := pdu["relative_path"]; ok {
		if relativePath, ok := x.(string); ok {
			res.relativePath = relativePath
		}
	}
	if x, ok := pdu["watch"]; ok {
		if watch, ok := x.(string); ok {
			res.watch = watch
		}
	}
	return
}

// RelativePath returns the difference between the requested directory
// and the watched directory actually chosen by the watch-project
// command.
func (res *WatchProjectResponse) RelativePath() string {
	return res.relativePath
}

// Watch returns the watched directory chosen by the watch-project
// command.
func (res *WatchProjectResponse) Watch() string {
	return res.watch
}
