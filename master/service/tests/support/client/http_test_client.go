// Package client is the HTTP client for the service.
package client

import (
	"context"
	"fmt"
	"net/http"

	registryAPI "github.com/mnabbasbaadi/distributedcache/master/api/v1"
)

// RegisterAPITestClient is the HTTP client for the service.
type RegisterAPITestClient struct {
	client *registryAPI.ClientWithResponses
}

// NewGradingAPITestClient creates a new GradingAPITestClient.
func NewGradingAPITestClient(url string, opts ...registryAPI.ClientOption) (*RegisterAPITestClient, error) {
	client, err := registryAPI.NewClientWithResponses(fmt.Sprintf("http://%s", url), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &RegisterAPITestClient{
		client: client,
	}, nil
}

// GetLiveness returns the liveness of the service.
func (c *RegisterAPITestClient) GetLiveness(ctx context.Context) error {
	resp, err := c.client.GetLivenessWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to get liveness: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	return nil
}

// RegisterNode ...
func (c *RegisterAPITestClient) RegisterNode(ctx context.Context) error {
	resp, err := c.client.RegisterNodeWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gpa: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil

}

// UnRegisterNode ...
func (c *RegisterAPITestClient) UnRegisterNode(ctx context.Context, addr string) error {
	resp, err := c.client.UnRegisterNodeWithResponse(ctx, addr)
	if err != nil {
		return fmt.Errorf("failed to get gpa: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil

}
