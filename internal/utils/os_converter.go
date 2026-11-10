// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/utils/os_converter.go
package utils

import (
	"github.com/mrb4el/terraform-provider-aeza/internal/models/legacy"
	"github.com/mrb4el/terraform-provider-aeza/internal/models/next"
)

func ConvertOsFromLegacy(legacy legacy.OperatingSystem) next.OperatingSystem {
	return next.OperatingSystem{
		ID:         legacy.ID,
		Slug:       legacy.Slug,
		Name:       legacy.Name,
		Repository: legacy.Repository,
		Group:      legacy.Group,
		Username:   legacy.Username,
		Enabled:    legacy.Enabled,
		Targets:    legacy.Targets,
		Order:      0, 	}
}
