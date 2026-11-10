// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package data_sources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mrb4el/terraform-provider-aeza/internal/interfaces"
)

func RecipesDataSource() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving list of Aeza recipes",

		ReadContext: recipesRead,

		Schema: map[string]*schema.Schema{
			"recipes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"icon": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"os": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"operating_system": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"targets": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
		},
	}
}

func recipesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(interfaces.DataClient)

	recipes, err := client.ListRecipes(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error fetching recipes: %w", err))
	}

	recipeList := make([]map[string]interface{}, len(recipes))
	for i, recipe := range recipes {
		recipeList[i] = map[string]interface{}{
			"id":               recipe.ID,
			"slug":             recipe.Slug,
			"name":             recipe.Name,
			"icon":             recipe.Icon,
			"description":      recipe.Description,
			"tags":             recipe.Tags,
			"enabled":          recipe.Enabled,
			"os":               recipe.OS,
			"operating_system": recipe.OperatingSystem,
			"targets":          recipe.Targets,
		}
	}

	if err := d.Set("recipes", recipeList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("recipes")
	return nil
}

