package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/scinfra-pro/terraform-provider-aeza/internal/models"
)

type productLookupClient struct {
	products []models.Product
	err      error
}

func (c productLookupClient) ListProducts(context.Context) ([]models.Product, error) {
	return c.products, c.err
}

func (productLookupClient) ListServices(context.Context) ([]models.TerraformService, error) {
	return nil, nil
}

func (productLookupClient) ListServiceTypes(context.Context) ([]models.ServiceType, error) {
	return nil, nil
}

func (productLookupClient) ListServiceGroups(context.Context, string) ([]models.ServiceGroup, error) {
	return nil, nil
}

func (productLookupClient) ListOS(context.Context) ([]models.OperatingSystem, error) {
	return nil, nil
}

func (productLookupClient) ListRecipes(context.Context) ([]models.Recipe, error) {
	return nil, nil
}

func TestResolveProductIDByName(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ServiceResource().Schema, map[string]interface{}{
		"name":         "example",
		"product_name": "Premium VPS",
	})
	client := productLookupClient{products: []models.Product{
		{ID: 10, Name: "Basic VPS"},
		{ID: 42, Name: "Premium VPS"},
	}}

	got, err := resolveProductID(context.Background(), resourceData, client)
	if err != nil {
		t.Fatalf("resolveProductID returned an error: %v", err)
	}
	if got != 42 {
		t.Fatalf("unexpected product ID: got %d, want 42", got)
	}
	if stateID := resourceData.Get("product_id").(int); stateID != 42 {
		t.Fatalf("resolved product ID was not written to state: %d", stateID)
	}
}

func TestResolveProductIDRejectsAmbiguousName(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ServiceResource().Schema, map[string]interface{}{
		"name":         "example",
		"product_name": "Premium VPS",
	})
	client := productLookupClient{products: []models.Product{
		{ID: 41, Name: "Premium VPS"},
		{ID: 42, Name: "premium vps"},
	}}

	_, err := resolveProductID(context.Background(), resourceData, client)
	if err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("expected ambiguous product error, got %v", err)
	}
}
