// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// /internal/resources/service_actions.go
package resources

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mrb4el/terraform-provider-aeza/internal/interfaces"
)

func ServiceActionsResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for managing service actions (on, off, reboot)",

		CreateContext: actionsCreate,
		ReadContext:   actionsRead,
		UpdateContext: actionsUpdate,
		DeleteContext: actionsDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the service to manage",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Power state of the service (true = resume/power on, false = suspend/power off)",
			},
			"restart": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Restart the service",
			},
		},
	}
}

func actionsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(interfaces.ResourceClient)
	serviceID := int64(d.Get("service_id").(int))

	service, err := client.GetService(ctx, serviceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting service status: %w", err))
	}

	isActive := d.Get("is_active").(bool)

	if (service.Status == "active" && !isActive) || (service.Status == "suspended" && isActive) {
		var action string
		if isActive {
			action = "resume"
		} else {
			action = "suspend"
		}

		if err := client.ControlService(ctx, serviceID, action); err != nil {
			return diag.FromErr(fmt.Errorf("error %s service: %w", action, err))
		}
		log.Printf("[INFO] Successfully %s service %d (status: %s -> desired: %v)",
			action, serviceID, service.Status, isActive)
	} else {
		log.Printf("[DEBUG] Skipping power operation - service %d already in desired state: %s",
			serviceID, service.Status)
	}

	if d.Get("restart").(bool) {
		if service.Status == "active" {
			if err := client.ControlService(ctx, serviceID, "restart"); err != nil {
				return diag.FromErr(fmt.Errorf("error restarting service: %w", err))
			}
			log.Printf("[INFO] Successfully restarted service %d", serviceID)
		} else {
			log.Printf("[WARN] Cannot restart service %d - current status is %s", serviceID, service.Status)
		}
	}

	d.SetId(fmt.Sprintf("%d", serviceID))
	return actionsRead(ctx, d, meta)
}

func actionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// For service actions, we don't have much to read back from API
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

func actionsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// For update, we can handle is_active or restart operation changes
	if d.HasChange("is_active") || d.HasChange("restart") {
		return actionsCreate(ctx, d, meta)
	}
	return nil
}

func actionsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// For service actions, delete doesn't actually delete anything
	// It just removes the resource from state
	d.SetId("")
	return nil
}
