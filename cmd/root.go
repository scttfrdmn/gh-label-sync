package cmd

import (
	"github.com/spf13/cobra"
)

var (
	repoFlag string
)

var rootCmd = &cobra.Command{
	Use:   "label-sync",
	Short: "Bulk label management and synchronization",
	Long: `A GitHub CLI extension for bulk label management and synchronization from YAML/JSON files.

Examples:
  gh label-sync export > labels.yml
  gh label-sync sync --file labels.yml
  gh label-sync clone source/repo --repo target/repo`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&repoFlag, "repo", "R", "", "Repository (owner/repo)")

	// Add subcommands
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(cloneCmd)
}
