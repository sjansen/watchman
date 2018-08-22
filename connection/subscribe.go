package connection

/*
["subscribe","/tmp","sub1",{"fields":["name"]}]
{"clock":"c:1531594843:978:9:826","subscribe":"sub1","version":"4.9.0"}
{"unilateral":true,"subscription":"sub1","root":"/tmp","files":["foo.go","bar.go"],"version":"4.9.0","clock":"c:1531594843:978:9:826","is_fresh_instance":true}
{"unilateral":true,"subscription":"sub1","root":"/tmp","files":["foo.go"],"version":"4.9.0","since":"c:1531594843:978:9:826","clock":"c:1531594843:978:9:827","is_fresh_instance":false}
{"unilateral":true,"subscription":"sub1","root":"/tmp","canceled":true,"version":"4.9.0"}
*/

// A Subscription represents changes observed as a result of the Watchman subscribe command.
type Subscription struct {
	unilateral
	Clock           string
	Root            string
	Subscription    string
	IsFreshInstance bool `json:"is_fresh_instance"`
	Files           []string
}

// A SubscribeRequest represents the Watchman subscribe command.
//
// See also: https://facebook.github.io/watchman/docs/cmd/subscribe.html
type SubscribeRequest struct {
	Root string
	Name string
}

// Args returns values used to encode a request PDU.
func (req *SubscribeRequest) Args() []interface{} {
	m := map[string]interface{}{"fields": []string{"name"}}
	return []interface{}{"subscribe", req.Root, req.Name, m}
}

type subscribeResponse struct {
	response
	Clock        string
	Subscription string `json:"subscribe"`
}

// A SubscribeResponse represents a response to the Watchman subscribe command.
type SubscribeResponse struct {
	subscribeResponse
}

// Version returns the Watchman server version.
func (res *SubscribeResponse) Version() string {
	return res.response.Version
}

// Warning returns a notice from the Watchman server that, if non-empty,
// should be shown to the user as an advisory so that the system can
// operate more effectively
func (res *SubscribeResponse) Warning() string {
	return res.response.Warning
}

// Clock returns a value represention when the subscription started.
func (res *SubscribeResponse) Clock() string {
	return res.subscribeResponse.Clock
}

// Subscription returns the name registered to the subscription.
func (res *SubscribeResponse) Subscription() string {
	return res.subscribeResponse.Subscription
}
