package utils

import (
	"reflect"
	"testing"

	"github.com/scinfra-pro/terraform-provider-aeza/internal/models/legacy"
)

func TestConvertLegacyServiceGetToTerraformAddresses(t *testing.T) {
	service := legacy.ServiceGet{ServiceVPS: legacy.ServiceVPS{
		Service: legacy.Service{IP: "192.0.2.10"},
		IPs: []legacy.IPAddress{
			{Value: "192.0.2.10"},
			{Value: "192.0.2.11"},
			{Value: ""},
		},
		IPv6: []legacy.IPv6Address{
			{
				Value: "2001:db8::10",
				IPs: []interface{}{
					"2001:db8::11",
					map[string]interface{}{"value": "2001:db8::12"},
					map[string]interface{}{"address": "2001:db8::13"},
				},
			},
		},
	}}

	got := ConvertLegacyServiceGetToTerraform(service)
	wantIPv4 := []string{"192.0.2.10", "192.0.2.11"}
	wantIPv6 := []string{"2001:db8::10", "2001:db8::11", "2001:db8::12", "2001:db8::13"}

	if !reflect.DeepEqual(got.IPv4Addresses, wantIPv4) {
		t.Fatalf("unexpected IPv4 addresses: got %v, want %v", got.IPv4Addresses, wantIPv4)
	}
	if !reflect.DeepEqual(got.IPv6Addresses, wantIPv6) {
		t.Fatalf("unexpected IPv6 addresses: got %v, want %v", got.IPv6Addresses, wantIPv6)
	}
	if got.IP != wantIPv4[0] {
		t.Fatalf("unexpected primary IP: got %q, want %q", got.IP, wantIPv4[0])
	}
}

func TestConvertLegacyServiceGetUsesFirstIPv4AsPrimary(t *testing.T) {
	service := legacy.ServiceGet{ServiceVPS: legacy.ServiceVPS{
		IPs: []legacy.IPAddress{{Value: "192.0.2.20"}},
	}}

	got := ConvertLegacyServiceGetToTerraform(service)
	if got.IP != "192.0.2.20" {
		t.Fatalf("unexpected primary IP: %q", got.IP)
	}
}
