package protocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

type FakeResponse struct {
	response
}

func NewFakeResponse(pdu ResponsePDU) (res *FakeResponse) {
	res = &FakeResponse{}
	res.response.init(pdu)
	return
}

func TestResponseMixin(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		response string
		expected *FakeResponse
	}{
		{
			response: `{"version":"4.9.0"}` + "\n",
			expected: &FakeResponse{
				response: response{
					pdu: ResponsePDU{
						"version": "4.9.0",
					},
					version: "4.9.0",
				},
			},
		}, {
			response: `{"version":"4.9.0","warning":"restart pending"}` + "\n",
			expected: &FakeResponse{
				response: response{
					pdu: ResponsePDU{
						"version": "4.9.0",
						"warning": "restart pending",
					},
					version: "4.9.0",
					warning: "restart pending",
				},
			},
		},
	} {
		c := &Connection{
			reader: bufio.NewReader(
				bytes.NewReader([]byte(tc.response)),
			),
		}

		pdu, err := c.Recv()
		require.NoError(err)
		require.NotNil(pdu)
		actual := NewFakeResponse(pdu)
		require.Equal(tc.expected, actual)
		require.Equal(tc.expected.response.pdu, actual.PDU())
		require.Equal(tc.expected.response.pdu["version"], actual.Version())
		if warning, ok := tc.expected.response.pdu["warning"]; ok {
			require.Equal(warning, actual.Warning())
		}
	}
}
