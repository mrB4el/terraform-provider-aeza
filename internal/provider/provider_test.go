package provider

import "testing"

func TestProviderInternalValidate(t *testing.T) {
	if err := New().InternalValidate(); err != nil {
		t.Fatalf("provider schema is invalid: %v", err)
	}
}
