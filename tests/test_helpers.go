// Package tests provides common test utilities and helper functions for the EVE-NG Terraform provider.
package tests

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eveng "github.com/nawada0615/terraform-provider-eve-ng/eve-ng"
)

// getProviderFactories returns the provider factories for testing
func getProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"eve": func() (*schema.Provider, error) {
			provider := eveng.Provider()
			if provider == nil {
				return nil, fmt.Errorf("failed to create provider")
			}
			return provider, nil
		},
	}
}

// createTestConfig creates a common test configuration with provider and lab
func createTestConfig(serverURL, resourceConfig string) string {
	return fmt.Sprintf(`
		provider "eve" {
			endpoint = "%s"
			username = "testuser"
			password = "testpass"
			insecure_skip_verify = true
		}
		resource "eve_lab" "test" {
			path = "/"
			name = "test-lab"
			author = "test"
			description = "test lab"
			version = "1"
		}
		%s
	`, serverURL, resourceConfig)
}

// createTestStep creates a common test step with the given config and checks
func createTestStep(config string, checks []resource.TestCheckFunc) resource.TestStep {
	return resource.TestStep{
		Config: config,
		Check:  resource.ComposeTestCheckFunc(checks...),
	}
}

// runResourceTest runs a resource test with the given server, config, and checks
func runResourceTest(t *testing.T, server *httptest.Server, resourceConfig string, checks []resource.TestCheckFunc) {
	defer server.Close()

	config := createTestConfig(server.URL, resourceConfig)
	resource.Test(t, resource.TestCase{
		ProviderFactories: getProviderFactories(),
		Steps: []resource.TestStep{
			createTestStep(config, checks),
		},
	})
}
