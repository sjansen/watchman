package watchman

import (
	"runtime"

	"github.com/sjansen/watchman/protocol"
)

type eventloop struct {
	requests  chan<- protocol.Request
	responses <-chan result
	updates   <-chan interface{}
}

type result struct {
	err *protocol.WatchmanError
	pdu protocol.ResponsePDU
}

func reader(conn *protocol.Connection) <-chan result {
	ch := make(chan result)
	go func() {
		defer close(ch)

		for {
			pdu, err := conn.Recv()
			result := result{}
			if err == nil {
				result.pdu = pdu
			} else if e, ok := err.(*protocol.WatchmanError); ok {
				result.err = e
			} else {
				return
			}
			ch <- result
		}
	}()
	return ch
}

func startEventLoop(conn *protocol.Connection) (l *eventloop, stop func(bool)) {
	/* SHUTDOWN
	requests:    closed by caller/stop()
	responses:   closed locally
	updates:     closed locally
	*/

	recv := reader(conn)
	requests := make(chan protocol.Request)
	responses := make(chan result)
	updates := make(chan interface{})
	l = &eventloop{
		requests:  requests,
		responses: responses,
		updates:   updates,
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
			case result, ok := <-recv:
				if ok {
					updates <- translateUnilateralPDU(result.pdu)
				} else {
					return false
				}
			}
		}
	}

	expectResponse := func() (ok bool) {
		for result := range recv {
			if result.err == nil && result.pdu.IsUnilateral() {
				updates <- translateUnilateralPDU(result.pdu)
			} else {
				responses <- result
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
		for range updates {
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
			close(updates)
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

func translateUnilateralPDU(pdu protocol.ResponsePDU) interface{} {
	if _, ok := pdu["subscription"]; ok {
		sub := protocol.NewSubscription(pdu)
		return newChangeNotification(sub)
	}
	return pdu
}
