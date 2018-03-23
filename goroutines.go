package watchman

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
)

type result struct {
	resp object
	err  error
}

func reader(ctx context.Context, socket net.Conn) <-chan result {
	// TODO log warnings
	events := make(chan result)
	go func() {
		defer close(events)
		bytes := producer(socket)

		for {
			result := result{}
			select {
			case pdu := <-bytes:
				var event object
				if err := json.Unmarshal(pdu, &event); err != nil {
					result.err = err
				} else {
					result.resp = event
				}
			case <-ctx.Done():
				return
			}
			select {
			case events <- result:
			case <-ctx.Done():
				return
			}
		}

	}()
	return events
}

func writer(ctx context.Context, socket net.Conn) chan<- []interface{} {
	commands := make(chan []interface{})
	go func() {
		defer close(commands)

		for command := range commands {
			pdu, err := json.Marshal(command)
			if err != nil {
				return
			}
			_, err = fmt.Fprintln(socket, string(pdu))
			if err != nil {
				return
			}
		}

	}()
	return commands
}
