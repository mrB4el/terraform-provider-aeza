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

func ProductsDataSource() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for retrieving list of Aeza products",

		ReadContext: productsRead,

		Schema: map[string]*schema.Schema{
			"products": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_handler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prices": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeFloat,
							},
						},
					},
				},
			},
		},
	}
}

func productsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(interfaces.DataClient)

	products, err := client.ListProducts(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error fetching products: %w", err))
	}

	productList := make([]map[string]interface{}, len(products))
	for i, product := range products {
		productList[i] = map[string]interface{}{
			"id":              product.ID,
			"name":            product.Name,
			"type":            product.Type,
			"group_id":        product.GroupID,
			"service_handler": product.ServiceHandler,
			"prices":          product.Prices,
		}
	}

	if err := d.Set("products", productList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("products")
	return nil
}
