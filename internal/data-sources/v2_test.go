package data_sources

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrb4el/terraform-provider-aeza/internal/client"
	"github.com/mrb4el/terraform-provider-aeza/internal/models"
)

func TestV2CatalogDataSources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v2/services/operating-systems":
			fmt.Fprint(writer, `[{"id":12,"slug":"ubuntu_2404","name":"Ubuntu","group":"linux","username":"root","enabled":true,"order":1}]`)
		case "/v2/services/operating-systems/recipes":
			fmt.Fprint(writer, `[{"id":23,"slug":"docker","name":"Docker","enabled":true,"operatingSystemSlug":"ubuntu_2404","description":"Docker","icon":"docker"}]`)
		case "/v2/services/products":
			fmt.Fprint(writer, `{"total":1,"items":[{"id":42,"name":"Example VPS","groupId":7,"type":{"slug":"vps"},"rawPrices":{"month":100}}]}`)
		case "/v2/services/types":
			fmt.Fprint(writer, `{"total":1,"items":[{"slug":"vps","name":"VPS","serviceHandler":"vm6"}]}`)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	apiClient, _ := client.NewClient(server.URL, "test-key")

	for _, test := range []struct {
		name     string
		resource *schema.Resource
		read     schema.ReadContextFunc
		field    string
	}{
		{name: "operating systems", resource: OSDataSource(), read: osRead, field: "os_list"},
		{name: "recipes", resource: RecipesDataSource(), read: recipesRead, field: "recipes"},
		{name: "products", resource: ProductsDataSource(), read: productsRead, field: "products"},
	} {
		t.Run(test.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, test.resource.Schema, nil)
			if diagnostics := test.read(context.Background(), data, apiClient); diagnostics.HasError() {
				t.Fatalf("read failed: %v", diagnostics)
			}
			if items := data.Get(test.field).([]interface{}); len(items) != 1 {
				t.Fatalf("unexpected %s: %v", test.field, items)
			}
		})
	}
}

func TestFlattenServicesIPv6Subnets(t *testing.T) {
	data := schema.TestResourceDataRaw(t, ServicesDataSource().Schema, nil)
	if err := flattenTerraformServices(data, []models.TerraformService{{IPv6Subnets: []string{"2001:db8::/48"}}}); err != nil {
		t.Fatal(err)
	}
	services := data.Get("services").([]interface{})
	fields := services[0].(map[string]interface{})
	if !fields["ipv6_subnets"].(*schema.Set).Contains("2001:db8::/48") {
		t.Fatalf("missing IPv6 subnet: %v", fields["ipv6_subnets"])
	}
}