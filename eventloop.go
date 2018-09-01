package watchman

import (
	"runtime"

	"github.com/sjansen/watchman/protocol"
)

type eventloop struct {
	requests    chan<- protocol.Request
	responses   <-chan protocol.ResponsePDU
	unilaterals <-chan protocol.ResponsePDU
}

func reader(conn *protocol.Connection) <-chan protocol.ResponsePDU {
	ch := make(chan protocol.ResponsePDU)
	go func() {
		defer close(ch)

		for {
			pdu, err := conn.Recv()
			if err != nil {
				return
			}
			ch <- pdu
		}
	}()
	return ch
}

func startEventLoop(conn *protocol.Connection) (l *eventloop, stop func(bool)) {
	/* SHUTDOWN
	requests:    closed by caller/stop()
	responses:   closed locally
	unilaterals: closed locally
	*/

	recv := reader(conn)
	requests := make(chan protocol.Request)
	responses := make(chan protocol.ResponsePDU)
	unilaterals := make(chan protocol.ResponsePDU)
	l = &eventloop{
		requests:    requests,
		responses:   responses,
		unilaterals: unilaterals,
	}

	expectRequest := func() (ok bool) {
		for {
			select {
			case req, ok := <-requests:
				if ok {
					err := conn.Send(req)
					ok = err == nil
				}
				return ok
			case pdu, ok := <-recv:
				if ok {
					unilaterals <- pdu
				} else {
					return false
				}
			}
		}
	}

	expectResponse := func() (ok bool) {
		for pdu := range recv {
			if pdu.IsUnilateral() {
				unilaterals <- pdu
			} else {
				responses <- pdu
				return true
			}
		}
		return false
	}

	// Close and empty channels so that other goroutines can shutdown.
	// IMPORTANT: delayClose should normally be false. It is useful in
	// tests that would be invalidated by closing `requests` too early.
	stop = func(delayClose bool) {
		if !delayClose {
			close(requests)
		}
		for range responses {
			continue
		}
		for range unilaterals {
			continue
		}
		if delayClose {
			close(requests)
		}
		// allow other goroutines to run their shutdown logic;
		// avoid false positives in tests to detect leaks
		runtime.Gosched()
	}

	go func() {
		defer func() {
			conn.Close()
			close(responses)
			close(unilaterals)
			for range requests {
				continue
			}
		}()
		for {
			if ok := expectRequest(); !ok {
				return
			}
			if ok := expectResponse(); !ok {
				return
			}
		}
	}()

	return
}
