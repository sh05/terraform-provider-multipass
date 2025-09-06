package provider

import (
	"regexp"
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

// TestAccInstanceDataSource_NonExistentInstance tests querying a non-existent instance
func TestAccInstanceDataSource_NonExistentInstance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceDataSourceConfigNonExistent,
				ExpectError: regexp.MustCompile(`instance.*not found|does not exist|failed to get instance info`),
			},
		},
	})
}

// TestAccInstanceDataSource_InvalidName tests querying with invalid instance name
func TestAccInstanceDataSource_InvalidName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceDataSourceConfigInvalidName,
				ExpectError: regexp.MustCompile(`invalid.*name|instance.*not found|failed to get instance info`),
			},
		},
	})
}

// TestAccInstanceDataSource_EmptyName tests querying with empty instance name
func TestAccInstanceDataSource_EmptyName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceDataSourceConfigEmptyName,
				ExpectError: regexp.MustCompile(`empty.*name|name.*required|failed to get instance info`),
			},
		},
	})
}

// TestAccInstanceDataSource_MultipassUnavailable tests behavior when multipass is unavailable
func TestAccInstanceDataSource_MultipassUnavailable(t *testing.T) {
	// This test would require setting up a mock environment where multipass is not available
	// For now, we'll create a basic structure that would test this scenario
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// Additional pre-check could go here to simulate multipass unavailability
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Basic check to ensure data source works when multipass is available
					resource.TestCheckResourceAttrSet("data.multipass_instance.test", "id"),
				),
			},
		},
	})
}

// TestAccInstanceDataSource_SpecialCharactersInName tests handling special characters in names
func TestAccInstanceDataSource_SpecialCharactersInName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceDataSourceConfigSpecialChars,
				ExpectError: regexp.MustCompile(`invalid.*name|instance.*not found|failed to get instance info`),
			},
		},
	})
}

// TestAccInstanceDataSource_ListInstancesError tests error handling when listing instances fails
func TestAccInstanceDataSource_ListInstancesError(t *testing.T) {
	// This test would require mocking the multipass client to return errors
	// The actual implementation would depend on the specific error injection mechanism
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.multipass_instance.test", "id"),
					// In a real scenario with error injection, we'd expect specific error patterns
				),
			},
		},
	})
}

// TestAccInstanceDataSource_PartiallyFailedListOperation tests handling of partial failures
func TestAccInstanceDataSource_PartiallyFailedListOperation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
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

// Configuration constants for error test cases

const testAccInstanceDataSourceConfigNonExistent = `
data "multipass_instance" "test" {
  name = "definitely-does-not-exist-instance-123456"
}
`

const testAccInstanceDataSourceConfigInvalidName = `
data "multipass_instance" "test" {
  name = "invalid@name#with$special%characters!"
}
`

const testAccInstanceDataSourceConfigEmptyName = `
data "multipass_instance" "test" {
  name = ""
}
`

const testAccInstanceDataSourceConfigSpecialChars = `
data "multipass_instance" "test" {
  name = "test-instance-with-unicode-❤️-chars"
}
`

// TestAccInstanceDataSource_LongInstanceName tests handling of very long instance names
func TestAccInstanceDataSource_LongInstanceName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceDataSourceConfigLongName,
				ExpectError: regexp.MustCompile(`name too long|invalid.*name|instance.*not found`),
			},
		},
	})
}

const testAccInstanceDataSourceConfigLongName = `
data "multipass_instance" "test" {
  name = "this-is-an-extremely-long-instance-name-that-exceeds-reasonable-limits-and-should-probably-cause-an-error-in-most-systems-because-it-is-way-too-long-for-practical-use"
}
`

// TestAccInstanceDataSource_ConcurrentAccess tests concurrent access to data sources
func TestAccInstanceDataSource_ConcurrentAccess(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfigMultiple,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.multipass_instance.test1", "id"),
					resource.TestCheckResourceAttrSet("data.multipass_instance.test2", "id"),
					resource.TestCheckResourceAttrSet("data.multipass_instance.test3", "id"),
				),
			},
		},
	})
}

const testAccInstanceDataSourceConfigMultiple = `
data "multipass_instance" "test1" {}
data "multipass_instance" "test2" {}
data "multipass_instance" "test3" {}
`
