package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/akitanak/watcher/app"
)

// cmd represents the watcher command
var cmd = &cobra.Command{
	Use:   "watcher",
	Short: "run some commands when some files are changed",
	Long: `Watcher is a command-line utility designed to monitor a specified directory and its subdirectories for any file modifications.
Once a change is detected, the tool automatically executes a user-specified command,
providing a flexible and automated solution for various tasks.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
	RunE: func(_ *cobra.Command, args []string) error {
		params.Command = args
		receiver := make(chan []byte)
		go app.NewStdoutPrinter(receiver).Print()
		if err := app.Watch(&params, receiver); err != nil {
			return fmt.Errorf("failed to watch: %w", err)
		}
		return nil
	},
}

var params = app.Params{}

func init() {
	cmd.Flags().StringVarP(&params.Directory, "directory", "d", "./", "directory to watch")
	cmd.Flags().StringSliceVarP(&params.Filters, "filters", "f", []string{}, "filters to watch files")
	cmd.Flags().BoolVarP(&params.ExitOnError, "exit-on-error", "e", false, "exit on error")
}

// Execute executes the watcher command.
func Execute() error {
	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute watcher: %w", err)
	}
	return nil
}
