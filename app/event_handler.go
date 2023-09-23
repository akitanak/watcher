package app

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/exp/slices"
)

// handleEvent handles event.
func handleEvent(params *Params, event fsnotify.Event) error {
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

			// write empty line for readability.
			fmt.Println()
		}
	}

	return nil
}

var targetOperations = []fsnotify.Op{
	fsnotify.Write,
	fsnotify.Create,
	fsnotify.Remove,
}

// eventMatchAsTargetOp checks if event matches as target op.
func eventMatchAsTargetOp(event fsnotify.Event) bool {
	return slices.ContainsFunc(targetOperations, event.Has)
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
