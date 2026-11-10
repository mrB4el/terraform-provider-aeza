// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/client/api_next.go
package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mrb4el/terraform-provider-aeza/internal/models"
	"github.com/mrb4el/terraform-provider-aeza/internal/models/next"
)

func (c *Client) ListServices_V2(ctx context.Context) ([]next.Service, error) {
	var services []next.Service
	for offset := 0; ; {
		var response next.ListServicesResponse
		request := c.NewRequest("GET", "/v2/services", nil)
		request.AddQueryParam("offset", strconv.Itoa(offset))
		request.AddQueryParam("limit", "50")
		if err := request.Do(ctx, &response); err != nil {
			return nil, err
		}
		services = append(services, response.Items...)
		offset += len(response.Items)
		if offset >= response.Total {
			return services, nil
		}
		if len(response.Items) == 0 {
			return nil, fmt.Errorf("service list ended before reaching total %d", response.Total)
		}
	}
}

func (c *Client) ListServiceIPv4Addresses(ctx context.Context, serviceID int64) ([]string, error) {
	var response struct {
		Items []struct {
			Value string `json:"value"`
		} `json:"items"`
	}
	path := fmt.Sprintf("/v2/services/%d/networks/ipv4", serviceID)
	if err := c.NewRequest("GET", path, nil).Do(ctx, &response); err != nil {
		return nil, fmt.Errorf("failed to get IPv4 addresses for service %d: %w", serviceID, err)
	}
	addresses := make([]string, 0, len(response.Items))
	for _, network := range response.Items {
		if network.Value != "" {
			addresses = append(addresses, network.Value)
		}
	}
	return addresses, nil
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

func (c *Client) ListServiceTypes_V2(ctx context.Context) ([]next.ServiceType, error) {
	var serviceTypes []next.ServiceType
	for offset := 0; ; {
		var response struct {
			Items []next.ServiceType `json:"items"`
			Total int                `json:"total"`
		}
		request := c.NewRequest("GET", "/v2/services/types", nil)
		request.AddQueryParam("offset", strconv.Itoa(offset))
		request.AddQueryParam("limit", "50")
		if err := request.Do(ctx, &response); err != nil {
			return nil, err
		}
		serviceTypes = append(serviceTypes, response.Items...)
		offset += len(response.Items)
		if offset >= response.Total {
			return serviceTypes, nil
		}
		if len(response.Items) == 0 {
			return nil, fmt.Errorf("service type list ended before reaching total %d", response.Total)
		}
	}
}

func (c *Client) ListProducts_V2(ctx context.Context) ([]models.Product, error) {
	var response struct {
		Items []struct {
			ID        int64              `json:"id"`
			Name      string             `json:"name"`
			TypeSlug  string             `json:"typeSlug"`
			GroupID   int64              `json:"groupId"`
			RawPrices map[string]float64 `json:"rawPrices"`
			Type      struct {
				Slug string `json:"slug"`
			} `json:"type"`
		} `json:"items"`
	}
	if err := c.NewRequest("GET", "/v2/services/products", nil).Do(ctx, &response); err != nil {
		return nil, err
	}
	products := make([]models.Product, 0, len(response.Items))
	for _, product := range response.Items {
		typeSlug := product.Type.Slug
		if typeSlug == "" {
			typeSlug = product.TypeSlug
		}
		products = append(products, models.Product{
			ID: product.ID, Name: product.Name, Type: typeSlug,
			GroupID: product.GroupID, Prices: product.RawPrices,
		})
	}
	return products, nil
}

func (c *Client) ListOS_V2(ctx context.Context) ([]next.OperatingSystem, error) {
	var operatingSystems []next.OperatingSystem
	for offset := 0; ; {
		var response []next.OperatingSystem
		request := c.NewRequest("GET", "/v2/services/operating-systems", nil)
		request.AddQueryParam("offset", strconv.Itoa(offset))
		request.AddQueryParam("limit", "50")
		if err := request.Do(ctx, &response); err != nil {
			return nil, err
		}
		operatingSystems = append(operatingSystems, response...)
		if len(response) < 50 {
			return operatingSystems, nil
		}
		offset += len(response)
	}
}

func (c *Client) ListRecipes_V2(ctx context.Context) ([]models.Recipe, error) {
	var response []struct {
		ID              int    `json:"id"`
		Slug            string `json:"slug"`
		Name            string `json:"name"`
		Icon            string `json:"icon"`
		Description     string `json:"description"`
		Enabled         bool   `json:"enabled"`
		OperatingSystem string `json:"operatingSystemSlug"`
	}
	if err := c.NewRequest("GET", "/v2/services/operating-systems/recipes", nil).Do(ctx, &response); err != nil {
		return nil, err
	}
	recipes := make([]models.Recipe, 0, len(response))
	for _, recipe := range response {
		recipes = append(recipes, models.Recipe{
			ID: recipe.ID, Slug: recipe.Slug, Name: recipe.Name, Icon: recipe.Icon,
			Description: recipe.Description, Enabled: recipe.Enabled, OperatingSystem: recipe.OperatingSystem,
		})
	}
	return recipes, nil
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

	var body interface{}
	if action == "suspend" || action == "restart" {
		body = struct {
			Force bool `json:"force"`
		}{Force: false}
	}
	err := c.NewRequest("POST", url, body).Do(ctx, nil)
	if err != nil {
		return fmt.Errorf("service control request failed for action %s: %w", action, err)
	}

	return nil
}
