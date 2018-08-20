package watchman

/*
["subscribe","/tmp","sub1",{"fields":["name"]}]
{"clock":"c:1531594843:978:9:826","subscribe":"sub1","version":"4.9.0"}
{"unilateral":true,"subscription":"sub1","root":"/tmp","files":["foo.go","bar.go"],"version":"4.9.0","clock":"c:1531594843:978:9:826","is_fresh_instance":true}
{"unilateral":true,"subscription":"sub1","root":"/tmp","files":["foo.go"],"version":"4.9.0","since":"c:1531594843:978:9:826","clock":"c:1531594843:978:9:827","is_fresh_instance":false}
*/

type Subscription struct {
	Clock           string
	Root            string
	Subscription    string
	IsFreshInstance bool `json:"is_fresh_instance"`
	Files           []string
}

type SubscribeRequest struct {
	Root string
	Name string
}

func (req *SubscribeRequest) Args() []interface{} {
	m := map[string]interface{}{"fields": []string{"name"}}
	return []interface{}{"subscribe", req.Root, req.Name, m}
}

type subscribeResponse struct {
	response
	Clock        string
	Subscription string `json:"subscribe"`
}

type SubscribeResponse struct {
	subscribeResponse
}

func (res *SubscribeResponse) Version() string {
	return res.response.Version
}

func (res *SubscribeResponse) Warning() string {
	return res.response.Warning
}

func (res *SubscribeResponse) Clock() string {
	return res.subscribeResponse.Clock
}

func (res *SubscribeResponse) Subscription() string {
	return res.subscribeResponse.Subscription
}
