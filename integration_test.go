// +build integration

package watchman

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendAndRecv(t *testing.T) {
	require := require.New(t)

	c, err := Connect()
	require.NoError(err)
	require.NotEmpty(c.Version())

	wd, err := os.Getwd()
	require.NoError(err)

	testdata := filepath.Join(wd, "testdata")
	err = c.Send(&WatchProjectRequest{testdata})
	require.NoError(err)

	watchProject := &WatchProjectResponse{}
	res, err := c.Recv(watchProject)
	require.NoError(err)
	require.Nil(res)
	require.NotEmpty(watchProject.Watch())

	err = c.Send(&WatchListRequest{})
	require.NoError(err)

	watchList := &WatchListResponse{}
	res, err = c.Recv(watchList)
	require.NoError(err)
	require.Nil(res)
	require.NotEmpty(watchList.Roots())
}
