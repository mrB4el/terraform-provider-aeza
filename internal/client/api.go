// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/client/api.go

package client

import (
	"context"
	"fmt"

	"github.com/scinfra-pro/terraform-provider-aeza/internal/models"
	"github.com/scinfra-pro/terraform-provider-aeza/internal/utils"
)

func (c *Client) ListServices(ctx context.Context) ([]models.TerraformService, error) {
	legacyServices, err := c.ListServicesVPS_legacy(ctx)
	if err == nil && len(legacyServices) > 0 {
		return utils.ConvertLegacyServices(legacyServices), nil
	}

	nextServices, err := c.ListServices_V2(ctx)
	if err != nil {
		return nil, err
	}

	return utils.ConvertNextServices(nextServices), nil
}

func (c *Client) ListServiceTypes(ctx context.Context) ([]models.ServiceType, error) {

	legacyTypes, err := c.ListServiceTypes_legacy(ctx)
	if err != nil {
		return nil, err
	}

	var result []models.ServiceType
	for _, legacyType := range legacyTypes {
		result = append(result, utils.ConvertLegacyServiceType(legacyType))
	}

	return result, nil
}

func (c *Client) ListServiceGroups(ctx context.Context, serviceType string) ([]models.ServiceGroup, error) {
	legacyGroups, err := c.ListServiceGroups_Legacy(ctx, serviceType)
	if err == nil && len(legacyGroups) > 0 {
		return utils.ConvertLegacyServiceGroups(legacyGroups), nil
	}

	nextGroups, err := c.ListServiceGroups_V2(ctx, serviceType)
	if err != nil {
		return nil, err
	}

	return utils.ConvertNextServiceGroups(nextGroups), nil
}

func (c *Client) ListProducts(ctx context.Context) ([]models.Product, error) {
	legacyProducts, err := c.ListProducts_Legacy(ctx)
	if err != nil {
		return nil, err
	}

	var result []models.Product
	for _, legacyProduct := range legacyProducts {
		result = append(result, utils.ConvertLegacyProduct(legacyProduct))
	}

	return result, nil
}

func (c *Client) ListOS(ctx context.Context) ([]models.OperatingSystem, error) {
	nextOS, err := c.ListOS_V2(ctx)
	if err == nil && len(nextOS) > 0 {
		var result []models.OperatingSystem
		for _, os := range nextOS {
			result = append(result, models.OperatingSystem{OperatingSystem: os})
		}
		return result, nil
	}

	legacyOS, err := c.ListOS_Legacy(ctx)
	if err != nil {
		return nil, err
	}

	var result []models.OperatingSystem
	for _, os := range legacyOS {
		result = append(result, models.OperatingSystem{OperatingSystem: utils.ConvertOsFromLegacy(os)})
	}

	return result, nil
}

func (c *Client) ListRecipes(ctx context.Context) ([]models.Recipe, error) {
	legacyRecipes, err := c.ListRecipes_Legacy(ctx)
	if err != nil {
		return nil, err
	}

	var result []models.Recipe
	for _, legacyRecipe := range legacyRecipes {
		result = append(result, utils.ConvertRecipeFromLegacy(legacyRecipe))
	}

	return result, nil
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
	return c.DeleteService_legacy(ctx, id)
}

func (c *Client) GetService(ctx context.Context, id int64) (*models.Service, error) {

		legacyResp, err := c.GetService_legacy(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

		if len(legacyResp.Data.Items) == 0 {
		return nil, fmt.Errorf("service with ID %d not found", id)
	}

	legacyService := legacyResp.Data.Items[0]
	terraformService := utils.ConvertLegacyServiceGetToTerraform(legacyService)

	ipv6Addresses, err := c.ListServiceIPv6Addresses(ctx, id)
	if err != nil {
		return nil, err
	}
	terraformService.IPv6Addresses = ipv6Addresses

	return &terraformService, nil
}

func (c *Client) UpdateService(ctx context.Context, id int64, req models.ServiceUpdateRequest) error {
	legacyReq := utils.ConvertToLegacy_ServiceUpdateRequest(req)
	return c.UpdateService_legacy(ctx, id, legacyReq)
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
