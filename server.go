package watchman

import (
	"bufio"
	"context"
	"fmt"
	"net"
)

type server struct {
	commands chan<- string
	events   <-chan []byte
}

func producer(socket net.Conn) <-chan []byte {
	ch := make(chan []byte)
	go func() {
		defer close(ch)
		r := bufio.NewReader(socket)

		for {
			if bytes, err := r.ReadBytes('\n'); err != nil {
				return
			} else {
				ch <- bytes
			}
		}
	}()
	return ch
}

func serverFromSocket(ctx context.Context, socket net.Conn) (s *server) {
	/* SHUTDOWN
	s.commands: closed by caller
	s.events:   closed locally
	*/

	commands := make(chan string)
	events := make(chan []byte)
	s = &server{
		commands: commands,
		events:   events,
	}

	go func() {
		defer func() {
			close(events)
			socket.Close()
			for _ = range <-commands {
				continue
			}
		}()

		bytes := producer(socket)

		for {
			select {
			case command, ok := <-commands:
				if !ok {
					return
				}
				_, err := fmt.Fprintln(socket, command)
				if err != nil {
					return
				}
			case pdu, ok := <-bytes:
				if !ok {
					return
				}
				select {
				case events <- pdu:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return
}
