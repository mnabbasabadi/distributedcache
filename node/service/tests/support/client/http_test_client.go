// Package client is the HTTP client for the service.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	cacheAPI "github.com/mnabbasbaadi/distributedcache/node/api/v1"
)

// CacheAPITestClient is the HTTP client for the service.
type CacheAPITestClient struct {
	client *cacheAPI.ClientWithResponses
}

// NewGradingAPITestClient creates a new GradingAPITestClient.
func NewGradingAPITestClient(url string, opts ...cacheAPI.ClientOption) (*CacheAPITestClient, error) {
	client, err := cacheAPI.NewClientWithResponses(fmt.Sprintf("http://%s", url), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &CacheAPITestClient{
		client: client,
	}, nil
}

// GetLiveness returns the liveness of the service.
func (c *CacheAPITestClient) GetLiveness(ctx context.Context) error {
	resp, err := c.client.GetLivenessWithResponse(ctx)
	if err != nil {
		return fmt.Errorf("failed to get liveness: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	return nil
}

// GetValue ...
func (c *CacheAPITestClient) GetValue(ctx context.Context, key string) (cacheAPI.KeyValue, error) {
	resp, err := c.client.GetValueWithResponse(ctx, key)
	if err != nil {
		return cacheAPI.KeyValue{}, fmt.Errorf("failed to get gpa: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return cacheAPI.KeyValue{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	var ret cacheAPI.KeyValue
	if err := json.Unmarshal(resp.Body, &ret); err != nil {
		return cacheAPI.KeyValue{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return ret, nil

}

// SetValue ...
func (c *CacheAPITestClient) SetValue(ctx context.Context, key string, value string) error {
	resp, err := c.client.AddKeyWithResponse(ctx, cacheAPI.AddKeyJSONRequestBody{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to get gpa: %w", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	return nil

}
