package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestGetServiceFetchesIPv6Addresses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/services/123":
			fmt.Fprint(writer, `{"data":{"items":[{"id":123,"ip":"192.0.2.1"}]}}`)
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
	if want := []string{"192.0.2.1"}; !reflect.DeepEqual(service.IPv4Addresses, want) {
		t.Fatalf("IPv4 addresses = %v, want %v", service.IPv4Addresses, want)
	}
}

func TestGetServiceIPv6RequestFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/services/123" {
			fmt.Fprint(writer, `{"data":{"items":[{"id":123}]}}`)
			return
		}
		http.NotFound(writer, request)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL, "test-key")
	_, err := client.GetService(context.Background(), 123)
	if err == nil || !strings.Contains(err.Error(), "IPv6") || strings.Contains(err.Error(), "not found") {
		t.Fatalf("unexpected IPv6 request error: %v", err)
	}
}