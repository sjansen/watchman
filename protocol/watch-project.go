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

type watchProjectResponse struct {
	response
	RelativePath string `json:"relative_path"`
	Watch        string
}

// A WatchProjectResponse represents a response to the Watchman watch-project command.
type WatchProjectResponse struct {
	watchProjectResponse
}

// Version returns the Watchman server version.
func (res *WatchProjectResponse) Version() string {
	return res.response.Version
}

// Warning returns a notice from the Watchman server that, if non-empty,
// should be shown to the user as an advisory so that the system can
// operate more effectively
func (res *WatchProjectResponse) Warning() string {
	return res.response.Warning
}

// RelativePath returns the difference between the requested directory
// and the watched directory actually chosen by the watch-project
// command.
func (res *WatchProjectResponse) RelativePath() string {
	return res.watchProjectResponse.RelativePath
}

// Watch returns the watched directory chosen by the watch-project
// command.
func (res *WatchProjectResponse) Watch() string {
	return res.watchProjectResponse.Watch
}
