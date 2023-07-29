package main

import (
	"fmt"
	"os"

	"github.com/akitanak/watcher/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("failed to execute watcher: %s", err)
		os.Exit(1)
	}
}
