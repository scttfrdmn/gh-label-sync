package cmd

import (
	"fmt"

	"github.com/scttfrdmn/gh-label-sync/pkg/api"
	"github.com/scttfrdmn/gh-label-sync/pkg/diff"
	"github.com/scttfrdmn/gh-label-sync/pkg/format"
	"github.com/spf13/cobra"
)

var (
	cloneForce bool
)

var cloneCmd = &cobra.Command{
	Use:   "clone <source-repo>",
	Short: "Clone labels from another repository",
	Long: `Clone all labels from a source repository to the target repository.

This is equivalent to exporting labels from the source and syncing to the target.

Examples:
  gh label-sync clone owner/source-repo --repo owner/target-repo
  gh label-sync clone owner/template --repo owner/new-project --force`,
	Args: cobra.ExactArgs(1),
	RunE: runClone,
}

func init() {
	cloneCmd.Flags().BoolVar(&cloneForce, "force", false, "Update existing labels that differ")
}

func runClone(cmd *cobra.Command, args []string) error {
	sourceRepo := args[0]

	if repoFlag == "" {
		return fmt.Errorf("target repository required (use --repo flag)")
	}

	// Get labels from source repository
	fmt.Printf("Fetching labels from %s...\n", sourceRepo)
	sourceClient, err := api.NewClient(sourceRepo)
	if err != nil {
		return fmt.Errorf("failed to connect to source repo: %w", err)
	}

	sourceLabels, err := sourceClient.ListLabels()
	if err != nil {
		return fmt.Errorf("failed to list source labels: %w", err)
	}

	if len(sourceLabels) == 0 {
		return fmt.Errorf("no labels found in source repository")
	}

	fmt.Printf("Found %d label(s) in source repository\n\n", len(sourceLabels))

	// Get labels from target repository
	fmt.Printf("Fetching labels from %s...\n", repoFlag)
	targetClient, err := api.NewClient(repoFlag)
	if err != nil {
		return fmt.Errorf("failed to connect to target repo: %w", err)
	}

	targetLabels, err := targetClient.ListLabels()
	if err != nil {
		return fmt.Errorf("failed to list target labels: %w", err)
	}

	// Compute diff
	diffs := diff.ComputeDiff(sourceLabels, targetLabels)

	// Display diff
	fmt.Print(format.FormatDiff(diffs, false))
	fmt.Print(format.FormatSummary(diffs, cloneForce, false))

	// Check if there are any changes to apply
	_, creates, updates, _ := diff.Summary(diffs)
	totalChanges := creates
	if cloneForce {
		totalChanges += updates
	}

	if totalChanges == 0 {
		fmt.Println("\n✓ All labels are already in sync")
		return nil
	}

	// Apply changes
	created, updated := 0, 0

	for _, d := range diffs {
		switch d.Type {
		case diff.DiffTypeCreate:
			input := api.LabelInput{
				Name:        d.Desired.Name,
				Color:       d.Desired.Color,
				Description: d.Desired.Description,
			}
			_, err := targetClient.CreateLabel(input)
			if err != nil {
				fmt.Printf("  ✗ Failed to create %s: %v\n", d.Name, err)
			} else {
				fmt.Printf("  ✓ Created %s\n", d.Name)
				created++
			}

		case diff.DiffTypeUpdate:
			if cloneForce {
				input := api.LabelInput{
					Name:        d.Desired.Name,
					Color:       d.Desired.Color,
					Description: d.Desired.Description,
				}
				_, err := targetClient.UpdateLabel(d.Name, input)
				if err != nil {
					fmt.Printf("  ✗ Failed to update %s: %v\n", d.Name, err)
				} else {
					fmt.Printf("  ✓ Updated %s\n", d.Name)
					updated++
				}
			}
		}
	}

	fmt.Println()
	fmt.Println(format.FormatResult(created, updated, 0))

	return nil
}
