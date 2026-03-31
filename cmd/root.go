package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/kostine/kbd/internal/app"
)

var (
	dbPath      string
	versionStr  string
	commitStr   string
	dateStr     string
)

// SetVersion sets version info from main (populated by goreleaser).
func SetVersion(version, commit, date string) {
	versionStr = version
	commitStr = commit
	dateStr = date
}

var rootCmd = &cobra.Command{
	Use:   "kbd",
	Short: "TUI for beads (bd)",
	Long:  "kbd is a terminal user interface for the bd issue tracker.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := exec.LookPath("bd"); err != nil {
			return fmt.Errorf("bd (beads) CLI not found in PATH. Install it first: https://github.com/derailed/beads")
		}
		a := app.New(dbPath)
		return a.Run()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbd %s (commit: %s, built: %s)\n", versionStr, commitStr, dateStr)
	},
}

func init() {
	rootCmd.Flags().StringVar(&dbPath, "db", "", "Database path (passed through to bd)")
	rootCmd.AddCommand(versionCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
