package protocol

/*
["unsubscribe", "/tmp", "sub1"]
{"unsubscribe":"sub1", "deleted":true, "version":"4.9.0"}
*/

// An UnsubscribeRequest represents the Watchman unsubscribe command.
//
// See also: https://facebook.github.io/watchman/docs/cmd/unsubscribe.html
type UnsubscribeRequest struct {
	Root string
	Name string
}

// Args returns values used to encode a request PDU.
func (req *UnsubscribeRequest) Args() []interface{} {
	return []interface{}{"unsubscribe", req.Root, req.Name}
}

type unsubscribeResponse struct {
	response
	Deleted      bool
	Subscription string `json:"unsubscribe"`
}

// An UnsubscribeResponse represents a response to the Watchman unsubscribe command.
type UnsubscribeResponse struct {
	unsubscribeResponse
}

// Version returns the Watchman server version.
func (res *UnsubscribeResponse) Version() string {
	return res.response.Version
}

// Warning returns a notice from the Watchman server that, if non-empty,
// should be shown to the user as an advisory so that the system can
// operate more effectively
func (res *UnsubscribeResponse) Warning() string {
	return res.response.Warning
}

// Deleted returns the status of the cancelled subscription.
func (res *UnsubscribeResponse) Deleted() bool {
	return res.unsubscribeResponse.Deleted
}

// Subscription returns the name registered to the cancelled subscription.
func (res *UnsubscribeResponse) Subscription() string {
	return res.unsubscribeResponse.Subscription
}
