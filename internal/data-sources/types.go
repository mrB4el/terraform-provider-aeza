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

func ServiceTypesDataSource() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving list of Aeza service types",

		ReadContext: serviceTypesRead,

		Schema: map[string]*schema.Schema{
			"service_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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

func serviceTypesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(interfaces.DataClient)

	serviceTypes, err := client.ListServiceTypes(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error fetching service types: %w", err))
	}

	serviceTypeList := make([]map[string]interface{}, len(serviceTypes))
	for i, serviceType := range serviceTypes {
		serviceTypeList[i] = map[string]interface{}{
			"slug":  serviceType.Slug,
			"name":  serviceType.Name,
			"order": serviceType.Order,
		}
	}

	if err := d.Set("service_types", serviceTypeList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("service_types")
	return nil
}
