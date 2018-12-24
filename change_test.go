package watchman

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	require := require.New(t)

	var c StateChange
	require.Equal("invalid", c.String())
	require.Equal("created", Created.String())
	require.Equal("removed", Removed.String())
	require.Equal("updated", Updated.String())
	require.Equal("ephemeral", Ephemeral.String())
}
