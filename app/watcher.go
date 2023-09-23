package app

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
				if eventMatchAsTargetOp(event) {
					matched, err := doesFileNameMatchFilters(event.Name, *params)
					if err != nil {
						log.Fatalf("failed to match file name: %s", err)
					}
					if matched {
						logEvent(event)
						log.Printf("run command: %v", strings.Join(params.Command, " "))
						cmd := exec.Command("sh", "-c", strings.Join(params.Command, " "))
						if err := redirectCommandStdoutAndStderr(cmd); err != nil {
							log.Fatalf("failed to redirect stdout and stderr: %s", err)
						}

						if err := cmd.Start(); err != nil {
							log.Fatalf("failed to run command: %s", err)
						}

						if err := cmd.Wait(); err != nil && params.ExitOnError {
							log.Fatalf("failed to wait command: %s", err)
						}
					}
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

// eventMatchAsTargetOp checks if event matches as target op.
func eventMatchAsTargetOp(event fsnotify.Event) bool {
	return event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove)
}

// redirectCommandStdoutAndStderr redirects stdout and stderr to os.Stdout and os.Stderr.
func redirectCommandStdoutAndStderr(cmd *exec.Cmd) error {
	stdout, stderr, err := getStdoutPipeAndStderrPipe(cmd)
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	go func() {
		_, err := io.Copy(os.Stdout, stdout)
		if err != nil {
			fmt.Println("failed to copy stdout: %w", err)
		}
	}()

	go func() {
		_, err := io.Copy(os.Stderr, stderr)
		if err != nil {
			fmt.Println("failed to copy stderr: %w", err)
		}
	}()
	return nil
}

// getStdoutPipeAndStderrPipe gets stdout pipe and stderr pipe.
func getStdoutPipeAndStderrPipe(cmd *exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	return stdout, stderr, nil
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
