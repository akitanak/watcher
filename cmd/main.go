package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/akitanak/watcher/app"
)

// rootCmd represents the base command when called without any subcommands
var cmd = &cobra.Command{
	Use:   "watcher",
	Short: "run some commands when some files are changed",
	Long: `Watcher is a command-line utility designed to monitor a specified directory and its subdirectories for any file modifications.
Once a change is detected, the tool automatically executes a user-specified command,
providing a flexible and automated solution for various tasks.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Args: cobra.MatchAll(cobra.MinimumNArgs(1)),
	RunE: func(_ *cobra.Command, args []string) error {
		params.Command = args
		if err := app.Watch(&params); err != nil {
			return fmt.Errorf("failed to watch: %w", err)
		}
		return nil
	},
}

var params = app.Params{}

func init() {
	cmd.Flags().StringVarP(&params.Directory, "directory", "d", "./", "directory to watch")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("failed to execute watcher: %s", err)
		os.Exit(1)
	}
}
