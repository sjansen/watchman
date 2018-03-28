package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sjansen/watchman"
)

func main() {
	var dir string
	var err error
	if len(os.Args) > 1 {
		dir, err = filepath.Abs(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		dir, err = filepath.EvalSymlinks(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		dir, err = os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	os.Stdout.Write([]byte("Connecting to watchman... "))
	os.Stdout.Sync()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	c, err := watchman.Connect(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("CONNECTED")

	w, err := c.WatchProject(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	value, err := w.Clock(0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("clock:", value)
}
