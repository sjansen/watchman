// +build integration

package watchman_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sjansen/watchman"
	"github.com/stretchr/testify/require"
)

func mkdir() (dir string, err error) {
	dir, err = ioutil.TempDir("", "watchman-client-test")
	if err != nil {
		return
	}

	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		return
	}

	path := filepath.Join(dir, ".watchmanconfig")
	err = ioutil.WriteFile(path, []byte(`{"idle_reap_age_seconds": 300}`+"\n"), os.ModePerm)
	return
}

func touch(dir string, names ...string) error {
	for _, name := range names {
		path := filepath.Join(dir, name)
		err := ioutil.WriteFile(path, []byte("Kilroy was here."), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestClient(t *testing.T) {
	require := require.New(t)

	dir, err := mkdir()
	require.NoError(err)

	// connect
	c, err := watchman.Connect()
	require.NoError(err)

	// connection metadata
	require.NotEmpty(c.SockName())
	require.NotEmpty(c.Version())

	// watch-project
	watch, err := c.WatchProject(dir)
	require.NoError(err)

	// watch-list
	roots, err := c.WatchList()
	require.NoError(err)
	require.NotEmpty(roots)

	// subscribe
	s, err := watch.Subscribe("Spoon!", dir)
	require.NoError(err)

	// clock
	clock1, err := watch.Clock(0)
	require.NoError(err)
	require.NotEmpty(clock1)

	err = touch(dir, "foo", "bar", "baz")
	require.NoError(err)

	clock2, err := watch.Clock(3000)
	require.NoError(err)
	require.NotEmpty(clock2)
	require.NotEqual(clock1, clock2)

	// unsubscribe
	err = s.Unsubscribe()
	require.NoError(err)

	// close
	err = c.Close()
	require.NoError(err)
}
