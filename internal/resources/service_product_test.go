package resources

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mrb4el/terraform-provider-aeza/internal/client"
	"github.com/mrb4el/terraform-provider-aeza/internal/models"
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

func TestServiceReadDistinguishesServiceAndNetworkNotFound(t *testing.T) {
	for _, missingService := range []bool{true, false} {
		t.Run(fmt.Sprintf("missing_service_%t", missingService), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if !missingService && request.URL.Path == "/v2/services/123" {
					fmt.Fprint(writer, `{"id":123,"capabilities":["ip"]}`)
					return
				}
				if !missingService && request.URL.Path == "/v2/services/123/networks/ipv4" {
					fmt.Fprint(writer, `{"total":0,"items":[]}`)
					return
				}
				http.NotFound(writer, request)
			}))
			defer server.Close()
			apiClient, _ := client.NewClient(server.URL, "test-key")
			data := schema.TestResourceDataRaw(t, ServiceResource().Schema, nil)
			data.SetId("123")
			diagnostics := serviceRead(context.Background(), data, apiClient)
			if missingService && (data.Id() != "" || diagnostics.HasError()) {
				t.Fatalf("missing service: ID=%q diagnostics=%v", data.Id(), diagnostics)
			}
			if !missingService && (data.Id() != "123" || !diagnostics.HasError()) {
				t.Fatalf("missing network: ID=%q diagnostics=%v", data.Id(), diagnostics)
			}
		})
	}
}

func TestServiceDeleteKeepsIDUntilDeletionCompletes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodDelete {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		fmt.Fprint(writer, `{"status":"active"}`)
	}))
	defer server.Close()
	apiClient, _ := client.NewClient(server.URL, "test-key")
	data := schema.TestResourceDataRaw(t, ServiceResource().Schema, nil)
	data.SetId("123")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	diagnostics := serviceDelete(ctx, data, apiClient)
	if data.Id() != "123" || !diagnostics.HasError() {
		t.Fatalf("unfinished deletion: ID=%q diagnostics=%v", data.Id(), diagnostics)
	}
}

func TestServiceUpdateRejectsPaymentTermChange(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		t.Errorf("unexpected API request: %s %s", request.Method, request.URL.Path)
		http.Error(writer, "unexpected request", http.StatusInternalServerError)
	}))
	defer server.Close()
	apiClient, _ := client.NewClient(server.URL, "test-key")
	data := schema.TestResourceDataRaw(t, ServiceResource().Schema, map[string]interface{}{
		"name": "server", "payment_term": "month",
	})
	data.SetId("123")
	if err := data.Set("payment_term", "year"); err != nil {
		t.Fatal(err)
	}
	diagnostics := serviceUpdate(context.Background(), data, apiClient)
	if !diagnostics.HasError() || !strings.Contains(fmt.Sprint(diagnostics), "payment_term") {
		t.Fatalf("expected payment term diagnostic, got %v", diagnostics)
	}
}

func TestSetServiceGetResponseIPv6Subnets(t *testing.T) {
	data := schema.TestResourceDataRaw(t, ServiceResource().Schema, nil)
	if err := setServiceGetResponse(data, models.Service{IPv6Subnets: []string{"2001:db8::/48"}}); err != nil {
		t.Fatal(err)
	}
	if !data.Get("ipv6_subnets").(*schema.Set).Contains("2001:db8::/48") {
		t.Fatalf("missing IPv6 subnet: %v", data.Get("ipv6_subnets"))
	}
	if data.Get("ipv6_addresses").(*schema.Set).Len() != 0 {
		t.Fatalf("subnet was treated as an address: %v", data.Get("ipv6_addresses"))
	}
}
