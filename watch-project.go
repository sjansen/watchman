package watchman

/*
$ watchman watch-project ~/www/some/child/dir
{
  "version": "3.0.1",
  "watch": "/Users/wez/www",
  "relative_path": "some/child/dir"
}
*/

type WatchProjectRequest struct {
	Path string
}

func (req *WatchProjectRequest) Args() []interface{} {
	return []interface{}{"watch-project", req.Path}
}

type watchProjectResponse struct {
	response
	RelativePath string `json:"relative_path"`
	Watch        string
}

type WatchProjectResponse struct {
	watchProjectResponse
}

func (res *WatchProjectResponse) Version() string {
	return res.response.Version
}

func (res *WatchProjectResponse) Warning() string {
	return res.response.Warning
}

func (res *WatchProjectResponse) RelativePath() string {
	return res.watchProjectResponse.RelativePath
}

func (res *WatchProjectResponse) Watch() string {
	return res.watchProjectResponse.Watch
}
