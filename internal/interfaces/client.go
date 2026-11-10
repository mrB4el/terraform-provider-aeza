// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/interfaces/client.go
package interfaces

import (
	"context"
	"errors"

	"github.com/mrb4el/terraform-provider-aeza/internal/models"
)

var ErrServiceNotFound = errors.New("service not found")

type DataClient interface {
	ListServices(ctx context.Context) ([]models.TerraformService, error)
	ListProducts(ctx context.Context) ([]models.Product, error)
	ListServiceTypes(ctx context.Context) ([]models.ServiceType, error)
	ListServiceGroups(ctx context.Context, serviceType string) ([]models.ServiceGroup, error)
	ListOS(ctx context.Context) ([]models.OperatingSystem, error)
	ListRecipes(ctx context.Context) ([]models.Recipe, error)
}

type ResourceClient interface {
	DataClient 
	CreateService(ctx context.Context, req models.ServiceCreateRequest) (*models.ServiceCreateResponse, error)
	GetService(ctx context.Context, id int64) (*models.Service, error)
	UpdateService(ctx context.Context, id int64, req models.ServiceUpdateRequest) error
	DeleteService(ctx context.Context, id int64) error

	ProlongService(ctx context.Context, serviceID int64, req models.ServiceProlongRequest) (*models.ServiceProlongResponse, error)
	ControlService(ctx context.Context, serviceID int64, action string) error
}
