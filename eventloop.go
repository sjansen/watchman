package watchman

import (
	"encoding/json"
	"runtime"
)

type object map[string]interface{}

type eventloop struct {
	commands    chan<- string
	results     <-chan object
	unilaterals <-chan object
}

func loop(s *server) (l *eventloop, stop func(bool)) {
	/* SHUTDOWN
	commands:    closed by caller/stop()
	results:     closed locally
	unilaterals: closed locally
	s.commands:  closed locally
	s.events:    closed by *server
	*/

	commands := make(chan string)
	results := make(chan object)
	unilaterals := make(chan object)
	l = &eventloop{
		commands:    commands,
		results:     results,
		unilaterals: unilaterals,
	}

	expectCommand := func() (ok bool) {
		for {
			select {
			case command, ok := <-commands:
				if ok {
					s.commands <- command
				}
				return ok
			case pdu, ok := <-s.events:
				// TODO log warnings
				if ok {
					var event object
					if err := json.Unmarshal(pdu, &event); err != nil {
						ok = false
						event = object{"error": err.Error()}
					}
					unilaterals <- event
				}
				return ok
			}
		}
	}

	expectResult := func() (ok bool) {
		for {
			select {
			case pdu, ok := <-s.events:
				// TODO log warnings
				if ok {
					var event object
					if err := json.Unmarshal(pdu, &event); err != nil {
						ok = false
						event = object{"error": err.Error()}
					}
					if _, u8l := event["log"]; u8l {
						unilaterals <- event
					} else if _, u8l := event["subscription"]; u8l {
						unilaterals <- event
					} else {
						results <- event
					}
				}
				return ok
			}
		}
	}

	// Close and empty channels so that other goroutines can shutdown.
	// IMPORTANT: delayClose should normally be false. It is useful in
	// tests that would be invalidated by close l.commands too early.
	stop = func(delayClose bool) {
		if !delayClose {
			close(l.commands)
		}
		for _ = range <-l.results {
			continue
		}
		for _ = range <-l.unilaterals {
			continue
		}
		if delayClose {
			close(l.commands)
		}
		// allow other goroutines to run their shutdown logic;
		// avoid false positives in tests detect leaks
		runtime.Gosched()
	}

	go func() {
		defer func() {
			close(results)
			close(unilaterals)
			close(s.commands)
			for _ = range <-commands {
				continue
			}
			for _ = range <-s.events {
				continue
			}
		}()
		for {
			if ok := expectCommand(); !ok {
				return
			}
			if ok := expectResult(); !ok {
				return
			}
		}
	}()

	return
}
