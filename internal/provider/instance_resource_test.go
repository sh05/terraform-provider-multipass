package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInstanceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccInstanceResourceConfig("test-instance"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("multipass_instance.test", "name", "test-instance"),
					resource.TestCheckResourceAttrSet("multipass_instance.test", "id"),
					resource.TestCheckResourceAttrSet("multipass_instance.test", "state"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "multipass_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore cloud_init as it's not stored in state after creation
				ImportStateVerifyIgnore: []string{"cloud_init"},
			},
			// Update and Read testing (this will trigger replacement due to schema)
			{
				Config: testAccInstanceResourceConfigUpdated("test-instance-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("multipass_instance.test", "name", "test-instance-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccInstanceResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "multipass_instance" "test" {
  name   = "%s"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}
`, name)
}

func testAccInstanceResourceConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "multipass_instance" "test" {
  name   = "%s"
  image  = "22.04"
  cpu    = "2"
  memory = "2G"
  disk   = "10G"
}
`, name)
}
