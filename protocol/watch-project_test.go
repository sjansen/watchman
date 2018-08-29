package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWatchProject(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		request  string
		response string
		req      *WatchProjectRequest
		res      *WatchProjectResponse
		rpath    string
	}{
		{
			request:  `["watch-project","/tmp"]` + "\n",
			response: `{"watcher":"fsevents","watch":"/tmp","version":"4.9.0"}` + "\n",
			req:      &WatchProjectRequest{"/tmp"},
			res: &WatchProjectResponse{
				response: response{
					pdu: ResponsePDU{
						"version": "4.9.0",
						"watch":   "/tmp",
						"watcher": "fsevents",
					},
					version: "4.9.0",
				},
				watch: "/tmp",
			},
			rpath: "",
		},
		{
			request:  `["watch-project","/tmp/testdata"]` + "\n",
			response: `{"watch":"/tmp","relative_path":"testdata","version":"4.9.0"}` + "\n",
			req:      &WatchProjectRequest{"/tmp/testdata"},
			res: &WatchProjectResponse{
				response: response{
					pdu: ResponsePDU{
						"relative_path": "testdata",
						"version":       "4.9.0",
						"watch":         "/tmp",
					},
					version: "4.9.0",
				},
				watch:        "/tmp",
				relativePath: "testdata",
			},
			rpath: "testdata",
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
		actual := NewWatchProjectResponse(pdu)
		require.Equal(tc.res, actual)
		require.Equal("", actual.Warning())
		require.Equal("4.9.0", actual.Version())
		require.Equal("/tmp", actual.Watch())
		require.Equal(tc.rpath, actual.RelativePath())
	}
}
