// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// internal/utils/recipe_converter.go
package utils

import (
	"github.com/mrb4el/terraform-provider-aeza/internal/models"
	"github.com/mrb4el/terraform-provider-aeza/internal/models/legacy"
)

func ConvertRecipeFromLegacy(legacyRecipe legacy.Recipe) models.Recipe {
	return models.Recipe{
		ID:              legacyRecipe.ID,
		Slug:            legacyRecipe.Slug,
		Name:            legacyRecipe.Name,
		Icon:            legacyRecipe.Icon,
		Description:     legacyRecipe.Description,
		Tags:            legacyRecipe.Tags,
		Enabled:         legacyRecipe.Enabled,
		OS:              legacyRecipe.OS,
		OperatingSystem: legacyRecipe.OperatingSystem,
		Targets:         legacyRecipe.Targets,
	}
}

