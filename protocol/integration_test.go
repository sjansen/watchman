// +build integration

package protocol_test

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

	// sockname
	sockname := c.SockName()
	require.NotEmpty(sockname)

	// capabilities
	require.Equal(true, c.HasCapability("cmd-watch-project"))
	require.Equal(false, c.HasCapability("fire-the-missles"))

	// watch-project
	wd, err := os.Getwd()
	require.NoError(err)

	testdata := filepath.Join(wd, "testdata")
	testdata, err = filepath.EvalSymlinks(testdata)
	require.NoError(err)

	err = c.Send(&protocol.WatchProjectRequest{testdata})
	require.NoError(err)

	pdu, err := c.Recv()
	require.NoError(err)
	require.NotNil(pdu)
	watchProject := protocol.NewWatchProjectResponse(pdu)
	watchRoot := watchProject.Watch()
	require.NotEmpty(watchRoot)

	// watch-list
	err = c.Send(&protocol.WatchListRequest{})
	require.NoError(err)

	pdu, err = c.Recv()
	require.NoError(err)
	require.NotNil(pdu)
	watchList := protocol.NewWatchListResponse(pdu)
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

		pdu, err = c.Recv()
		require.NoError(err)
		require.NotNil(pdu)
		clock := protocol.NewClockResponse(pdu)
		require.NotEmpty(clock.Clock())
	}

	// subscribe
	err = c.Send(&protocol.SubscribeRequest{
		Root: testdata,
		Name: subName,
	})
	require.NoError(err)

	pdu, err = c.Recv()
	require.NoError(err)
	require.NotNil(pdu)
	sub := protocol.NewSubscribeResponse(pdu)
	require.NotEmpty(sub.Clock())
	require.Equal(subName, sub.Subscription())

	// unsubscribe
	err = c.Send(&protocol.UnsubscribeRequest{
		Root: testdata,
		Name: subName,
	})
	require.NoError(err)

	var unilaterals []protocol.ResponsePDU
	var unsub *protocol.UnsubscribeResponse
	for {
		pdu, err = c.Recv()
		require.NoError(err)
		if pdu.IsUnilateral() {
			unilaterals = append(unilaterals, pdu)
		} else {
			unsub = protocol.NewUnsubscribeResponse(pdu)
			break
		}
	}
	require.Equal(subName, unsub.Subscription())
	require.NotEmpty(unilaterals)

	pdu = unilaterals[0]
	require.Equal(pdu["subscription"], subName)
	require.Equal(pdu["root"], testdata)
	require.NotEmpty(pdu["clock"])
	require.NotEmpty(pdu["files"])

	err = c.Close()
	require.NoError(err)
}

type NoopRequest struct{}

func (req *NoopRequest) Args() []interface{} {
	return []interface{}{"noop"}
}

type VersionRequest struct{}

func (req *VersionRequest) Args() []interface{} {
	return []interface{}{"version"}
}

func TestInvalidCommand(t *testing.T) {
	require := require.New(t)

	c, err := protocol.Connect()
	require.NoError(err)
	require.NotEmpty(c.Version())

	// invalid command should result in WatchmanError
	err = c.Send(&NoopRequest{})
	require.NoError(err)

	pdu, err := c.Recv()
	require.Nil(pdu)
	require.NotNil(err)
	require.IsType(&protocol.WatchmanError{}, err)
	require.NotEmpty(err.Error())

	// connection should still be valid
	err = c.Send(&VersionRequest{})
	require.NoError(err)

	pdu, err = c.Recv()
	require.NotNil(pdu)
	require.NoError(err)

	err = c.Close()
	require.NoError(err)
}
