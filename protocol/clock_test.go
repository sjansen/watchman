package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClock(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		request  string
		response string
		req      *ClockRequest
		res      *ClockResponse
	}{
		{
			request:  `["clock","/tmp"]` + "\n",
			response: `{"clock":"c:1531594843:978:9:345","version":"4.9.0"}` + "\n",
			req:      &ClockRequest{Path: "/tmp"},
			res: &ClockResponse{
				response: response{
					pdu: ResponsePDU{
						"version": "4.9.0",
						"clock":   "c:1531594843:978:9:345",
					},
					version: "4.9.0",
				},
				clock: "c:1531594843:978:9:345",
			},
		},
		{
			request:  `["clock","/tmp",{"sync_timeout":1234}]` + "\n",
			response: `{"clock":"c:1531594843:978:9:345","version":"4.9.0"}` + "\n",
			req: &ClockRequest{
				Path:        "/tmp",
				SyncTimeout: 1234,
			},
			res: &ClockResponse{
				response: response{
					pdu: ResponsePDU{
						"version": "4.9.0",
						"clock":   "c:1531594843:978:9:345",
					},
					version: "4.9.0",
				},
				clock: "c:1531594843:978:9:345",
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
		actual := NewClockResponse(pdu)
		require.Equal(tc.res, actual)
		require.Equal("", actual.Warning())
		require.Equal("4.9.0", actual.Version())
		require.Equal("c:1531594843:978:9:345", actual.Clock())
	}
}
