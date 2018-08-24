package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sjansen/watchman"
)

func connect() *watchman.Client {
	os.Stdout.Write([]byte("Connecting to Watchman... "))
	os.Stdout.Sync()
	c, err := watchman.Connect()
	if err != nil {
		fmt.Println("FAILURE")
		die(err)
	}
	fmt.Println("SUCCESS")

	return c
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func mkdir() string {
	dir, err := ioutil.TempDir("", "watchman-watch-demo")
	if err != nil {
		die(err)
	}

	path := filepath.Join(dir, ".watchmanconfig")
	err = ioutil.WriteFile(path, []byte(`{"idle_reap_age_seconds": 300}`+"\n"), os.ModePerm)
	if err != nil {
		die(err)
	}

	return dir
}

func main() {
	dir := mkdir()

	c := connect()
	fmt.Printf("version: %s\n\n", c.Version())

	if roots, err := c.WatchList(); err != nil {
		die(err)
	} else {
		fmt.Println("Watches before:")
		for _, root := range roots {
			fmt.Println(" - ", root)
		}
	}

	if watch, err := c.WatchProject(dir); err != nil {
		die(err)
	} else if clock, err := watch.Clock(0); err != nil {
		die(err)
	} else {
		fmt.Printf("\nClock: %s\n\n", clock)
	}

	if roots, err := c.WatchList(); err != nil {
		die(err)
	} else {
		fmt.Println("Watches after:")
		for _, root := range roots {
			fmt.Println(" - ", root)
		}
	}
}
