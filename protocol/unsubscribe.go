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

// An UnsubscribeResponse represents a response to the Watchman unsubscribe command.
type UnsubscribeResponse struct {
	response
	deleted      bool
	subscription string
}

// NewUnsubscribeResponse converts a ResponsePDU to UnsubscribeResponse
func NewUnsubscribeResponse(pdu ResponsePDU) (res *UnsubscribeResponse) {
	res = &UnsubscribeResponse{}
	res.response.init(pdu)

	if x, ok := pdu["deleted"]; ok {
		if deleted, ok := x.(bool); ok {
			res.deleted = deleted
		}
	}
	if x, ok := pdu["unsubscribe"]; ok {
		if subscription, ok := x.(string); ok {
			res.subscription = subscription
		}
	}
	return
}

// Deleted returns the status of the cancelled subscription.
func (res *UnsubscribeResponse) Deleted() bool {
	return res.deleted
}

// Subscription returns the name registered to the cancelled subscription.
func (res *UnsubscribeResponse) Subscription() string {
	return res.subscription
}
