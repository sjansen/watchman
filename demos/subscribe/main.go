package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sjansen/watchman"
	"github.com/sjansen/watchman/protocol"
)

func coalesce(watch *watchman.Watch) {
	time.Sleep(500 * time.Millisecond)
	if clock, err := watch.Clock(0); err != nil {
		die(err)
	} else {
		fmt.Println("Clock:", clock)
	}
}

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

func touch(dir string, names ...string) {
	for _, name := range names {
		path := filepath.Join(dir, name)
		err := ioutil.WriteFile(path, []byte(name), os.ModePerm)
		if err != nil {
			die(err)
		}
	}
}

func main() {
	c := connect()
	fmt.Printf("version: %s\n\n", c.Version())

	go func() {
		for pdu := range c.Updates() {
			s := protocol.NewSubscription(pdu)
			fmt.Printf(
				"Update detected: (clock=%q)\n",
				s.Clock(),
			)
			for _, filename := range s.Files() {
				fmt.Printf("\t%s\n", filename)
			}
		}
	}()

	dir := mkdir()
	dir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		die(err)
	}

	watch, err := c.WatchProject(dir)
	if err != nil {
		die(err)
	}
	coalesce(watch)

	touch(dir, "foo", "bar", "baz")
	coalesce(watch)

	subscription, err := watch.Subscribe("demo", dir)
	if err != nil {
		die(err)
	}

	touch(dir, "qux", "quux", "corge")
	coalesce(watch)

	err = subscription.Unsubscribe()
	if err != nil {
		die(err)
	}

	touch(dir, "grault", "garply", "waldo")
	coalesce(watch)
}
