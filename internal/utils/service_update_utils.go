// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/utils/service_update_utils.go
package utils

import (
	"github.com/mrb4el/terraform-provider-aeza/internal/models"
	"github.com/mrb4el/terraform-provider-aeza/internal/models/legacy"
)

func ConvertToLegacy_ServiceUpdateRequest(req models.ServiceUpdateRequest) legacy.ServiceUpdateRequest {
	return legacy.ServiceUpdateRequest{
		Name:        req.Name,
		AutoProlong: req.AutoProlong,
	}
}
