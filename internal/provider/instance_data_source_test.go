package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInstanceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test data source without specific instance (list all)
			{
				Config: testAccInstanceDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.multipass_instance.test", "id"),
					resource.TestCheckResourceAttrSet("data.multipass_instance.test", "instances.#"),
				),
			},
		},
	})
}

func TestAccInstanceDataSourceWithName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// First create an instance
			{
				Config: testAccInstanceDataSourceWithResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("multipass_instance.test", "name", "test-datasource"),
					resource.TestCheckResourceAttr("data.multipass_instance.test", "name", "test-datasource"),
					resource.TestCheckResourceAttrSet("data.multipass_instance.test", "instance.name"),
					resource.TestCheckResourceAttrSet("data.multipass_instance.test", "instance.state"),
				),
			},
		},
	})
}

const testAccInstanceDataSourceConfig = `
data "multipass_instance" "test" {}
`

const testAccInstanceDataSourceWithResourceConfig = `
resource "multipass_instance" "test" {
  name   = "test-datasource"
  image  = "22.04"
}

data "multipass_instance" "test" {
  name = multipass_instance.test.name
}
`
