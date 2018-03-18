package watchman

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
)

type result struct {
	resp map[string]interface{}
	err  error
}

func reader(ctx context.Context, socket net.Conn) <-chan result {
	// TODO log warnings
	events := make(chan result)
	go func() {
		defer close(events)
		r := bufio.NewReader(socket)

		for {
			result := result{}
			if pdu, err := r.ReadBytes('\n'); err != nil {
				result.err = err
			} else {
				var event map[string]interface{}
				if err = json.Unmarshal(pdu, &event); err != nil {
					result.err = err
				} else {
					result.resp = event
				}
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
