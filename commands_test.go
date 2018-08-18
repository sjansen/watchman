package watchman

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	require := require.New(t)

	for _, tc := range []struct {
		expected []byte
		command  interface{}
	}{
		{
			expected: []byte(`["get-config", "/tmp"]` + "\n"),
			command: &GetConfig{
				Path: "/tmp",
			},
		},
		{
			expected: []byte(`["watch-list"]` + "\n"),
			command:  &WatchList{},
		},
		{
			expected: []byte(`["watch-project", "/tmp"]` + "\n"),
			command: &WatchProject{
				Path: "/tmp",
			},
		},
	} {
		actual, err := Marshal(tc.command)
		require.NoError(err)
		require.Equal(tc.expected, actual)
	}
}
