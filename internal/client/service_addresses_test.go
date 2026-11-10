package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/mrb4el/terraform-provider-aeza/internal/models"
	"github.com/mrb4el/terraform-provider-aeza/internal/interfaces"
)

func TestGetServiceFetchesIPv6Addresses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v2/services/123":
			fmt.Fprint(writer, `{"id":123,"name":"server","productId":42,"product":{"id":42,"name":"VPS","typeSlug":"vps","groupId":7,"prices":{"base":{"month":100}}},"ip":"192.0.2.1","paymentTerm":"month","parameters":{"os":"ubuntu_2404"},"status":"active","capabilities":["ip"]}`)
		case "/v2/services/123/networks/ipv4":
			fmt.Fprint(writer, `{"total":2,"items":[{"value":"192.0.2.1"},{"value":"192.0.2.2"}]}`)
		case "/v2/services/123/networks/ipv6":
			if request.Header.Get("X-API-Key") != "test-key" {
				t.Error("missing API key")
			}
			fmt.Fprint(writer, `{"total":1,"items":[{"key":"1","value":"2001:db8::/64","prefix":64,"gateway":"2001:db8::1","ips":["2001:db8::10","2001:db8::11","2001:db8::10"]}]}`)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client, _ := NewClient(server.URL, "test-key")
	service, err := client.GetService(context.Background(), 123)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"2001:db8::10", "2001:db8::11"}; !reflect.DeepEqual(service.IPv6Addresses, want) {
		t.Fatalf("IPv6 addresses = %v, want %v", service.IPv6Addresses, want)
	}
	if want := []string{"192.0.2.1", "192.0.2.2"}; !reflect.DeepEqual(service.IPv4Addresses, want) {
		t.Fatalf("IPv4 addresses = %v, want %v", service.IPv4Addresses, want)
	}
	if service.ID != 123 || service.ProductID != 42 || service.Product.Name != "VPS" || service.Parameters["os"] != "ubuntu_2404" {
		t.Fatalf("unexpected service fields: %+v", service)
	}
}

func TestGetServiceIPv6RequestFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/v2/services/123" {
			fmt.Fprint(writer, `{"id":123,"capabilities":["ip"]}`)
			return
		}
		if request.URL.Path == "/v2/services/123/networks/ipv4" {
			fmt.Fprint(writer, `{"total":0,"items":[]}`)
			return
		}
		http.NotFound(writer, request)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL, "test-key")
	_, err := client.GetService(context.Background(), 123)
	if err == nil || !strings.Contains(err.Error(), "IPv6") || errors.Is(err, interfaces.ErrServiceNotFound) {
		t.Fatalf("unexpected IPv6 request error: %v", err)
	}
}

func TestGetServiceNotFound(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	_, err := client.GetService(context.Background(), 123)
	if !errors.Is(err, interfaces.ErrServiceNotFound) {
		t.Fatalf("expected service not found error, got %v", err)
	}
}

func TestGetServiceWithoutIPCapability(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v2/services/123" {
			t.Errorf("unexpected network request: %s", request.URL.Path)
			http.NotFound(writer, request)
			return
		}
		fmt.Fprint(writer, `{"id":123,"ip":"192.0.2.1","capabilities":["rename"]}`)
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	service, err := client.GetService(context.Background(), 123)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"192.0.2.1"}; !reflect.DeepEqual(service.IPv4Addresses, want) {
		t.Fatalf("IPv4 addresses = %v, want %v", service.IPv4Addresses, want)
	}
}

func TestGetServicePreservesPrimaryIPWhenNetworksEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v2/services/123":
			fmt.Fprint(writer, `{"id":123,"ip":"192.0.2.1","capabilities":["ip"]}`)
		case "/v2/services/123/networks/ipv4", "/v2/services/123/networks/ipv6":
			fmt.Fprint(writer, `{"total":0,"items":[]}`)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	service, err := client.GetService(context.Background(), 123)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"192.0.2.1"}; !reflect.DeepEqual(service.IPv4Addresses, want) {
		t.Fatalf("IPv4 addresses = %v, want %v", service.IPv4Addresses, want)
	}
}

func TestGetServiceReturnsIPv6SubnetWithoutHostAddresses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v2/services/123":
			fmt.Fprint(writer, `{"id":123,"ip":"192.0.2.1","capabilities":["ip"]}`)
		case "/v2/services/123/networks/ipv4":
			fmt.Fprint(writer, `{"total":1,"items":[{"value":"192.0.2.1"}]}`)
		case "/v2/services/123/networks/ipv6":
			fmt.Fprint(writer, `{"total":1,"items":[{"value":"2001:db8::/48","prefix":"48","ips":[]}]}`)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	service, err := client.GetService(context.Background(), 123)
	if err != nil {
		t.Fatal(err)
	}
	if len(service.IPv6Addresses) != 0 || !reflect.DeepEqual(service.IPv6Subnets, []string{"2001:db8::/48"}) {
		t.Fatalf("unexpected IPv6 addresses %v or subnets %v", service.IPv6Addresses, service.IPv6Subnets)
	}
}

func TestControlServiceV2RequestBody(t *testing.T) {
	for _, action := range []string{"resume", "suspend", "restart"} {
		t.Run(action, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if request.Method != http.MethodPost || request.URL.Path != "/v2/services/123/ctl/"+action {
					t.Errorf("unexpected request: %s %s", request.Method, request.URL.Path)
				}
				if action != "resume" {
					var body map[string]bool
					if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
						t.Errorf("decode request body: %v", err)
					} else if force, ok := body["force"]; !ok || force {
						t.Errorf("unexpected force value: %v", body)
					}
				}
				writer.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()
			client, _ := NewClient(server.URL, "test-key")
			if err := client.ControlService_V2(context.Background(), 123, action); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestListServicesV2Pagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v2/services" || request.URL.Query().Get("limit") != "50" {
			t.Errorf("unexpected request: %s", request.URL.String())
		}
		switch request.URL.Query().Get("offset") {
		case "0":
			fmt.Fprint(writer, `{"total":2,"items":[{"id":1,"ip":"192.0.2.1"}]}`)
		case "1":
			fmt.Fprint(writer, `{"total":2,"items":[{"id":2,"ip":"192.0.2.2"}]}`)
		default:
			t.Errorf("unexpected offset: %s", request.URL.Query().Get("offset"))
			http.Error(writer, "unexpected offset", http.StatusBadRequest)
		}
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	services, err := client.ListServices(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(services) != 2 || services[0].ID != 1 || services[1].ID != 2 {
		t.Fatalf("unexpected services: %+v", services)
	}
}

func TestListServicesV2FetchesNetworkAddresses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v2/services":
			fmt.Fprint(writer, `{"total":2,"items":[{"id":1,"ip":"192.0.2.1","capabilities":["ip"]},{"id":2,"capabilities":["rename"]}]}`)
		case "/v2/services/1/networks/ipv4":
			fmt.Fprint(writer, `{"items":[{"value":"192.0.2.1"},{"value":"192.0.2.2"}],"total":2}`)
		case "/v2/services/1/networks/ipv6":
			fmt.Fprint(writer, `{"items":[{"value":"2001:db8::/48","prefix":"48","ips":["2001:db8::1"]}],"total":1}`)
		default:
			t.Errorf("unexpected request: %s", request.URL.Path)
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	services, err := client.ListServices(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(services) != 2 || !reflect.DeepEqual(services[0].IPv4Addresses, []string{"192.0.2.1", "192.0.2.2"}) || !reflect.DeepEqual(services[0].IPv6Addresses, []string{"2001:db8::1"}) || !reflect.DeepEqual(services[0].IPv6Subnets, []string{"2001:db8::/48"}) {
		t.Fatalf("unexpected services: %+v", services)
	}
}

func TestListServiceGroupsV2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v2/services/groups" || request.URL.Query().Get("type") != "vps" {
			t.Errorf("unexpected request: %s", request.URL.String())
		}
		fmt.Fprint(writer, `{"total":1,"items":[{"id":42,"name":"Example","type":{"slug":"vps","name":"VPS","serviceHandler":"vps"},"description":"test","payload":{},"localedPayload":{}}]}`)
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	groups, err := client.ListServiceGroups(context.Background(), "vps")
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 || groups[0].ID != 42 || groups[0].Type != "vps" {
		t.Fatalf("unexpected groups: %+v", groups)
	}
}

func TestListServiceTypesV2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v2/services/types" || request.URL.Query().Get("limit") != "50" {
			t.Errorf("unexpected request: %s", request.URL.String())
		}
		switch request.URL.Query().Get("offset") {
		case "0":
			fmt.Fprint(writer, `{"total":2,"items":[{"slug":"vps","name":"VPS"}]}`)
		case "1":
			fmt.Fprint(writer, `{"total":2,"items":[{"slug":"server","name":"Server"}]}`)
		default:
			http.Error(writer, "unexpected offset", http.StatusBadRequest)
		}
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	types, err := client.ListServiceTypes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 2 || types[0].Slug != "vps" || types[1].Name != "Server" {
		t.Fatalf("unexpected service types: %+v", types)
	}
}

func TestUpdateServiceV2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPatch || request.URL.Path != "/v2/services/123" {
			t.Errorf("unexpected request: %s %s", request.Method, request.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
		} else if body["name"] != "server" || body["autoProlong"] != false || len(body) != 2 {
			t.Errorf("unexpected body: %v", body)
		}
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	name, autoProlong := "server", false
	if err := client.UpdateService(context.Background(), 123, models.ServiceUpdateRequest{Name: &name, AutoProlong: &autoProlong}); err != nil {
		t.Fatal(err)
	}
}

func TestListOSAndRecipesV2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v2/services/operating-systems":
			fmt.Fprint(writer, `[{"id":12,"slug":"ubuntu_2404","name":"Ubuntu","group":"linux","username":"root","enabled":true,"order":1,"supportsSshKeys":true}]`)
		case "/v2/services/operating-systems/recipes":
			fmt.Fprint(writer, `[{"id":23,"slug":"docker","name":"Docker","group":"linux","enabled":true,"order":1,"operatingSystemSlug":"ubuntu_2404","description":"Docker","icon":"docker"}]`)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	operatingSystems, err := client.ListOS(context.Background())
	if err != nil || len(operatingSystems) != 1 || operatingSystems[0].Slug != "ubuntu_2404" {
		t.Fatalf("unexpected operating systems: %+v, error: %v", operatingSystems, err)
	}
	recipes, err := client.ListRecipes(context.Background())
	if err != nil || len(recipes) != 1 || recipes[0].OperatingSystem != "ubuntu_2404" || recipes[0].OS != 12 {
		t.Fatalf("unexpected recipes: %+v, error: %v", recipes, err)
	}
}

func TestListProductsV2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/v2/services/products":
			fmt.Fprint(writer, `{"total":1,"items":[{"id":42,"name":"Example VPS","groupId":7,"type":{"slug":"vps","name":"VPS"},"parameters":{"cpu":2},"rawPrices":{"month":100,"year":1000},"prices":{"base":{"month":100}},"resources":{"cpu":{}}}]}`)
		case "/v2/services/types":
			fmt.Fprint(writer, `{"total":1,"items":[{"slug":"vps","name":"VPS","serviceHandler":"vm6"}]}`)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	products, err := client.ListProducts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 1 || products[0].ID != 42 || products[0].Type != "vps" || products[0].GroupID != 7 || products[0].Prices["month"] != 100 || products[0].ServiceHandler != "vm6" {
		t.Fatalf("unexpected products: %+v", products)
	}
}

func TestListOSV2Pagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v2/services/operating-systems" || request.URL.Query().Get("limit") != "50" {
			t.Errorf("unexpected request: %s", request.URL.String())
		}
		var response []map[string]interface{}
		switch request.URL.Query().Get("offset") {
		case "0":
			for index := 1; index <= 50; index++ {
				response = append(response, map[string]interface{}{"id": index, "slug": fmt.Sprintf("os-%d", index)})
			}
		case "50":
			response = append(response, map[string]interface{}{"id": 51, "slug": "os-51"})
		default:
			t.Errorf("unexpected offset: %s", request.URL.Query().Get("offset"))
			http.Error(writer, "unexpected offset", http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(writer).Encode(response); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	operatingSystems, err := client.ListOS(context.Background())
	if err != nil || len(operatingSystems) != 51 || operatingSystems[50].ID != 51 {
		t.Fatalf("unexpected operating systems: count=%d, error=%v", len(operatingSystems), err)
	}
}

func TestListOSV2RejectsRepeatedPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		response := make([]map[string]interface{}, 0, 50)
		for index := 1; index <= 50; index++ {
			response = append(response, map[string]interface{}{"id": index, "slug": fmt.Sprintf("os-%d", index)})
		}
		if err := json.NewEncoder(writer).Encode(response); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client, _ := NewClient(server.URL, "test-key")
	_, err := client.ListOS(context.Background())
	if err == nil || !strings.Contains(err.Error(), "duplicate ID") {
		t.Fatalf("expected duplicate page error, got %v", err)
	}
}

func TestDeleteServiceV2WaitsForDeletion(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprintf("missing_%t", missing), func(t *testing.T) {
			var deleted bool
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if request.URL.Path != "/v2/services/123" {
					t.Errorf("unexpected path: %s", request.URL.Path)
				}
				switch request.Method {
				case http.MethodDelete:
					deleted = true
					writer.WriteHeader(http.StatusNoContent)
				case http.MethodGet:
					if !deleted {
						t.Error("status checked before deletion was requested")
					}
					if missing {
						http.NotFound(writer, request)
					} else {
						fmt.Fprint(writer, `{"status":"deleted"}`)
					}
				default:
					t.Errorf("unexpected method: %s", request.Method)
				}
			}))
			defer server.Close()
			client, _ := NewClient(server.URL, "test-key")
			if err := client.DeleteService(context.Background(), 123); err != nil {
				t.Fatal(err)
			}
		})
	}
}