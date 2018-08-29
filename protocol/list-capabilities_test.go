package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListCapabilities(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		request  string
		response string
		req      *ListCapabilitiesRequest
		res      *ListCapabilitiesResponse
	}{
		{
			request:  `["list-capabilities"]` + "\n",
			response: `{"capabilities":["cmd-watch-project","cmd-subscribe"],"version":"4.9.0"}` + "\n",
			req:      &ListCapabilitiesRequest{},
			res: &ListCapabilitiesResponse{
				response: response{
					pdu: ResponsePDU{
						"version":      "4.9.0",
						"capabilities": []interface{}{"cmd-watch-project", "cmd-subscribe"},
					},
					version: "4.9.0",
				},
				capabilities: []string{
					"cmd-watch-project",
					"cmd-subscribe",
				},
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
		actual := NewListCapabilitiesResponse(pdu)
		require.Equal(tc.res, actual)
		require.Equal("", actual.Warning())
		require.Equal("4.9.0", actual.Version())
		require.Equal(
			[]string{"cmd-watch-project", "cmd-subscribe"},
			actual.Capabilities(),
		)
	}
}
