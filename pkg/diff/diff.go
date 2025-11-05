package diff

import (
	"github.com/scttfrdmn/gh-label-sync/pkg/api"
)

type DiffType string

const (
	DiffTypeCreate DiffType = "create"
	DiffTypeUpdate DiffType = "update"
	DiffTypeMatch  DiffType = "match"
	DiffTypeExtra  DiffType = "extra"
)

type LabelDiff struct {
	Type        DiffType
	Name        string
	Desired     *api.Label
	Current     *api.Label
	ColorChange bool
	DescChange  bool
}

// ComputeDiff compares desired labels with current labels
func ComputeDiff(desired, current []api.Label) []LabelDiff {
	var diffs []LabelDiff

	// Create maps for quick lookup
	currentMap := make(map[string]api.Label)
	for _, label := range current {
		currentMap[label.Name] = label
	}

	desiredMap := make(map[string]api.Label)
	for _, label := range desired {
		desiredMap[label.Name] = label
	}

	// Check desired labels
	for _, desiredLabel := range desired {
		if currentLabel, exists := currentMap[desiredLabel.Name]; exists {
			// Label exists, check if it matches
			colorMatch := api.NormalizeColor(currentLabel.Color) == api.NormalizeColor(desiredLabel.Color)
			descMatch := currentLabel.Description == desiredLabel.Description

			if colorMatch && descMatch {
				diffs = append(diffs, LabelDiff{
					Type:    DiffTypeMatch,
					Name:    desiredLabel.Name,
					Desired: &desiredLabel,
					Current: &currentLabel,
				})
			} else {
				diffs = append(diffs, LabelDiff{
					Type:        DiffTypeUpdate,
					Name:        desiredLabel.Name,
					Desired:     &desiredLabel,
					Current:     &currentLabel,
					ColorChange: !colorMatch,
					DescChange:  !descMatch,
				})
			}
		} else {
			// Label doesn't exist, needs to be created
			diffs = append(diffs, LabelDiff{
				Type:    DiffTypeCreate,
				Name:    desiredLabel.Name,
				Desired: &desiredLabel,
			})
		}
	}

	// Check for extra labels (in repo but not in file)
	for _, currentLabel := range current {
		if _, exists := desiredMap[currentLabel.Name]; !exists {
			diffs = append(diffs, LabelDiff{
				Type:    DiffTypeExtra,
				Name:    currentLabel.Name,
				Current: &currentLabel,
			})
		}
	}

	return diffs
}

// Summary returns counts for each diff type
func Summary(diffs []LabelDiff) (matches, creates, updates, extras int) {
	for _, diff := range diffs {
		switch diff.Type {
		case DiffTypeMatch:
			matches++
		case DiffTypeCreate:
			creates++
		case DiffTypeUpdate:
			updates++
		case DiffTypeExtra:
			extras++
		}
	}
	return
}
