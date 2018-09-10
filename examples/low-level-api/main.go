package main

import (
	"fmt"
	"os"

	"github.com/sjansen/watchman/protocol"
)

func main() {
	os.Stdout.Write([]byte("Connecting to Watchman... "))
	os.Stdout.Sync()

	c, err := protocol.Connect()
	if err != nil {
		fmt.Println("FAILURE")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("SUCCESS")
	fmt.Println("sockname:", c.SockName())
	fmt.Println("version: ", c.Version())
	fmt.Println()

	err = c.Send(&protocol.WatchListRequest{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if pdu, err := c.Recv(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		watchList := protocol.NewWatchListResponse(pdu)
		fmt.Println("Watches:")
		for _, root := range watchList.Roots() {
			fmt.Println("  ", root)
		}
	}
}
