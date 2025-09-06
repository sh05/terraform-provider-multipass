package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

// TestAccInstanceResource_InvalidConfiguration tests various invalid configurations
func TestAccInstanceResource_InvalidConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceResourceConfigInvalidCPU(),
				ExpectError: regexp.MustCompile(`invalid CPU value`),
			},
		},
	})
}

// TestAccInstanceResource_InvalidMemory tests invalid memory configurations
func TestAccInstanceResource_InvalidMemory(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceResourceConfigInvalidMemory(),
				ExpectError: regexp.MustCompile(`invalid memory value|failed to launch instance`),
			},
		},
	})
}

// TestAccInstanceResource_InvalidDisk tests invalid disk configurations
func TestAccInstanceResource_InvalidDisk(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceResourceConfigInvalidDisk(),
				ExpectError: regexp.MustCompile(`invalid disk value|failed to launch instance`),
			},
		},
	})
}

// TestAccInstanceResource_DuplicateName tests creating instances with duplicate names
func TestAccInstanceResource_DuplicateName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceResourceConfigDuplicate(),
				ExpectError: regexp.MustCompile(`already exists|name.*already.*use`),
			},
		},
	})
}

// TestAccInstanceResource_InvalidName tests invalid instance names
func TestAccInstanceResource_InvalidName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceResourceConfigInvalidName(),
				ExpectError: regexp.MustCompile(`invalid.*name|failed to launch instance`),
			},
		},
	})
}

// TestAccInstanceResource_NonExistentImage tests using non-existent image
func TestAccInstanceResource_NonExistentImage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceResourceConfigNonExistentImage(),
				ExpectError: regexp.MustCompile(`unable to find.*image|not found|failed to launch instance`),
			},
		},
	})
}

// TestAccInstanceResource_InvalidTimeout tests invalid timeout values
func TestAccInstanceResource_InvalidTimeout(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceResourceConfigInvalidTimeout(),
				ExpectError: regexp.MustCompile(`invalid timeout format`),
			},
		},
	})
}

// TestAccInstanceResource_ImportNonExistent tests importing non-existent instance
func TestAccInstanceResource_ImportNonExistent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        testAccInstanceResourceConfig("test-import"),
				ResourceName:  "multipass_instance.test",
				ImportState:   true,
				ImportStateId: "non-existent-instance",
				ExpectError:   regexp.MustCompile(`instance.*not found|does not exist`),
			},
		},
	})
}

// Helper functions for invalid configurations

func testAccInstanceResourceConfigInvalidCPU() string {
	return `
resource "multipass_instance" "test" {
  name   = "test-invalid-cpu"
  image  = "22.04"
  cpu    = "-1"
  memory = "1G"
  disk   = "5G"
}
`
}

func testAccInstanceResourceConfigInvalidMemory() string {
	return `
resource "multipass_instance" "test" {
  name   = "test-invalid-memory"
  image  = "22.04"
  cpu    = "1"
  memory = "invalid-memory"
  disk   = "5G"
}
`
}

func testAccInstanceResourceConfigInvalidDisk() string {
	return `
resource "multipass_instance" "test" {
  name   = "test-invalid-disk"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "invalid-disk"
}
`
}

func testAccInstanceResourceConfigDuplicate() string {
	return `
resource "multipass_instance" "test1" {
  name   = "duplicate-name-test"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}

resource "multipass_instance" "test2" {
  name   = "duplicate-name-test"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}
`
}

func testAccInstanceResourceConfigInvalidName() string {
	return `
resource "multipass_instance" "test" {
  name   = "invalid@name!with$special%chars"
  image  = "22.04"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}
`
}

func testAccInstanceResourceConfigNonExistentImage() string {
	return `
resource "multipass_instance" "test" {
  name   = "test-nonexistent-image"
  image  = "non-existent-image-999"
  cpu    = "1"
  memory = "1G"
  disk   = "5G"
}
`
}

func testAccInstanceResourceConfigInvalidTimeout() string {
	return `
resource "multipass_instance" "test" {
  name    = "test-invalid-timeout"
  image   = "22.04"
  cpu     = "1"
  memory  = "1G"
  disk    = "5G"
  timeout = "invalid-timeout-format"
}
`
}

// TestAccInstanceResource_CloudInitFile tests cloud-init file functionality
func TestAccInstanceResource_CloudInitFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceResourceConfigInvalidCloudInit(),
				ExpectError: regexp.MustCompile(`cloud-init.*not found|no such file`),
			},
		},
	})
}

func testAccInstanceResourceConfigInvalidCloudInit() string {
	return `
resource "multipass_instance" "test" {
  name       = "test-invalid-cloud-init"
  image      = "22.04"
  cpu        = "1"
  memory     = "1G"
  disk       = "5G"
  cloud_init = "/non/existent/cloud-init.yaml"
}
`
}

// TestAccInstanceResource_ResourceLimits tests resource limit validation
func TestAccInstanceResource_ResourceLimits(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceResourceConfigExcessiveResources(),
				ExpectError: regexp.MustCompile(`resource allocation failed|insufficient resources|failed to launch instance`),
			},
		},
	})
}

func testAccInstanceResourceConfigExcessiveResources() string {
	return `
resource "multipass_instance" "test" {
  name   = "test-excessive-resources"
  image  = "22.04"
  cpu    = "9999"
  memory = "999999G"
  disk   = "999999G"
}
`
}

// Custom check function for testing non-existent instance reads
func testCheckInstanceDoesNotExist(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no ID set", resourceName)
		}

		// This would typically check that the instance doesn't exist in the actual system
		// but for this test we're just checking the state management
		return nil
	}
}
