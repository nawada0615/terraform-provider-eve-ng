package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eveng "github.com/nawada0615/terraform-provider-eve-ng/eve-ng"
)

func TestProvider(t *testing.T) {
	if err := eveng.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(_ *testing.T) {
	var _ *schema.Provider = eveng.Provider()
}

func TestProviderSchema(t *testing.T) {
	provider := eveng.Provider()

	// Test required fields
	requiredFields := []string{"endpoint", "username", "password"}
	for _, field := range requiredFields {
		if provider.Schema[field] == nil {
			t.Errorf("Required field %s is missing from provider schema", field)
		}
		if !provider.Schema[field].Required {
			t.Errorf("Field %s should be required", field)
		}
	}

	// Test optional fields
	optionalFields := []string{"insecure_skip_verify", "timeout"}
	for _, field := range optionalFields {
		if provider.Schema[field] == nil {
			t.Errorf("Optional field %s is missing from provider schema", field)
		}
		if !provider.Schema[field].Optional {
			t.Errorf("Field %s should be optional", field)
		}
	}
}
