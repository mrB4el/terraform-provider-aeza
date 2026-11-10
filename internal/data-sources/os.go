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

func OSDataSource() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving list of Aeza operating systems",

		ReadContext: osRead,

		Schema: map[string]*schema.Schema{
			"os_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
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
						"order": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func osRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(interfaces.DataClient)

	osList, err := client.ListOS(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error fetching operating systems: %w", err))
	}

	osListItem := make([]map[string]interface{}, len(osList))
	for i, os := range osList {
		osListItem[i] = map[string]interface{}{
			"id":         os.ID,
			"name":       os.Name,
			"repository": os.Repository,
			"group":      os.Group,
			"enabled":    os.Enabled,
			"slug":       os.Slug,
			"username":   os.Username,
			"targets":    os.Targets,
			"order":      os.Order, 		}
	}

	if err := d.Set("os_list", osListItem); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("os_list")
	return nil
}
