package format

import (
	"fmt"
	"strings"

	"github.com/scttfrdmn/gh-label-sync/pkg/diff"
)

// FormatDiff formats a diff for display
func FormatDiff(diffs []diff.LabelDiff, verbose bool) string {
	var sb strings.Builder

	sb.WriteString("Analyzing labels...\n")

	for _, d := range diffs {
		switch d.Type {
		case diff.DiffTypeMatch:
			if verbose {
				sb.WriteString(fmt.Sprintf("  ✓ %s - matches\n", d.Name))
			}
		case diff.DiffTypeCreate:
			sb.WriteString(fmt.Sprintf("  + %s - will create (color: %s)\n", d.Name, d.Desired.Color))
		case diff.DiffTypeUpdate:
			changes := []string{}
			if d.ColorChange {
				changes = append(changes, fmt.Sprintf("color: %s → %s", d.Current.Color, d.Desired.Color))
			}
			if d.DescChange {
				changes = append(changes, "description")
			}
			sb.WriteString(fmt.Sprintf("  ~ %s - differs (%s)\n", d.Name, strings.Join(changes, ", ")))
		case diff.DiffTypeExtra:
			sb.WriteString(fmt.Sprintf("  ⚠ %s - exists but not in file\n", d.Name))
		}
	}

	return sb.String()
}

// FormatSummary formats a summary of changes
func FormatSummary(diffs []diff.LabelDiff, force, deleteUnmanaged bool) string {
	matches, creates, updates, extras := diff.Summary(diffs)

	var sb strings.Builder
	sb.WriteString("\nSummary:\n")

	if matches > 0 {
		sb.WriteString(fmt.Sprintf("  %d label(s) match\n", matches))
	}
	if creates > 0 {
		sb.WriteString(fmt.Sprintf("  %d label(s) to create\n", creates))
	}
	if updates > 0 {
		if force {
			sb.WriteString(fmt.Sprintf("  %d label(s) to update\n", updates))
		} else {
			sb.WriteString(fmt.Sprintf("  %d label(s) differ (use --force to update)\n", updates))
		}
	}
	if extras > 0 {
		if deleteUnmanaged {
			sb.WriteString(fmt.Sprintf("  %d unmanaged label(s) to delete\n", extras))
		} else {
			sb.WriteString(fmt.Sprintf("  %d unmanaged label(s) (use --delete-unmanaged to remove)\n", extras))
		}
	}

	return sb.String()
}

// FormatResult formats the result of sync operation
func FormatResult(created, updated, deleted int) string {
	var parts []string
	if created > 0 {
		parts = append(parts, fmt.Sprintf("%d created", created))
	}
	if updated > 0 {
		parts = append(parts, fmt.Sprintf("%d updated", updated))
	}
	if deleted > 0 {
		parts = append(parts, fmt.Sprintf("%d deleted", deleted))
	}

	if len(parts) == 0 {
		return "✓ No changes made"
	}

	return fmt.Sprintf("✓ Synced labels (%s)", strings.Join(parts, ", "))
}
