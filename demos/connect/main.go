package main

import (
	"fmt"
	"os"

	"github.com/sjansen/watchman"
)

func main() {
	os.Stdout.Write([]byte("Connecting to watchman... "))
	os.Stdout.Sync()

	c, err := watchman.Connect()
	if err != nil {
		fmt.Println("FAILURE")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("SUCCESS")
	fmt.Println("version:", c.Version())
	fmt.Println("")
	fmt.Println("capabilities:")
	fmt.Println("    cmd-subscribe\t\t", c.HasCapability("cmd-subscribe"))
	fmt.Println("    field-symlink_target\t", c.HasCapability("field-symlink_target"))
	fmt.Println("    relative_root\t\t", c.HasCapability("relative_root"))
	fmt.Println("    suffix-set\t\t\t", c.HasCapability("suffix-set"))
	fmt.Println("    term-pcre\t\t\t", c.HasCapability("term-pcre"))
	fmt.Println("    wildmatch\t\t\t", c.HasCapability("wildmatch"))
	fmt.Println("")
	fmt.Println("sockname:", c.SockName())
}
