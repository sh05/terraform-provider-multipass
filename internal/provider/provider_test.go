package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"multipass": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add checks for required environment variables or other prerequisites here
	// For example:
	// if v := os.Getenv("MULTIPASS_BINARY_PATH"); v == "" {
	//     t.Fatal("MULTIPASS_BINARY_PATH must be set for acceptance tests")
	// }
}

func TestMain(m *testing.M) {
	// You can use this to setup and teardown for tests
	m.Run()
}
