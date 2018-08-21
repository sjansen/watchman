// +build integration

package watchman

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	conn "github.com/sjansen/watchman/connection"
)

const subName = "TANSTAAFL"

func TestSendAndRecv(t *testing.T) {
	require := require.New(t)

	// connect
	c, err := conn.New()
	require.NoError(err)
	require.NotEmpty(c.Version())

	// watch-project
	wd, err := os.Getwd()
	require.NoError(err)

	testdata := filepath.Join(wd, "testdata")
	err = c.Send(&conn.WatchProjectRequest{testdata})
	require.NoError(err)

	watchProject := &conn.WatchProjectResponse{}
	event, err := c.Recv(watchProject)
	require.NoError(err)
	require.Nil(event)
	watchRoot := watchProject.Watch()
	require.NotEmpty(watchRoot)

	// watch-list
	err = c.Send(&conn.WatchListRequest{})
	require.NoError(err)

	watchList := &conn.WatchListResponse{}
	event, err = c.Recv(watchList)
	require.NoError(err)
	require.Nil(event)
	require.NotEmpty(watchList.Roots())

	// clock
	for _, req := range []*conn.ClockRequest{
		{
			Path: watchRoot,
		},
		{
			Path:        watchRoot,
			SyncTimeout: 1000,
		},
	} {
		err = c.Send(req)
		require.NoError(err)

		clock := &conn.ClockResponse{}
		event, err = c.Recv(clock)
		require.NoError(err)
		require.Nil(event)
		require.NotEmpty(clock.Clock())
	}

	// subscribe
	err = c.Send(&conn.SubscribeRequest{
		Root: testdata,
		Name: subName,
	})
	require.NoError(err)

	sub := &conn.SubscribeResponse{}
	event, err = c.Recv(sub)
	require.NoError(err)
	require.Nil(event)
	require.NotEmpty(sub.Clock())
	require.Equal(subName, sub.Subscription())

	// unsubscribe
	err = c.Send(&conn.UnsubscribeRequest{
		Root: testdata,
		Name: subName,
	})
	require.NoError(err)

	var unilaterals []conn.Unilateral
	unsub := &conn.UnsubscribeResponse{}
	for {
		event, err = c.Recv(unsub)
		require.NoError(err)
		if event == nil {
			break
		}
		unilaterals = append(unilaterals, event)
	}
	require.Equal(subName, unsub.Subscription())
	require.NotEmpty(unilaterals)

	pdu := unilaterals[0].PDU()
	require.Equal(pdu["subscription"], subName)
	require.Equal(pdu["root"], testdata)
	require.NotEmpty(pdu["clock"])
	require.NotEmpty(pdu["files"])
}
