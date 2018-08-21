package connection

/*
$ watchman watch-list
{
  "version": "1.9",
  "roots": [
    "/home/wez/watchman"
  ]
}
*/

type WatchListRequest struct{}

func (req *WatchListRequest) Args() []interface{} {
	return []interface{}{"watch-list"}
}

type watchListResponse struct {
	response
	Roots []string
}

type WatchListResponse struct {
	watchListResponse
}

func (res *WatchListResponse) Version() string {
	return res.response.Version
}

func (res *WatchListResponse) Warning() string {
	return res.response.Warning
}

func (res *WatchListResponse) Roots() []string {
	return res.watchListResponse.Roots
}
