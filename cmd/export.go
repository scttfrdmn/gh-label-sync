package cmd

import (
	"fmt"
	"os"

	"github.com/scttfrdmn/gh-label-sync/pkg/api"
	"github.com/scttfrdmn/gh-label-sync/pkg/parser"
	"github.com/spf13/cobra"
)

var (
	exportFormat string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export labels to YAML or JSON",
	Long: `Export repository labels to YAML or JSON format.

Examples:
  gh label-sync export > labels.yml
  gh label-sync export --format json > labels.json
  gh label-sync export --repo owner/repo > labels.yml`,
	RunE: runExport,
}

func init() {
	exportCmd.Flags().StringVar(&exportFormat, "format", "yaml", "Output format (yaml or json)")
}

func runExport(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(repoFlag)
	if err != nil {
		return err
	}

	labels, err := client.ListLabels()
	if err != nil {
		return err
	}

	switch exportFormat {
	case "yaml", "yml":
		err = parser.WriteYAML(os.Stdout, labels)
	case "json":
		err = parser.WriteJSON(os.Stdout, labels)
	default:
		return fmt.Errorf("unsupported format: %s (use yaml or json)", exportFormat)
	}

	if err != nil {
		return err
	}

	return nil
}
