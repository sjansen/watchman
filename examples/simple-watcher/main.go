package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/sjansen/watchman"
	"github.com/sjansen/watchman/protocol"
)

var (
	DIR  = flag.String("d", "", "directory to watch")
	LIST = flag.Bool("l", false, "list watches instead of starting a watch")
)

func init() {
	flag.Parse()
}

type byFilepath []string

func (x byFilepath) Len() int      { return len(x) }
func (x byFilepath) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x byFilepath) Less(i, j int) bool {
	a := x[i]
	b := x[j]
	if strings.ContainsRune(a, filepath.Separator) {
		if !strings.ContainsRune(b, filepath.Separator) {
			return true
		}
	} else if strings.ContainsRune(b, filepath.Separator) {
		return false
	}
	return a < b
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

func listWatches(c *watchman.Client) {
	if roots, err := c.ListWatches(); err != nil {
		die(err)
	} else {
		fmt.Println("Watches:")
		for _, root := range roots {
			fmt.Println("  ", root)
		}
	}
}

func resolveDir() string {
	dir := *DIR
	if dir == "" {
		if wd, err := os.Getwd(); err != nil {
			die(err)
		} else {
			dir = wd
		}
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		die(err)
	}

	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		die(err)
	}

	return dir
}

func main() {
	c := connect()
	fmt.Printf("version: %s\n\n", c.Version())

	if *LIST {
		listWatches(c)
		return
	}

	dir := resolveDir()
	fmt.Println("Watching:", dir)
	watch, err := c.AddWatch(dir)
	if err != nil {
		die(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for pdu := range c.Updates() {
			s := protocol.NewSubscription(pdu)
			fmt.Printf(
				"Update: (clock=%q)\n",
				s.Clock(),
			)
			files := s.Files()
			sort.Sort(byFilepath(files))
			for _, filename := range files {
				fmt.Println(" ", filename)
			}
			fmt.Println()
		}
	}()

	if _, err = watch.Subscribe("example", dir); err != nil {
		die(err)
	}

	wg.Wait()
}
