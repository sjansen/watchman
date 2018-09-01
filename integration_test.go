// +build integration

package watchman_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sjansen/watchman"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	require := require.New(t)

	// connect
	c, err := watchman.Connect()
	require.NoError(err)

	// connection metadata
	require.NotEmpty(c.SockName())
	require.NotEmpty(c.Version())

	// watch-project
	wd, err := os.Getwd()
	require.NoError(err)

	testdata := filepath.Join(wd, "protocol", "testdata")
	testdata, err = filepath.EvalSymlinks(testdata)
	require.NoError(err)

	watch, err := c.WatchProject(testdata)
	require.NoError(err)

	// watch-list
	roots, err := c.WatchList()
	require.NoError(err)
	require.NotEmpty(roots)

	// clock
	require.NotEmpty(watch.Clock(0))

	// subscribe
	s, err := watch.Subscribe("Spoon!", testdata)
	require.NoError(err)

	// unsubscribe
	err = s.Unsubscribe()
	require.NoError(err)
}
