package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnsubscribe(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		request  string
		response string
		req      *UnsubscribeRequest
		res      *UnsubscribeResponse
	}{
		{
			request:  `["unsubscribe","/tmp","sub1"]` + "\n",
			response: `{"unsubscribe":"sub1", "deleted":true, "version":"4.9.0"}` + "\n",
			req: &UnsubscribeRequest{
				Root: "/tmp",
				Name: "sub1",
			},
			res: &UnsubscribeResponse{
				unsubscribeResponse: unsubscribeResponse{
					response: response{
						Version: "4.9.0",
					},
					Deleted:      true,
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

		actual := &UnsubscribeResponse{}
		event, err := c.Recv(actual)
		require.NoError(err)
		require.Nil(event)
		require.Equal(tc.res, actual)
		require.Equal("", actual.Warning())
		require.Equal("4.9.0", actual.Version())
		require.Equal(true, actual.Deleted())
		require.Equal("sub1", actual.Subscription())
	}
}
