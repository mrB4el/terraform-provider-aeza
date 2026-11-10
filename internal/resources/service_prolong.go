// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// /internal/resources/service_prolong.go
package resources

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrb4el/terraform-provider-aeza/internal/interfaces"
	"github.com/mrb4el/terraform-provider-aeza/internal/models"
)

func ServiceProlongResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for prolonging Aeza services",

		CreateContext: serviceProlongCreate,
		ReadContext:   serviceProlongRead,
		UpdateContext: serviceProlongUpdate,
		DeleteContext: serviceProlongDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the service to prolong",
			},
			"method": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "balance",
				Description: "Payment method (balance, sms, etc.)",
			},
			"term": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Payment term (hour, day, month)",
			},
			"term_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of terms to prolong",
			},
			"force": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force prolongation on every apply",
			},
						"transaction_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the transaction created by the prolong operation",
			},
			"transaction_amount": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Amount of the transaction",
			},
			"transaction_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the transaction (created, pending, completed)",
			},
			"transaction_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the transaction (prolong, etc.)",
			},
			"transaction_created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the transaction was created",
			},
		},
	}
}

func serviceProlongCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(interfaces.ResourceClient)
	serviceID := int64(d.Get("service_id").(int))

	req := models.ServiceProlongRequest{
		Method: d.Get("method").(string),
		Term:   d.Get("term").(string),
		Count:  int64(d.Get("term_count").(int)),
	}

	resp, err := client.ProlongService(ctx, serviceID, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error prolonging service: %w", err))
	}

	if resp != nil && resp.Transaction != nil {
		if err := d.Set("transaction_id", resp.Transaction.ID); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("transaction_amount", resp.Transaction.Amount); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("transaction_status", resp.Transaction.Status); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("transaction_type", resp.Transaction.Type); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("transaction_created_at", resp.Transaction.CreatedAt); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO] Successfully prolonged service %d for %d %s", serviceID, req.Count, req.Term)

	d.SetId(fmt.Sprintf("%d", serviceID))
	return serviceProlongRead(ctx, d, meta)
}

func serviceProlongRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// For service prolong, we don't have much to read back from API
	// Just verify the service exists
	client := meta.(interfaces.ResourceClient)

	serviceID, err := parseIntID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Verify service exists
	_, err = client.GetService(ctx, serviceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading service: %w", err))
	}

	return nil
}

func serviceProlongUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

		if d.Get("force").(bool) {
		return serviceProlongCreate(ctx, d, meta)
	}

	if d.HasChange("method") || d.HasChange("term") || d.HasChange("term_count") {
		return serviceProlongCreate(ctx, d, meta)
	}
	return nil
}

func serviceProlongDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// For service actions, delete doesn't actually delete anything
	// It just removes the resource from state
	d.SetId("")
	return nil
}
