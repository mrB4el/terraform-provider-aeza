// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// /internal/resources/service.go
package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/scinfra-pro/terraform-provider-aeza/internal/interfaces"
	"github.com/scinfra-pro/terraform-provider-aeza/internal/models"
)

func ServiceResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for managing Aeza services",

		CreateContext: serviceCreate,
		ReadContext:   serviceRead,
		UpdateContext: serviceUpdate,
		DeleteContext: serviceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"product_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"product_name"},
				ValidateFunc:  validation.IntAtLeast(1),
				Description:   "Aeza product ID. Use either product_id or product_name.",
			},
			"payment_term": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "month",
			},
			"auto_prolong": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"os": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true},
			"recipe": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true},
			"iso_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true},
			"order_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"product_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"product_id"},
				ValidateFunc:  validation.StringIsNotEmpty,
				Description:   "Exact Aeza product name. The provider resolves it to product_id.",
			},
			"location_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"price": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transaction_amount": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_addresses": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "All IPv4 addresses assigned to the service.",
			},
			"ipv6_addresses": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "All IPv6 addresses assigned to the service.",
			},
			"wait_for_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Wait for at least one IPv4 or IPv6 address after creating the service.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service create timestamp",
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service expiration timestamp",
			},
			"purchased_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service purchase timestamp",
			},
		},
	}
}

func serviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(interfaces.ResourceClient)

	productID, err := resolveProductID(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	req := models.ServiceCreateRequest{
		Name:        d.Get("name").(string),
		ProductID:   productID,
		PaymentTerm: d.Get("payment_term").(string),
		AutoProlong: d.Get("auto_prolong").(bool),
		OS:          d.Get("os").(string),
		Recipe:      d.Get("recipe").(string),
		IsoURL:      d.Get("iso_url").(string),
	}

	resp, err := client.CreateService(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating service: %w", err))
	}

	if resp.ID == 0 {
		return diag.FromErr(fmt.Errorf("API returned ID=0, creation failed"))
	}

	d.SetId(fmt.Sprintf("%d", resp.ID))

	if err := setServiceCreateResponse(d, resp); err != nil {
		return diag.FromErr(fmt.Errorf("error setting service response data: %w", err))
	}

	if d.Get("wait_for_ip").(bool) {
		return waitForServiceIP(ctx, d, client)
	}

	return serviceRead(ctx, d, meta)
}

func resolveProductID(ctx context.Context, d *schema.ResourceData, client interfaces.DataClient) (int64, error) {
	if rawID, ok := d.GetOk("product_id"); ok {
		return int64(rawID.(int)), nil
	}

	rawName, ok := d.GetOk("product_name")
	if !ok || strings.TrimSpace(rawName.(string)) == "" {
		return 0, fmt.Errorf("one of product_id or product_name must be configured")
	}

	requestedName := strings.TrimSpace(rawName.(string))
	products, err := client.ListProducts(ctx)
	if err != nil {
		return 0, fmt.Errorf("error resolving product_name %q: %w", requestedName, err)
	}

	var matches []models.Product
	for _, product := range products {
		if strings.EqualFold(strings.TrimSpace(product.Name), requestedName) {
			matches = append(matches, product)
		}
	}

	if len(matches) == 0 {
		return 0, fmt.Errorf("Aeza product named %q was not found", requestedName)
	}
	if len(matches) > 1 {
		return 0, fmt.Errorf("Aeza product name %q is ambiguous (%d matches); use product_id instead", requestedName, len(matches))
	}

	if err := d.Set("product_id", matches[0].ID); err != nil {
		return 0, fmt.Errorf("error setting resolved product_id: %w", err)
	}
	if err := d.Set("product_name", matches[0].Name); err != nil {
		return 0, fmt.Errorf("error setting canonical product_name: %w", err)
	}

	return matches[0].ID, nil
}

func waitForServiceIP(ctx context.Context, d *schema.ResourceData, client interfaces.ResourceClient) diag.Diagnostics {
	id, err := parseIntID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	timer := time.NewTimer(d.Timeout(schema.TimeoutCreate))
	defer timer.Stop()

	var lastErr error
	for {
		service, getErr := client.GetService(ctx, id)
		switch {
		case getErr != nil:
			lastErr = getErr
		case service == nil:
			lastErr = fmt.Errorf("service is not available yet")
		case len(service.IPv4Addresses) == 0 && len(service.IPv6Addresses) == 0 && service.IP == "":
			lastErr = fmt.Errorf("IP addresses are not available yet")
		default:
			if setErr := setServiceGetResponse(d, *service); setErr != nil {
				return diag.FromErr(setErr)
			}
			return nil
		}

		select {
		case <-ctx.Done():
			return diag.FromErr(fmt.Errorf("waiting for service %d IP addresses: %w", id, ctx.Err()))
		case <-timer.C:
			return diag.FromErr(fmt.Errorf("service was created, but its IP addresses were not available before the create timeout: %w", lastErr))
		case <-ticker.C:
		}
	}
}

func serviceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(interfaces.ResourceClient)

	id, err := parseIntID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	service, err := client.GetService(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading service: %w", err))
	}

	if service == nil {
		d.SetId("")
		return nil
	}

	if err := setServiceGetResponse(d, *service); err != nil {
		return diag.FromErr(fmt.Errorf("error setting service data: %w", err))
	}

	return nil
}

func setServiceGetResponse(d *schema.ResourceData, service models.Service) error {
	if err := setFieldIfNotEmpty(d, "name", service.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("product_id", service.ProductID); err != nil {
		return fmt.Errorf("error setting product_id: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "product_name", service.Product.Name); err != nil {
		return fmt.Errorf("error setting product_name: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "payment_term", service.PaymentTerm); err != nil {
		return fmt.Errorf("error setting payment_term: %w", err)
	}
	if err := d.Set("auto_prolong", service.AutoProlong); err != nil {
		return fmt.Errorf("error setting auto_prolong: %w", err)
	}

	if err := setServiceParameters(d, service.Parameters); err != nil {
		return fmt.Errorf("error setting service parameters: %w", err)
	}

	if err := setFieldIfNotEmpty(d, "status", service.Status); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	if err := d.Set("ip", service.IP); err != nil {
		return fmt.Errorf("error setting ip: %w", err)
	}
	if err := d.Set("ipv4_addresses", service.IPv4Addresses); err != nil {
		return fmt.Errorf("error setting ipv4_addresses: %w", err)
	}
	if err := d.Set("ipv6_addresses", service.IPv6Addresses); err != nil {
		return fmt.Errorf("error setting ipv6_addresses: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "created_at", service.CreatedAt); err != nil {
		return fmt.Errorf("error setting created_at: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "updated_at", service.UpdatedAt); err != nil {
		return fmt.Errorf("error setting updated_at: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "expires_at", service.ExpiresAt); err != nil {
		return fmt.Errorf("error setting expires_at: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "purchased_at", service.PurchasedAt); err != nil {
		return fmt.Errorf("error setting purchased_at: %w", err)
	}

	return nil
}

func setServiceParameters(d *schema.ResourceData, parameters map[string]interface{}) error {
	if parameters == nil {
		return nil
	}

	if os, exists := parameters["os"]; exists {
		if osStr, ok := os.(string); ok && osStr != "" {
			if err := d.Set("os", osStr); err != nil {
				return fmt.Errorf("error setting os: %w", err)
			}
		}
	}

	if isoURL, exists := parameters["isoUrl"]; exists {
		if isoURLStr, ok := isoURL.(string); ok {
			if err := d.Set("iso_url", isoURLStr); err != nil {
				return fmt.Errorf("error setting iso_url: %w", err)
			}
		}
	}

	if recipe, exists := parameters["recipe"]; exists && recipe != nil {
		if recipeStr, ok := recipe.(string); ok && recipeStr != "" {
			if err := d.Set("recipe", recipeStr); err != nil {
				return fmt.Errorf("error setting recipe: %w", err)
			}
		}
	}

	return nil
}

func setFieldIfNotEmpty(d *schema.ResourceData, key, value string) error {
	if value == "" {
		return nil
	}
	return d.Set(key, value)
}

func serviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(interfaces.ResourceClient)

	id, err := parseIntID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	req := models.ServiceUpdateRequest{}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		req.Name = &name
	}

	if d.HasChange("auto_prolong") {
		autoProlong := d.Get("auto_prolong").(bool)
		req.AutoProlong = &autoProlong
	}

	if req.Name != nil || req.AutoProlong != nil {
		err := client.UpdateService(ctx, id, req)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating service: %w", err))
		}
	}

	return serviceRead(ctx, d, meta)
}

func serviceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(interfaces.ResourceClient)

	id, err := parseIntID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteService(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting service: %w", err))
	}

	d.SetId("")
	return nil
}

func setServiceCreateResponse(d *schema.ResourceData, resp *models.ServiceCreateResponse) error {
	if err := d.Set("order_id", resp.OrderID); err != nil {
		return fmt.Errorf("error setting order_id: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "date", resp.Date); err != nil {
		return fmt.Errorf("error setting date: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "product_type", resp.ProductType); err != nil {
		return fmt.Errorf("error setting product_type: %w", err)
	}
	if resp.GroupId != nil {
		if err := d.Set("group_id", *resp.GroupId); err != nil {
			return fmt.Errorf("error setting group_id: %w", err)
		}
	}
	if err := setFieldIfNotEmpty(d, "product_name", resp.ProductName); err != nil {
		return fmt.Errorf("error setting product_name: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "location_name", resp.LocationName); err != nil {
		return fmt.Errorf("error setting location_name: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "price", resp.Price); err != nil {
		return fmt.Errorf("error setting price: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "transaction_amount", resp.TransactionAmount); err != nil {
		return fmt.Errorf("error setting transaction_amount: %w", err)
	}
	if err := setFieldIfNotEmpty(d, "status", resp.Status); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}

	return nil
}

func parseIntID(id string) (int64, error) {
	var intID int64
	_, err := fmt.Sscanf(id, "%d", &intID)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %s", id)
	}
	return intID, nil
}
