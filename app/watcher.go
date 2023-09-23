package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Params represents parameters for watcher.
type Params struct {
	Directory   string
	Filters     []string
	Command     []string
	ExitOnError bool
}

// Watch watches directory and runs command.
func Watch(params *Params) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if err := handleEvent(params, event); err != nil {
					log.Fatalf("failed to handle event: %s", err)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("error: %s", err)
			}
		}
	}()

	err = addDirectoryRecursively(watcher, params.Directory)
	if err != nil {
		return fmt.Errorf("failed to add directory recursively: %w", err)
	}

	<-make(chan struct{})

	return nil
}

func addDirectoryRecursively(watcher *fsnotify.Watcher, watchDir string) error {
	if err := filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if err := watcher.Add(path); err != nil {
				return fmt.Errorf("failed to add directory: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	return nil
}
