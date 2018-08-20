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

	// connect
	c, err := Connect()
	require.NoError(err)
	require.NotEmpty(c.Version())

	// watch-project
	wd, err := os.Getwd()
	require.NoError(err)

	testdata := filepath.Join(wd, "testdata")
	err = c.Send(&WatchProjectRequest{testdata})
	require.NoError(err)

	watchProject := &WatchProjectResponse{}
	event, err := c.Recv(watchProject)
	require.NoError(err)
	require.Nil(event)
	watchRoot := watchProject.Watch()
	require.NotEmpty(watchRoot)

	// watch-list
	err = c.Send(&WatchListRequest{})
	require.NoError(err)

	watchList := &WatchListResponse{}
	event, err = c.Recv(watchList)
	require.NoError(err)
	require.Nil(event)
	require.NotEmpty(watchList.Roots())

	// clock
	for _, req := range []*ClockRequest{
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

		clock := &ClockResponse{}
		event, err = c.Recv(clock)
		require.NoError(err)
		require.Nil(event)
		require.NotEmpty(clock.Clock())
	}

	// subscribe
	err = c.Send(&SubscribeRequest{
		Root: testdata,
		Name: "sub:testdata",
	})
	require.NoError(err)

	sub := &SubscribeResponse{}
	event, err = c.Recv(sub)
	require.NoError(err)
	require.Nil(event)
	require.NotEmpty(sub.Clock())
	require.Equal("sub:testdata", sub.Subscription())

	// unsubscribe
	err = c.Send(&UnsubscribeRequest{
		Root: testdata,
		Name: "sub:testdata",
	})
	require.NoError(err)

	var events []*Subscription
	unsub := &UnsubscribeResponse{}
	for {
		event, err = c.Recv(unsub)
		require.NoError(err)
		if event == nil {
			break
		}
		events = append(events, event)
	}
	require.Equal("sub:testdata", unsub.Subscription())
	require.NotEmpty(events)
}
