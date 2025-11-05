package parser

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/scttfrdmn/gh-label-sync/pkg/api"
	"gopkg.in/yaml.v3"
)

type LabelFile struct {
	Labels []api.Label `json:"labels" yaml:"labels"`
}

// ParseFile parses a label file (YAML, JSON, or CSV) based on file extension
func ParseFile(filename string) ([]api.Label, error) {
	var labels []api.Label
	var err error

	// Handle stdin
	if filename == "-" {
		return parseYAML(os.Stdin)
	}

	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Determine format by extension
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".yml", ".yaml":
		labels, err = parseYAML(file)
	case ".json":
		labels, err = parseJSON(file)
	case ".csv":
		labels, err = parseCSV(file)
	default:
		// Try YAML as default
		labels, err = parseYAML(file)
		if err != nil {
			return nil, fmt.Errorf("unsupported file format (use .yml, .json, or .csv): %w", err)
		}
	}

	if err != nil {
		return nil, err
	}

	// Normalize colors
	for i := range labels {
		labels[i].Color = api.NormalizeColor(labels[i].Color)
	}

	return labels, nil
}

func parseYAML(r io.Reader) ([]api.Label, error) {
	var labelFile LabelFile
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&labelFile); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return labelFile.Labels, nil
}

func parseJSON(r io.Reader) ([]api.Label, error) {
	var labelFile LabelFile
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&labelFile); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return labelFile.Labels, nil
}

func parseCSV(r io.Reader) ([]api.Label, error) {
	reader := csv.NewReader(r)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Find column indices
	nameIdx, colorIdx, descIdx := -1, -1, -1
	for i, col := range header {
		switch strings.ToLower(strings.TrimSpace(col)) {
		case "name":
			nameIdx = i
		case "color":
			colorIdx = i
		case "description", "desc":
			descIdx = i
		}
	}

	if nameIdx == -1 || colorIdx == -1 {
		return nil, fmt.Errorf("CSV must have 'name' and 'color' columns")
	}

	// Read rows
	var labels []api.Label
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		label := api.Label{
			Name:  row[nameIdx],
			Color: row[colorIdx],
		}
		if descIdx != -1 && descIdx < len(row) {
			label.Description = row[descIdx]
		}

		labels = append(labels, label)
	}

	return labels, nil
}

// WriteYAML writes labels to YAML format
func WriteYAML(w io.Writer, labels []api.Label) error {
	labelFile := LabelFile{Labels: labels}
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	if err := encoder.Encode(labelFile); err != nil {
		return fmt.Errorf("failed to write YAML: %w", err)
	}
	return nil
}

// WriteJSON writes labels to JSON format
func WriteJSON(w io.Writer, labels []api.Label) error {
	labelFile := LabelFile{Labels: labels}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(labelFile); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}
	return nil
}
