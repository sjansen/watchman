package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWatchList(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		request  string
		response string
		req      *WatchListRequest
		res      *WatchListResponse
	}{
		{
			request:  `["watch-list"]` + "\n",
			response: `{"roots":["/tmp"],"version":"4.9.0"}` + "\n",
			req:      &WatchListRequest{},
			res: &WatchListResponse{
				response: response{
					pdu: ResponsePDU{
						"version": "4.9.0",
						"roots":   []interface{}{"/tmp"},
					},
					version: "4.9.0",
				},
				roots: []string{"/tmp"},
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
		actual := NewWatchListResponse(pdu)
		require.Equal(tc.res, actual)
		require.Equal("", actual.Warning())
		require.Equal("4.9.0", actual.Version())
		require.Equal([]string{"/tmp"}, actual.Roots())
	}
}
