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
			request:  `["subscribe","/tmp","sub1",{"fields":["name"]}]` + "\n",
			response: `{"clock":"c:1531594843:978:9:345","subscribe":"sub1","version":"4.9.0"}` + "\n",
			req: &SubscribeRequest{
				Root: "/tmp",
				Name: "sub1",
			},
			res: &SubscribeResponse{
				subscribeResponse: subscribeResponse{
					response: response{
						Version: "4.9.0",
					},
					Clock:        "c:1531594843:978:9:345",
					Subscription: "sub1",
				}},
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

		actual := &SubscribeResponse{}
		event, err := c.Recv(actual)
		require.NoError(err)
		require.Nil(event)
		require.Equal(tc.res, actual)
		require.Equal("", actual.Warning())
		require.Equal("4.9.0", actual.Version())
		require.Equal("c:1531594843:978:9:345", actual.Clock())
		require.Equal("sub1", actual.Subscription())
	}
}
