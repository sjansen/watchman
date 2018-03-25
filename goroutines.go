package watchman

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
)

func reader(ctx context.Context, socket net.Conn) <-chan object {
	// TODO log warnings
	events := make(chan object)
	go func() {
		defer close(events)
		bytes := producer(socket)

		for {
			var event object
			select {
			case pdu := <-bytes:
				if err := json.Unmarshal(pdu, &event); err != nil {
					event = object{"error": err.Error()}
				}
			case <-ctx.Done():
				return
			}
			select {
			case events <- event:
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
