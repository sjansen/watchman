package connection

/*
["unsubscribe", "/tmp", "sub1"]
{"unsubscribe":"sub1", "deleted":true, "version":"4.9.0"}
*/

type UnsubscribeRequest struct {
	Root string
	Name string
}

func (req *UnsubscribeRequest) Args() []interface{} {
	return []interface{}{"unsubscribe", req.Root, req.Name}
}

type unsubscribeResponse struct {
	response
	Deleted      bool
	Subscription string `json:"unsubscribe"`
}

type UnsubscribeResponse struct {
	unsubscribeResponse
}

func (res *UnsubscribeResponse) Version() string {
	return res.response.Version
}

func (res *UnsubscribeResponse) Warning() string {
	return res.response.Warning
}

func (res *UnsubscribeResponse) Deleted() bool {
	return res.unsubscribeResponse.Deleted
}

func (res *UnsubscribeResponse) Subscription() string {
	return res.unsubscribeResponse.Subscription
}
