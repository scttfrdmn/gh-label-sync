package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/scttfrdmn/gh-label-sync/pkg/api"
	"github.com/scttfrdmn/gh-label-sync/pkg/diff"
	"github.com/scttfrdmn/gh-label-sync/pkg/format"
	"github.com/scttfrdmn/gh-label-sync/pkg/parser"
	"github.com/spf13/cobra"
)

var (
	syncFile           string
	syncDryRun         bool
	syncForce          bool
	syncDeleteUnmanaged bool
	syncYes            bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync labels from a file",
	Long: `Synchronize repository labels from a YAML, JSON, or CSV file.

By default, this command:
- Creates missing labels
- Skips labels that differ (use --force to update)
- Keeps unmanaged labels (use --delete-unmanaged to remove)

Examples:
  gh label-sync sync --file labels.yml
  gh label-sync sync --file labels.json --force
  gh label-sync sync --file labels.yml --dry-run
  gh label-sync sync --file labels.csv --delete-unmanaged --yes`,
	RunE: runSync,
}

func init() {
	syncCmd.Flags().StringVarP(&syncFile, "file", "f", "", "Label definition file (YAML, JSON, or CSV)")
	syncCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "Show what would change without applying")
	syncCmd.Flags().BoolVar(&syncForce, "force", false, "Update existing labels that differ")
	syncCmd.Flags().BoolVar(&syncDeleteUnmanaged, "delete-unmanaged", false, "Delete labels not in file")
	syncCmd.Flags().BoolVarP(&syncYes, "yes", "y", false, "Skip confirmation prompt")
	syncCmd.MarkFlagRequired("file")
}

func runSync(cmd *cobra.Command, args []string) error {
	// Parse label file
	desiredLabels, err := parser.ParseFile(syncFile)
	if err != nil {
		return err
	}

	if len(desiredLabels) == 0 {
		return fmt.Errorf("no labels found in file")
	}

	// Create API client
	client, err := api.NewClient(repoFlag)
	if err != nil {
		return err
	}

	// Get current labels
	currentLabels, err := client.ListLabels()
	if err != nil {
		return err
	}

	// Compute diff
	diffs := diff.ComputeDiff(desiredLabels, currentLabels)

	// Display diff
	fmt.Print(format.FormatDiff(diffs, false))
	fmt.Print(format.FormatSummary(diffs, syncForce, syncDeleteUnmanaged))

	// Check if there are any changes to apply
	_, creates, updates, extras := diff.Summary(diffs)
	totalChanges := creates
	if syncForce {
		totalChanges += updates
	}
	if syncDeleteUnmanaged {
		totalChanges += extras
	}

	if totalChanges == 0 {
		fmt.Println("\n✓ All labels are in sync")
		return nil
	}

	if syncDryRun {
		fmt.Println("\n(dry-run mode: no changes applied)")
		return nil
	}

	// Prompt for confirmation
	if !syncYes {
		fmt.Printf("\n? Apply changes? (y/N) ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Apply changes
	created, updated, deleted := 0, 0, 0

	for _, d := range diffs {
		switch d.Type {
		case diff.DiffTypeCreate:
			input := api.LabelInput{
				Name:        d.Desired.Name,
				Color:       d.Desired.Color,
				Description: d.Desired.Description,
			}
			_, err := client.CreateLabel(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  ✗ Failed to create %s: %v\n", d.Name, err)
			} else {
				fmt.Printf("  ✓ Created %s\n", d.Name)
				created++
			}

		case diff.DiffTypeUpdate:
			if syncForce {
				input := api.LabelInput{
					Name:        d.Desired.Name,
					Color:       d.Desired.Color,
					Description: d.Desired.Description,
				}
				_, err := client.UpdateLabel(d.Name, input)
				if err != nil {
					fmt.Fprintf(os.Stderr, "  ✗ Failed to update %s: %v\n", d.Name, err)
				} else {
					fmt.Printf("  ✓ Updated %s\n", d.Name)
					updated++
				}
			}

		case diff.DiffTypeExtra:
			if syncDeleteUnmanaged {
				err := client.DeleteLabel(d.Name)
				if err != nil {
					fmt.Fprintf(os.Stderr, "  ✗ Failed to delete %s: %v\n", d.Name, err)
				} else {
					fmt.Printf("  ✓ Deleted %s\n", d.Name)
					deleted++
				}
			}
		}
	}

	fmt.Println()
	fmt.Println(format.FormatResult(created, updated, deleted))

	return nil
}
