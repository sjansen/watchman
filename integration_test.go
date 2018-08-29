// +build integration

package watchman

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjansen/watchman/protocol"
)

const subName = "TANSTAAFL"

func TestSendAndRecv(t *testing.T) {
	require := require.New(t)

	// connect
	c, err := protocol.Connect()
	require.NoError(err)
	require.NotEmpty(c.Version())

	// watch-project
	wd, err := os.Getwd()
	require.NoError(err)

	testdata := filepath.Join(wd, "testdata")
	testdata, err = filepath.EvalSymlinks(testdata)
	require.NoError(err)

	err = c.Send(&protocol.WatchProjectRequest{testdata})
	require.NoError(err)

	watchProject := &protocol.WatchProjectResponse{}
	event, err := c.Recv(watchProject)
	require.NoError(err)
	require.Nil(event)
	watchRoot := watchProject.Watch()
	require.NotEmpty(watchRoot)

	// watch-list
	err = c.Send(&protocol.WatchListRequest{})
	require.NoError(err)

	watchList := &protocol.WatchListResponse{}
	event, err = c.Recv(watchList)
	require.NoError(err)
	require.Nil(event)
	require.NotEmpty(watchList.Roots())

	// clock
	for _, req := range []*protocol.ClockRequest{
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

		clock := &protocol.ClockResponse{}
		event, err = c.Recv(clock)
		require.NoError(err)
		require.Nil(event)
		require.NotEmpty(clock.Clock())
	}

	// subscribe
	err = c.Send(&protocol.SubscribeRequest{
		Root: testdata,
		Name: subName,
	})
	require.NoError(err)

	sub := &protocol.SubscribeResponse{}
	event, err = c.Recv(sub)
	require.NoError(err)
	require.Nil(event)
	require.NotEmpty(sub.Clock())
	require.Equal(subName, sub.Subscription())

	// unsubscribe
	err = c.Send(&protocol.UnsubscribeRequest{
		Root: testdata,
		Name: subName,
	})
	require.NoError(err)

	var unilaterals []protocol.Unilateral
	unsub := &protocol.UnsubscribeResponse{}
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
