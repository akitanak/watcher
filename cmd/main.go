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
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
