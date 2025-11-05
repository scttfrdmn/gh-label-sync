package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
)

type Client struct {
	restClient *api.RESTClient
	repo       repository.Repository
}

type Label struct {
	Name        string `json:"name" yaml:"name"`
	Color       string `json:"color" yaml:"color"`
	Description string `json:"description" yaml:"description"`
}

type LabelInput struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"`
}

// NewClient creates a new API client
func NewClient(repoOverride string) (*Client, error) {
	restClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create REST client: %w", err)
	}

	var repo repository.Repository
	if repoOverride != "" {
		repo, err = repository.Parse(repoOverride)
		if err != nil {
			return nil, fmt.Errorf("invalid repository format: %w", err)
		}
	} else {
		repo, err = repository.Current()
		if err != nil {
			return nil, fmt.Errorf("could not determine repository (use --repo flag): %w", err)
		}
	}

	return &Client{
		restClient: restClient,
		repo:       repo,
	}, nil
}

// ListLabels lists all labels in the repository
func (c *Client) ListLabels() ([]Label, error) {
	var labels []Label
	path := fmt.Sprintf("repos/%s/%s/labels?per_page=100", c.repo.Owner, c.repo.Name)

	err := c.restClient.Get(path, &labels)
	if err != nil {
		return nil, fmt.Errorf("failed to list labels: %w", err)
	}

	return labels, nil
}

// GetLabel retrieves a specific label by name
func (c *Client) GetLabel(name string) (*Label, error) {
	var label Label
	encodedName := url.PathEscape(name)
	path := fmt.Sprintf("repos/%s/%s/labels/%s", c.repo.Owner, c.repo.Name, encodedName)

	err := c.restClient.Get(path, &label)
	if err != nil {
		return nil, fmt.Errorf("failed to get label: %w", err)
	}

	return &label, nil
}

// CreateLabel creates a new label
func (c *Client) CreateLabel(input LabelInput) (*Label, error) {
	var label Label
	path := fmt.Sprintf("repos/%s/%s/labels", c.repo.Owner, c.repo.Name)

	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	err = c.restClient.Post(path, bytes.NewReader(body), &label)
	if err != nil {
		return nil, fmt.Errorf("failed to create label: %w", err)
	}

	return &label, nil
}

// UpdateLabel updates an existing label
func (c *Client) UpdateLabel(name string, input LabelInput) (*Label, error) {
	var label Label
	encodedName := url.PathEscape(name)
	path := fmt.Sprintf("repos/%s/%s/labels/%s", c.repo.Owner, c.repo.Name, encodedName)

	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	err = c.restClient.Patch(path, bytes.NewReader(body), &label)
	if err != nil {
		return nil, fmt.Errorf("failed to update label: %w", err)
	}

	return &label, nil
}

// DeleteLabel deletes a label
func (c *Client) DeleteLabel(name string) error {
	encodedName := url.PathEscape(name)
	path := fmt.Sprintf("repos/%s/%s/labels/%s", c.repo.Owner, c.repo.Name, encodedName)

	err := c.restClient.Delete(path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}

	return nil
}

// NormalizeColor normalizes color values (strips #, converts to lowercase)
func NormalizeColor(color string) string {
	// Remove # prefix if present
	if len(color) > 0 && color[0] == '#' {
		color = color[1:]
	}
	return color
}
