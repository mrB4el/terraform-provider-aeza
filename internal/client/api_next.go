// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/client/api_next.go
package client

import (
	"context"
	"fmt"

	"github.com/scinfra-pro/terraform-provider-aeza/internal/models/next"
)

func (c *Client) ListServices_V2(ctx context.Context) ([]next.Service, error) {
	var response next.ListServicesResponse
	err := c.NewRequest("GET", "/v2/services", nil).Do(ctx, &response)
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}

func (c *Client) ListServiceIPv6Addresses(ctx context.Context, serviceID int64) ([]string, error) {
	var response struct {
		Items []struct {
			IPs []string `json:"ips"`
		} `json:"items"`
	}
	path := fmt.Sprintf("/v2/services/%d/networks/ipv6", serviceID)
	if err := c.NewRequest("GET", path, nil).Do(ctx, &response); err != nil {
		return nil, fmt.Errorf("failed to get IPv6 addresses for service %d: %w", serviceID, err)
	}

	addresses := make([]string, 0)
	seen := make(map[string]struct{})
	for _, network := range response.Items {
		for _, address := range network.IPs {
			if address == "" {
				continue
			}
			if _, exists := seen[address]; !exists {
				seen[address] = struct{}{}
				addresses = append(addresses, address)
			}
		}
	}
	return addresses, nil
}

func (c *Client) ListServiceGroups_V2(ctx context.Context, serviceType string) ([]next.ServiceGroup, error) {
	var response next.ServiceGroupsResponse

	req := c.NewRequest("GET", "/v2/services/groups", nil)
	if serviceType != "" {
		req.AddQueryParam("type", serviceType)
	}

	err := req.Do(ctx, &response)
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}

func (c *Client) ListOS_V2(ctx context.Context) ([]next.OperatingSystem, error) {
	var response []next.OperatingSystem
	err := c.NewRequest("GET", "/v2/services/operating-systems", nil).Do(ctx, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ControlService_V2(ctx context.Context, serviceID int64, action string) error {
		allowedActions := map[string]bool{
		"suspend": true,
		"resume":  true,
		"restart": true,
	}

	if !allowedActions[action] {
		return fmt.Errorf("invalid service control action: %s", action)
	}

	url := fmt.Sprintf("/v2/services/%d/ctl/%s", serviceID, action)

	err := c.NewRequest("POST", url, nil).Do(ctx, nil)
	if err != nil {
		return fmt.Errorf("service control request failed for action %s: %w", action, err)
	}

	return nil
}
