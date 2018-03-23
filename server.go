package watchman

import (
	"bufio"
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
