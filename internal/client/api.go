// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/client/api.go

package client

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/mrb4el/terraform-provider-aeza/internal/interfaces"
	"github.com/mrb4el/terraform-provider-aeza/internal/models"
	"github.com/mrb4el/terraform-provider-aeza/internal/utils"
)

func (c *Client) ListServices(ctx context.Context) ([]models.TerraformService, error) {
	nextServices, err := c.ListServices_V2(ctx)
	if err != nil {
		return nil, err
	}

	services := utils.ConvertNextServices(nextServices)
	for index, nextService := range nextServices {
		if !hasIPCapability(nextService.Capabilities) {
			continue
		}
		ipv4Addresses, err := c.ListServiceIPv4Addresses(ctx, int64(nextService.ID))
		if err != nil {
			return nil, err
		}
		ipv6Addresses, ipv6Subnets, err := c.ListServiceIPv6Networks(ctx, int64(nextService.ID))
		if err != nil {
			return nil, err
		}
		services[index].IPv4Addresses = appendPrimaryIP(ipv4Addresses, services[index].IPv4Addresses)
		services[index].IPv6Addresses = appendPrimaryIP(ipv6Addresses, services[index].IPv6Addresses)
		services[index].IPv6Subnets = ipv6Subnets
	}
	return services, nil
}

func (c *Client) ListServiceTypes(ctx context.Context) ([]models.ServiceType, error) {
	nextTypes, err := c.ListServiceTypes_V2(ctx)
	if err != nil {
		return nil, err
	}

	var result []models.ServiceType
	for _, serviceType := range nextTypes {
		result = append(result, models.ServiceType{
			Slug:    serviceType.Slug,
			Name:    serviceType.Name,
			Payload: serviceType.Payload,
		})
	}

	return result, nil
}

func (c *Client) ListServiceGroups(ctx context.Context, serviceType string) ([]models.ServiceGroup, error) {
	nextGroups, err := c.ListServiceGroups_V2(ctx, serviceType)
	if err != nil {
		return nil, err
	}

	return utils.ConvertNextServiceGroups(nextGroups), nil
}

func (c *Client) ListProducts(ctx context.Context) ([]models.Product, error) {
	products, err := c.ListProducts_V2(ctx)
	if err != nil {
		return nil, err
	}
	serviceTypes, err := c.ListServiceTypes_V2(ctx)
	if err != nil {
		return nil, err
	}
	handlersByType := make(map[string]string, len(serviceTypes))
	for _, serviceType := range serviceTypes {
		handlersByType[serviceType.Slug] = serviceType.ServiceHandler
	}
	for index := range products {
		products[index].ServiceHandler = handlersByType[products[index].Type]
	}
	return products, nil
}

func (c *Client) ListOS(ctx context.Context) ([]models.OperatingSystem, error) {
	nextOS, err := c.ListOS_V2(ctx)
	if err != nil {
		return nil, err
	}

	var result []models.OperatingSystem
	for _, os := range nextOS {
		result = append(result, models.OperatingSystem{OperatingSystem: os})
	}

	return result, nil
}

func (c *Client) ListRecipes(ctx context.Context) ([]models.Recipe, error) {
	recipes, err := c.ListRecipes_V2(ctx)
	if err != nil {
		return nil, err
	}

	operatingSystems, err := c.ListOS_V2(ctx)
	if err != nil {
		return nil, err
	}
	idsBySlug := make(map[string]int, len(operatingSystems))
	for _, operatingSystem := range operatingSystems {
		idsBySlug[operatingSystem.Slug] = operatingSystem.ID
	}
	for index := range recipes {
		recipes[index].OS = idsBySlug[recipes[index].OperatingSystem]
	}
	return recipes, nil
}

func (c *Client) CreateService(ctx context.Context, req models.ServiceCreateRequest) (*models.ServiceCreateResponse, error) {
		legacyReq := utils.ConvertToLegacy_ServiceCreateRequest(req)

		legacyResp, err := c.CreateService_legacy(ctx, legacyReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create service via legacy API: %w", err)
	}

		terraformResp := utils.ConvertFromLegacy_ServiceCreateResponse(*legacyResp)

	return &terraformResp, nil
}

func (c *Client) DeleteService(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/v2/services/%d", id)
	if err := c.NewRequest("DELETE", path, nil).Do(ctx, nil); err != nil {
		var apiError *APIError
		if errors.As(err, &apiError) && apiError.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("failed to request deletion of service %d: %w", id, err)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		var response struct {
			Status string `json:"status"`
		}
		if err := c.NewRequest("GET", path, nil).Do(ctx, &response); err != nil {
			var apiError *APIError
			if errors.As(err, &apiError) && apiError.StatusCode == http.StatusNotFound {
				return nil
			}
			return fmt.Errorf("failed to check deletion of service %d: %w", id, err)
		}
		if response.Status == "deleted" {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for deletion of service %d: %w", id, ctx.Err())
		case <-ticker.C:
		}
	}
}

func (c *Client) GetService(ctx context.Context, id int64) (*models.Service, error) {
	var response struct {
		ID               int64                  `json:"id"`
		OwnerID          int64                  `json:"ownerId"`
		ProductID        int64                  `json:"productId"`
		Name             string                 `json:"name"`
		IP               string                 `json:"ip"`
		PaymentTerm      string                 `json:"paymentTerm"`
		Parameters       map[string]interface{} `json:"parameters"`
		SecureParameters map[string]interface{} `json:"secureParameters"`
		AutoProlong      bool                   `json:"autoProlong"`
		Backups          bool                   `json:"backups"`
		Status           string                 `json:"status"`
		Capabilities     []string               `json:"capabilities"`
		LocationCode     string                 `json:"locationCode"`
		CreatedAt        string                 `json:"createdAt"`
		ExpiresAt        string                 `json:"expiresAt"`
		Product          struct {
			ID       int64  `json:"id"`
			Name     string `json:"name"`
			TypeSlug string `json:"typeSlug"`
			GroupID  int64  `json:"groupId"`
		} `json:"product"`
	}
	path := fmt.Sprintf("/v2/services/%d", id)
	if err := c.NewRequest("GET", path, nil).Do(ctx, &response); err != nil {
		var apiError *APIError
		if errors.As(err, &apiError) && apiError.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%w: %d", interfaces.ErrServiceNotFound, id)
		}
		return nil, fmt.Errorf("failed to get service: %w", err)
	}
	terraformService := models.Service{
		ID: response.ID, OwnerID: response.OwnerID, ProductID: response.ProductID,
		Name: response.Name, IP: response.IP, PaymentTerm: response.PaymentTerm,
		Parameters: response.Parameters, SecureParameters: response.SecureParameters,
		AutoProlong: response.AutoProlong, Backups: response.Backups, Status: response.Status,
		LocationCode: response.LocationCode, CreatedAt: response.CreatedAt, ExpiresAt: response.ExpiresAt,
		Product: models.Product{ID: response.Product.ID, Name: response.Product.Name, Type: response.Product.TypeSlug, GroupID: response.Product.GroupID},
	}
	if ip := net.ParseIP(terraformService.IP); ip != nil {
		if ip.To4() != nil {
			terraformService.IPv4Addresses = []string{terraformService.IP}
		} else {
			terraformService.IPv6Addresses = []string{terraformService.IP}
		}
	}
	if !hasIPCapability(response.Capabilities) {
		return &terraformService, nil
	}
	ipv4Addresses, err := c.ListServiceIPv4Addresses(ctx, id)
	if err != nil {
		return nil, err
	}
	terraformService.IPv4Addresses = appendPrimaryIP(ipv4Addresses, terraformService.IPv4Addresses)

	ipv6Addresses, ipv6Subnets, err := c.ListServiceIPv6Networks(ctx, id)
	if err != nil {
		return nil, err
	}
	terraformService.IPv6Addresses = appendPrimaryIP(ipv6Addresses, terraformService.IPv6Addresses)
	terraformService.IPv6Subnets = ipv6Subnets

	return &terraformService, nil
}

func hasIPCapability(capabilities []string) bool {
	for _, capability := range capabilities {
		if capability == "ip" {
			return true
		}
	}
	return false
}

func appendPrimaryIP(addresses, primary []string) []string {
	if len(primary) == 0 {
		return addresses
	}
	for _, address := range addresses {
		if address == primary[0] {
			return addresses
		}
	}
	return append(addresses, primary[0])
}

func (c *Client) UpdateService(ctx context.Context, id int64, req models.ServiceUpdateRequest) error {
	body := struct {
		Name        *string `json:"name,omitempty"`
		AutoProlong *bool   `json:"autoProlong,omitempty"`
	}{Name: req.Name, AutoProlong: req.AutoProlong}
	path := fmt.Sprintf("/v2/services/%d", id)
	return c.NewRequest("PATCH", path, body).Do(ctx, nil)
}

func (c *Client) ProlongService(ctx context.Context, serviceID int64, req models.ServiceProlongRequest) (*models.ServiceProlongResponse, error) {

	legacyReq := utils.ConvertToLegacy_ServiceProlongRequest(req)

		legacyResp, err := c.ProlongService_Legacy(ctx, serviceID, legacyReq)
	if err != nil {
		return nil, err
	}

	resp := utils.ConvertFromLegacy_ServiceProlongResponse(legacyResp)

	return resp, nil
}

func (c *Client) ControlService(ctx context.Context, serviceID int64, action string) error {

	return c.ControlService_V2(ctx, serviceID, action)
}
