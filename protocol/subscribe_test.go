package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubscribe(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		request  string
		response string
		req      *SubscribeRequest
		res      *SubscribeResponse
	}{
		{
			request: `["subscribe","/tmp","sub1",{"fields":[` +
				`"cclock","ctime","exists","gid","mode","mtime","name",` +
				`"nlink","oclock","size","symlink_target","type","uid"` +
				"]}]\n",
			response: `{"clock":"c:1531594843:978:9:345","subscribe":"sub1","version":"4.9.0"}` + "\n",
			req: &SubscribeRequest{
				Root: "/tmp",
				Name: "sub1",
			},
			res: &SubscribeResponse{
				response: response{
					pdu: ResponsePDU{
						"version":   "4.9.0",
						"clock":     "c:1531594843:978:9:345",
						"subscribe": "sub1",
					},
					version: "4.9.0",
				},
				clock:        "c:1531594843:978:9:345",
				subscription: "sub1",
			},
		},
	} {
		requested := &bytes.Buffer{}
		c := &Connection{
			reader: bufio.NewReader(
				bytes.NewReader([]byte(tc.response)),
			),
			socket: requested,
		}

		err := c.Send(tc.req)
		require.NoError(err)
		require.Equal(tc.request, requested.String())

		pdu, err := c.Recv()
		require.NoError(err)
		require.NotNil(pdu)
		actual := NewSubscribeResponse(pdu)
		require.Equal(tc.res, actual)
		require.Equal("", actual.Warning())
		require.Equal("4.9.0", actual.Version())
		require.Equal("c:1531594843:978:9:345", actual.Clock())
		require.Equal("sub1", actual.Subscription())
	}
}

func TestNewSubscription(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		pdu ResponsePDU
		sub *Subscription
	}{
		{
			pdu: ResponsePDU{
				"unilateral":        true,
				"subscription":      "sub2",
				"root":              "/tmp",
				"version":           "4.9.0",
				"clock":             "c:1531594843:978:9:826",
				"is_fresh_instance": true,
				"files": []interface{}{
					map[string]interface{}{"name": "foo/main.go", "exists": true},
					map[string]interface{}{"name": "bar/main.go", "exists": true},
				},
			},
			sub: &Subscription{
				response: response{
					pdu: ResponsePDU{
						"unilateral":        true,
						"subscription":      "sub2",
						"root":              "/tmp",
						"version":           "4.9.0",
						"clock":             "c:1531594843:978:9:826",
						"is_fresh_instance": true,
						"files": []interface{}{
							map[string]interface{}{
								"name": "foo/main.go", "exists": true,
							},
							map[string]interface{}{
								"name": "bar/main.go", "exists": true,
							},
						},
					},
					version: "4.9.0",
				},
				clock:           "c:1531594843:978:9:826",
				root:            "/tmp",
				subscription:    "sub2",
				isFreshInstance: true,
				files: []map[string]interface{}{
					{"name": "foo/main.go", "exists": true},
					{"name": "bar/main.go", "exists": true},
				},
			},
		},
	} {
		actual := NewSubscription(tc.pdu)
		require.Equal(tc.sub, actual)
	}
}

func TestSubscription(t *testing.T) {
	require := require.New(t)

	s := &Subscription{
		clock:        "c:2642605954:867:8:937",
		root:         "/projects/x",
		subscription: "sub42",
		files: []map[string]interface{}{
			{"name": "secrets.txt", "exists": true},
		},
		isFreshInstance: true,
	}
	require.Equal("c:2642605954:867:8:937", s.Clock())
	require.Equal(true, s.IsFreshInstance())
	require.Equal("/projects/x", s.Root())
	require.Equal("sub42", s.Subscription())
	require.Equal([]map[string]interface{}{
		{"name": "secrets.txt", "exists": true},
	}, s.Files())
}
